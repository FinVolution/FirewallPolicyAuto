package h3c_v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"
	"github.com/FinVolution/FirewallPolicyAuto/service/utils"
)

// H3C防火墙
// 策略列表接口返回的源地址、源地址组、目的地址、目的地址组、服务端口、服务端口组是有数量限制的，最多显示7个。
// 完整的数据需要调用相关接口获取数据，根据策略id进行的关联
// H3C Firewall
// The policy list interface returns a limited number of source addresses, source address groups, destination addresses, destination address groups, service ports, and service port groups, with a maximum display of 7.
// Complete data needs to be obtained by calling the relevant interface to retrieve the data, associated based on the policy ID.

// GetPolicyList 获取策略列表
func (firewall *FirewallH3CV1) ListPolicy(filters map[string]string) (policyList []dto.Policy, err error) {
	rules, err := firewall.getRules()
	if err != nil {
		return
	}

	// if rules list is empty, return
	if len(rules.GetRules) == 0 {
		return
	}

	// 获取地址组名和IP地址 map group:addr
	groupAddrMap, err := firewall.getPolicyGroupAddress()
	if err != nil {
		return
	}
	// 获取服务组名和端口 map group:service port
	groupServiceMap, err := firewall.getPolicyGroupServicePort()
	if err != nil {
		return
	}

	// 获取每条策略包含的所有源地址组 map policy_id:group
	srcIDGroupMap, err := firewall.getSrcPolicyIDGroupMap()
	if err != nil {
		return
	}

	// 获取每条策略包含的所有目的地址组 map policy_id:group
	destIDGroupMap, err := firewall.getDestPolicyIDGroupMap()
	if err != nil {
		return
	}

	// 获取每条策略包含的所有服务组 map policy_id:serviceobj
	servIDGroupMap, err := firewall.getServicePolicyIDGroupMap()
	if err != nil {
		return
	}

	// 获取每条策略直接包含的所有源ip地址 simple address
	policyIDSrcSimpleAddrMap, err := firewall.getSrcPolicySimpleAddrMap()
	if err != nil {
		return
	}

	// 获取信息  每条策略直接包含的所有目的ip地址 simple address
	policyIDDestSimpleAddrMap, err := firewall.getDestPolicySimpleAddrMap()
	if err != nil {
		return
	}

	// 获取每条策略直接包含的所有服务端口 simple service port
	policyIDServPortMap, err := firewall.getServPolicyPortMap()
	if err != nil {
		return
	}

	// 完善策略信息
	for _, rule := range rules.GetRules {
		// 获取源地址组
		srcAddressList := []string{}
		for _, group := range srcIDGroupMap[rule.ID] {
			srcAddressList = append(srcAddressList, groupAddrMap[group]...)
		}
		// 获取simple源地址
		srcAddressList = append(srcAddressList, policyIDSrcSimpleAddrMap[rule.ID]...)
		// 数据去重
		srcAddressList = utils.RemoveDuplicateElement(srcAddressList)
		srcAddr, ok := filters["srcAddr"]
		if ok && srcAddr != "" && !utils.ContainsAny(srcAddressList, strings.Split(srcAddr, ",")) {
			continue
		}

		// 获取目的地址组
		destAddressList := []string{}
		for _, group := range destIDGroupMap[rule.ID] {
			destAddressList = append(destAddressList, groupAddrMap[group]...)
		}
		// 获取simple目的地址
		destAddressList = append(destAddressList, policyIDDestSimpleAddrMap[rule.ID]...)
		dstAddr, ok := filters["dstAddr"]
		if ok && dstAddr != "" && !utils.ContainsAny(destAddressList, strings.Split(dstAddr, ",")) {
			continue
		}

		// 服务端口服务端口
		servicePortList := []string{}
		for _, group := range servIDGroupMap[rule.ID] {
			servicePortList = append(servicePortList, groupServiceMap[group]...)
		}
		// 获取simple服务端口
		servicePortList = append(servicePortList, policyIDServPortMap[rule.ID]...)
		servicePort, ok := filters["servicePort"]
		if ok && servicePort != "" && !utils.ContainsAny(servicePortList, strings.Split(servicePort, ",")) {
			continue
		}

		// 源区域
		srcZone, ok := filters["srcZone"]
		if ok && srcZone != "" && !utils.ContainsAny(rule.SrcZoneList["DestZoneItem"], strings.Split(servicePort, ",")) {
			continue
		}

		// 目的区域
		destZone, ok := filters["destZone"]
		if ok && destZone != "" && !utils.ContainsAny(rule.DestZoneList["DestZoneItem"], strings.Split(servicePort, ",")) {
			continue
		}

		// 策略查询结果
		policyList = append(policyList, dto.Policy{
			ID:           rule.ID,
			Name:         rule.Name,
			Action:       rule.Action,
			Enable:       rule.Enable,
			SrcZone:      utils.AdditionalPolicyItem(rule.SrcZoneList["DestZoneItem"]),
			DestZone:     utils.AdditionalPolicyItem(rule.DestZoneList["DestZoneItem"]),
			SrcAddress:   utils.AdditionalPolicyItem(srcAddressList),
			DestAddress:  utils.AdditionalPolicyItem(destAddressList),
			ServicePort:  utils.AdditionalPolicyItem(servicePortList),
			FirewallName: firewall.Name,
		})

	}
	return
}

func (firewall *FirewallH3CV1) CreatePolicy(params dto.CreatePolicyParams) (err error) {

	now := fmt.Sprintf("-%s-%s", params.Title, time.Now().Format("20060102"))
	srcAddrGroupName := "SRC" + now     // 源地址组名称
	destAddrGroupName := "DEST" + now   // 目的地址组名称
	serviceGroupName := "SERVICE" + now // 服务组名称
	// 创建源地址组
	if err = firewall.createPolicyAddressGroup(srcAddrGroupName); err != nil {
		return
	}
	// 创建目的地址组
	if err = firewall.createPolicyAddressGroup(destAddrGroupName); err != nil {
		return
	}
	// 在源地址组中添加ip
	for _, addr := range params.SrcAddr {
		if err = firewall.addPolicyAddressGroup(srcAddrGroupName, addr); err != nil {
			return
		}
	}
	// 在目的地址组中添加ip
	for _, addr := range params.DestAddr {
		if err = firewall.addPolicyAddressGroup(destAddrGroupName, addr); err != nil {
			return
		}
	}
	// 创建服务组
	if err = firewall.createPolicyServiceGroup(serviceGroupName); err != nil {
		return
	}
	// 服务组添加端口
	for _, service := range params.Service {
		if err = firewall.addPolicyServiceGroup(serviceGroupName, service); err != nil {
			return
		}
	}
	// 创建策略规则
	policyName := "POLICY" + now // 策略名称
	policyID, err := firewall.createPolicyRules(policyName, params.Action)
	if err != nil {
		return
	}
	// 策略添加源地址组
	if err = firewall.addPolicySrcGroupToPolicy(srcAddrGroupName, policyID); err != nil {
		return
	}
	// 策略添加目的地址组
	if err = firewall.addPolicyDestGroupToPolicy(destAddrGroupName, policyID); err != nil {
		return
	}
	// 策略添加服务组
	if err = firewall.addPolicyServiceToPolicy(serviceGroupName, policyID); err != nil {
		return
	}
	// 策略添加源区域
	if err = firewall.addPolicySrcZoneToPolicy(params.SrcZone, policyID); err != nil {
		return
	}
	// 策略添加目的区域
	if err = firewall.addPolicyDestZoneToPolicy(params.DestZone, policyID); err != nil {
		return
	}
	return
}
