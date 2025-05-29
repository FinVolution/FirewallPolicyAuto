package cisco_v1

import (
	"fmt"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"
	"github.com/FinVolution/FirewallPolicyAuto/service/utils"
	"strconv"
	"strings"
)

func (firewall *FirewallCiscoV1) ListPolicy(filters map[string]string) (policyList []dto.Policy, err error) {
	// get policy list  name
	policyNameList, err := firewall.GetPolicyACLName()
	if err != nil {
		return
	}
	// if policy list is empty, return
	if len(policyNameList) == 0 {
		return
	}

	for _, aclName := range policyNameList {
		// 获取每个acl的策略
		var aclItems []ciscoPolicy
		aclItems, err = firewall.GetFirewallAclItems(aclName, 0)
		if err != nil {
			return
		}
		fmt.Printf("%s : %d\n", aclName, len(aclItems))
		// 处理每条策略
		for _, policy := range aclItems {
			id, err := strconv.ParseInt(policy.ObjectId, 10, 64)

			// 处理源地址
			var srcAddressList []string
			srcAddressList, err = firewall.DealACLSourceOrDestAddr(policy.SourceAddress)
			if err != nil {
				return
			}
			srcAddr, ok := filters["srcAddr"]
			if ok && srcAddr != "" && !utils.ContainsAny(srcAddressList, strings.Split(srcAddr, ",")) {
				continue
			}
			// 处理目的地址
			var destAddressList []string
			destAddressList, err = firewall.DealACLSourceOrDestAddr(policy.DestinationAddress)
			if err != nil {
				return
			}
			dstAddr, ok := filters["dstAddr"]
			if ok && dstAddr != "" && !utils.ContainsAny(destAddressList, strings.Split(dstAddr, ",")) {
				continue
			}
			var servicePortList []string
			servicePortList, err = firewall.DealServicePort(policy.DestinationService)
			servicePort, ok := filters["servicePort"]
			if ok && servicePort != "" && !utils.ContainsAny(servicePortList, strings.Split(servicePort, ",")) {
				continue
			}

			// 策略查询结果
			policyList = append(policyList, dto.Policy{
				ID:           id,
				Name:         aclName,
				Action:       ciscoActionExchange(policy.Permit),
				Enable:       policy.Active,
				SrcZone:      utils.AdditionalPolicyItem([]string{aclName}),
				DestZone:     utils.AdditionalPolicyItem(nil),
				SrcAddress:   utils.AdditionalPolicyItem(srcAddressList),
				DestAddress:  utils.AdditionalPolicyItem(destAddressList),
				ServicePort:  utils.AdditionalPolicyItem(servicePortList),
				FirewallName: firewall.Name,
			})
		}
	}
	return
}

func (firewall *FirewallCiscoV1) CreatePolicy(params dto.CreatePolicyParams) (err error) {
	return
}
