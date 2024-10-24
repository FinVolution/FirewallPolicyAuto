package fortinet_v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"
	"github.com/FinVolution/FirewallPolicyAuto/service/utils"
)

// Fortinet防火墙
// 服务端口或者服务端口组、地址或者地址组都有一个别名name，这个别名是唯一的，在获取的策略列表中，地址字段和地址组展示的都是别名
// 要获取到真实的数据需要调用另外的接口获取到name和ip/service/port的映射关系
// Fortinet firewall
// Service ports or service port groups, addresses or address groups all have a unique alias name. In the obtained policy list, the address field and address groups display the alias.
// To obtain the actual data, you need to call another interface to get the mapping relationship between the name and the IP/service/port.

// ListPolicy
func (firewall *FirewallFortinetV1) ListPolicy(filters map[string]string) (policyList []dto.Policy, err error) {
	// get policy list
	policyResponse, err := firewall.getPolicy()
	if err != nil {
		return
	}
	// if policy list is empty, return
	if len(policyResponse.Results) == 0 {
		return
	}
	// 获取策略的地址相关信息  地址名称:IP地址()
	addressNameMap, err := firewall.getPolicyNameAddressMap()
	if err != nil {
		return
	}
	// 获取策略地址组相关信息 地址组:地址名称()
	addressGroupNameMap, err := firewall.getPolicyGroupNameAddressMap()
	if err != nil {
		return
	}
	// 映射地址组名到地址名
	for k, v := range addressGroupNameMap {
		for _, vv := range v {
			addressNameMap[k] = append(addressNameMap[k], addressNameMap[vv]...)
		}
	}
	// 获取服务端口相关信息 服务名称:协议/端口(ALL_TCP:tcp:1-65535)
	serviceNamePortMap, err := firewall.getPolicyNameServicePortMap()
	if err != nil {
		return
	}
	// 获取服务端口相关信息 服务组名称：服务名称
	serviceGroupNameMap, err := firewall.getPolicyGroupNameServiceMap()
	if err != nil {
		return
	}
	// 映射服务组名到服务名称
	for k, v := range serviceGroupNameMap {
		for _, vv := range v {
			serviceNamePortMap[k] = append(serviceNamePortMap[k], serviceNamePortMap[vv]...)
		}
	}
	// 完善策略信息
	for _, fortinetPolicy := range policyResponse.Results {
		// 源地址
		srcAddressList := []string{}
		for _, item := range fortinetPolicy.SrcAddr {
			srcAddressList = append(srcAddressList, addressNameMap[item.Name]...)
		}
		srcAddressList = utils.RemoveDuplicateElement(srcAddressList)
		srcAddr, ok := filters["srcAddr"]
		if ok && srcAddr != "" && !utils.ContainsAny(srcAddressList, strings.Split(srcAddr, ",")) {
			continue
		}

		// 目的地址
		destAddressList := []string{}
		for _, item := range fortinetPolicy.DstAddr {
			destAddressList = append(destAddressList, addressNameMap[item.Name]...)
		}
		destAddressList = utils.RemoveDuplicateElement(destAddressList)
		dstAddr, ok := filters["dstAddr"]
		if ok && dstAddr != "" && !utils.ContainsAny(destAddressList, strings.Split(dstAddr, ",")) {
			continue
		}

		// 服务端口
		servicePortList := []string{}
		for _, item := range fortinetPolicy.Service {
			servicePortList = append(servicePortList, serviceNamePortMap[item.Name]...)
		}
		servicePortList = utils.RemoveDuplicateElement(servicePortList)
		servicePort, ok := filters["servicePort"]
		if ok && servicePort != "" && !utils.ContainsAny(servicePortList, strings.Split(servicePort, ",")) {
			continue
		}

		// 源区域
		srcZoneList := []string{}
		for _, v := range fortinetPolicy.SrcIntf {
			srcZoneList = append(srcZoneList, v.Name)
		}
		srcZone, ok := filters["srcZone"]
		if ok && srcZone != "" && !utils.ContainsAny(srcZoneList, strings.Split(srcZone, ",")) {
			continue
		}

		// 目的区域
		destZoneList := []string{}
		for _, v := range fortinetPolicy.DstIntf {
			destZoneList = append(destZoneList, v.Name)
		}
		destZone, ok := filters["destZone"]
		if ok && destZone != "" && !utils.ContainsAny(destZoneList, strings.Split(destZone, ",")) {
			continue
		}

		// 策略查询结果
		policyList = append(policyList, dto.Policy{
			ID:           fortinetPolicy.ID,
			Name:         fortinetPolicy.Name,
			Action:       actionExchange(fortinetPolicy.Action),
			Enable:       statusExchange(fortinetPolicy.Status),
			SrcZone:      utils.AdditionalPolicyItem(srcZoneList),
			DestZone:     utils.AdditionalPolicyItem(destZoneList),
			SrcAddress:   utils.AdditionalPolicyItem(srcAddressList),
			DestAddress:  utils.AdditionalPolicyItem(destAddressList),
			ServicePort:  utils.AdditionalPolicyItem(servicePortList),
			FirewallName: firewall.Name,
		})
	}
	return
}

func (firewall *FirewallFortinetV1) CreatePolicy(params dto.CreatePolicyParams) (err error) {

	now := fmt.Sprintf("-%s-%s", params.Title, time.Now().Format("20060102"))
	srcAddrGroupName := "SRC" + now     // 源地址组名称
	destAddrGroupName := "DEST" + now   // 目的地址组名称
	serviceGroupName := "SERVICE" + now // 服务组名称
	// 创建源地址
	for _, addr := range params.SrcAddr {
		if err := firewall.createPolicyAddress(addr); err != nil {
			return err
		}
	}
	// 创建目的地址
	for _, addr := range params.DestAddr {
		if err := firewall.createPolicyAddress(addr); err != nil {
			return err
		}
	}
	// 创建源地址组&添加源地址到地址组
	if err := firewall.addPolicyAddrToGroup(srcAddrGroupName, params.SrcAddr); err != nil {
		return err
	}
	// 创建目的地址组&添加目的地址到地址组
	if err := firewall.addPolicyAddrToGroup(destAddrGroupName, params.DestAddr); err != nil {
		return err
	}
	// 创建服务
	for _, service := range params.Service {
		if err := firewall.createPolicyService(service); err != nil {
			return err
		}
	}
	// 创建服务组&添加服务到服务组
	if err := firewall.addPolicyServiceToGroup(serviceGroupName, params.Service); err != nil {
		return err
	}

	policyName := "POLICY" + now // 策略名称
	action := actionAccept
	if params.Action == 1 {
		action = actionDeny
	}
	// 创建策略
	if err := firewall.addPolicy(policyName, action, params.SrcZone, params.DestZone, srcAddrGroupName, destAddrGroupName, serviceGroupName); err != nil {
		return err
	}
	return
}
