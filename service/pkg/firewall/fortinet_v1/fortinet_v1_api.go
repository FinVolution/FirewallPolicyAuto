package fortinet_v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/requests"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/logger"
	"github.com/FinVolution/FirewallPolicyAuto/service/utils"
)

type FirewallFortinetV1 struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Protocol    string `json:"protocol"`
	TokenID     string `json:"tokenID"`
	VirtualZone string `json:"virtualZone"`
}

// createURL generates a complete URL based on the given path.
func (firewall *FirewallFortinetV1) createURL(path string) string {
	return fmt.Sprintf("%s://%s/api/v2/cmdb%s", firewall.Protocol, firewall.Address, path)
}

// createRequestParams creates a request object with common headers.
func (firewall *FirewallFortinetV1) createRequestParams(tokenID string, method, url string, body interface{}) *requests.RequestParams {
	requestParams := &requests.RequestParams{
		URL:     url,
		Method:  method,
		Headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", tokenID)},
		Body:    body,
	}
	if firewall.VirtualZone != "" {
		requestParams.QueryParams = map[string]string{
			"vdom": firewall.VirtualZone,
		}
	}
	return requestParams
}

// doRequest sends a request to the Fortinet API.
func (firewall *FirewallFortinetV1) doRequest(path, method string, body map[string]interface{}) (*requests.Response, error) {
	if firewall.TokenID == "" {
		return nil, fmt.Errorf("FirewallFortinetV1 tokenID can not be empty!!")
	}
	// create base url
	url := firewall.createURL(path)
	// create request params
	requestParams := firewall.createRequestParams(firewall.TokenID, method, url, body)
	logger.Infof("Fortinet requestParams: %+v", requestParams)
	// create http client
	client := requests.NewHTTPClient(true)
	// request and receive response
	resp, err := client.Request(requestParams)
	logger.Infof("Fortinet request statusCode: %d, responseBody: %s", resp.StatusCode, string(resp.Body))
	return resp, err
}

// getPolicy
func (firewall *FirewallFortinetV1) getPolicy() (policyResult policyResp, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/firewall/policy/", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("Fortinet getPolicy", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		if err = json.Unmarshal(resp.Body, &policyResult); err != nil {
			return
		}
	}
	return
}

// getPolicyNameAddressMap 获取信息  地址名称：IP地址（10.2.3.0/20：[10.2.3.0 255.255.255.0]）
func (firewall *FirewallFortinetV1) getPolicyNameAddressMap() (resMap map[string][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/firewall/address/", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("Fortinet getPolicyNameAddressMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		result := vdomResp{}
		if err = json.Unmarshal(resp.Body, &result); err != nil {
			return
		}
		resMap = make(map[string][]string)
		for _, v := range result.Results {
			switch v.Type {
			case "ipmask":
				resMap[v.Name] = []string{v.Subnet}
			case "iprange":
				resMap[v.Name] = []string{fmt.Sprintf("%s-%s", v.StartIp, v.EndIp)}
			}
		}
	}
	return
}

// GetPolicyGroupNameAddressMap 获取信息  地址组名称：地址名称（test_group：[10.114.1.1/32, 10.114.1.2/32]）
func (firewall *FirewallFortinetV1) getPolicyGroupNameAddressMap() (resMap map[string][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/firewall/addrgrp/", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("Fortinet getPolicyNameAddressMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		result := groupResp{}
		if err = json.Unmarshal(resp.Body, &result); err != nil {
			return
		}

		resMap = make(map[string][]string)
		for _, v := range result.Results {
			for _, vv := range v.Member {
				resMap[v.Name] = append(resMap[v.Name], vv.Name)
			}
		}
	}
	return
}

// getPolicyNameServicePortMap 获取信息  服务名称：协议/端口 （test_111：[tcp：111]）
func (firewall *FirewallFortinetV1) getPolicyNameServicePortMap() (resMap map[string][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/firewall.service/custom/", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("Fortinet getPolicyNameAddressMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {

		result := serviceResp{}
		if err = json.Unmarshal(resp.Body, &result); err != nil {
			return
		}

		resMap = make(map[string][]string)
		for _, v := range result.Results {
			if v.Protocol == protocolTCPUDP {
				if v.TcpPortRange != "" {
					resMap[v.Name] = []string{fmt.Sprintf("%s:%s", protocolTCP, v.TcpPortRange)}
				} else if v.UdpPortRange != "" {
					resMap[v.Name] = []string{fmt.Sprintf("%s:%s", protocolUDP, v.UdpPortRange)}
				} else if v.SctpPortRange != "" {
					resMap[v.Name] = []string{fmt.Sprintf("%s:%s", protocolSTCP, v.SctpPortRange)}
				}
			} else if v.Protocol == protocolIP {
				resMap[v.Name] = []string{fmt.Sprintf("%s:%s", protocolIp, v.IpRange)}
			}
		}
	}
	return
}

// getPolicyGroupNameServiceMap 获取信息  服务组名称：服务名称（test_group1:[test_111,test_222]）
func (firewall *FirewallFortinetV1) getPolicyGroupNameServiceMap() (resMap map[string][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/firewall.service/group/", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("Fortinet getPolicyNameAddressMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		result := groupResp{}
		if err = json.Unmarshal(resp.Body, &result); err != nil {
			return
		}

		resMap = make(map[string][]string)
		for _, v := range result.Results {
			for _, vv := range v.Member {
				resMap[v.Name] = append(resMap[v.Name], vv.Name)
			}
		}
	}
	return
}

// createPolicyAddress 创建IP地址
func (firewall *FirewallFortinetV1) createPolicyAddress(addr string) (err error) {

	requestBody := map[string]interface{}{
		"end-ip":   "255.255.255.255",
		"start-ip": addr,
		"name":     addr,
		"type":     "ipmask",
	}
	if strings.Contains(addr, "/") {
		addrSplit := strings.Split(addr, "/")
		subnetIPv4Address := addrSplit[0]
		ipv4MaskInt, err := strconv.Atoi(addrSplit[1])
		if err != nil {
			err = fmt.Errorf("createPolicyAddress failed, ipv4Mask int change error: %v", err)
			return err
		}
		ipv4MaskStr, err := utils.CIDRToMask(ipv4MaskInt)
		if err != nil {
			err = fmt.Errorf("createPolicyAddress failed, ipv4Mask str change error: %v", err)
			return err
		}
		requestBody = map[string]interface{}{
			"end-ip":   ipv4MaskStr,
			"start-ip": subnetIPv4Address,
			"name":     addr,
			"type":     "ipmask",
		}
	} else if strings.Contains(addr, "-") {
		addrSplit := strings.Split(addr, "-")
		requestBody = map[string]interface{}{
			"end-ip":   addrSplit[1],
			"start-ip": addrSplit[0],
			"name":     addr,
			"type":     "iprange",
		}
	}
	// send request and receive response
	resp, err := firewall.doRequest("/firewall/address/", http.MethodPost, requestBody)
	if err != nil {
		err = fmt.Errorf("Fortinet createPolicyAddress failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("Fortinet createPolicyAddress", resp.Body, resp.StatusCode, http.StatusOK); err != nil {
		return
	}
	return
}

// addPolicyAddrToGroup 添加源地址到地址组
func (firewall *FirewallFortinetV1) addPolicyAddrToGroup(groupName string, addrList []string) (err error) {
	member := []map[string]interface{}{}
	for _, addr := range addrList {
		member = append(member, map[string]interface{}{
			"name": addr,
		})
	}
	requestBody := map[string]interface{}{
		"name":   groupName,
		"member": member,
	}
	// send request and receive response
	resp, err := firewall.doRequest("/firewall/addrgrp/", http.MethodPost, requestBody)
	if err != nil {
		err = fmt.Errorf("Fortinet addPolicyAddrToGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("Fortinet addPolicyAddrToGroup", resp.Body, resp.StatusCode, http.StatusOK); err != nil {
		return
	}
	return
}

// createPolicyService 创建服务
func (firewall *FirewallFortinetV1) createPolicyService(service string) (err error) {

	var tcpPortRange string
	var udpPortRange string
	var protocol string
	var protocolNumber int
	serviceSplit := strings.Split(service, ":")
	switch serviceSplit[0] {
	case protocolTCP:
		protocol = protocolTCPUDP
		tcpPortRange = strings.ReplaceAll(serviceSplit[1], "/", "-")
	case protocolUDP:
		protocol = protocolTCPUDP
		udpPortRange = strings.ReplaceAll(serviceSplit[1], "/", "-")
	default:
		protocol = protocolIP
		protocolNumber = 0
	}
	requestBody := map[string]interface{}{
		"name":           service,
		"protocol":       protocol,
		"tcp-portrange":  tcpPortRange,
		"udp-portrange":  udpPortRange,
		"protocolNumber": protocolNumber,
	}
	// send request and receive response
	resp, err := firewall.doRequest("/firewall.service/custom/", http.MethodPost, requestBody)
	if err != nil {
		err = fmt.Errorf("Fortinet createPolicyService failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("Fortinet createPolicyService", resp.Body, resp.StatusCode, http.StatusOK); err != nil {
		return
	}
	return
}

// addPolicyServiceToGroup 添加服务到服务组
func (firewall *FirewallFortinetV1) addPolicyServiceToGroup(groupName string, serviceList []string) (err error) {
	member := []map[string]interface{}{}
	for _, service := range serviceList {
		member = append(member, map[string]interface{}{
			"name": service,
		})
	}
	requestBody := map[string]interface{}{
		"name":   groupName,
		"member": member,
	}

	// send request and receive response
	resp, err := firewall.doRequest("/firewall.service/group/", http.MethodPost, requestBody)
	if err != nil {
		err = fmt.Errorf("Fortinet addPolicyServiceToGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("Fortinet addPolicyServiceToGroup", resp.Body, resp.StatusCode, http.StatusOK); err != nil {
		return
	}
	return
}

// addPolicy 创建策略
func (firewall *FirewallFortinetV1) addPolicy(policyName, action, srcZone, destZone, srcAddrGroupName, destAddrGroupName, serviceGroupName string) (err error) {
	requestBody := map[string]interface{}{
		"policyid": 0,
		"action":   action,
		"name":     policyName,
		"srcintf": []map[string]string{
			{
				"name": srcZone,
			},
		},
		"dstintf": []map[string]string{
			{
				"name": destZone,
			},
		},
		"srcaddr": []map[string]string{
			{
				"name": srcAddrGroupName,
			},
		},
		"dstaddr": []map[string]string{
			{
				"name": destAddrGroupName,
			},
		},
		"service": []map[string]string{
			{
				"name": serviceGroupName,
			},
		},
	}

	// send request and receive response
	resp, err := firewall.doRequest("/firewall/policy/", http.MethodPost, requestBody)
	if err != nil {
		err = fmt.Errorf("Fortinet addPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
	}
	if err = requests.RequestPostStatusCodeCheck("Fortinet addPolicy", resp.Body, resp.StatusCode, http.StatusOK); err != nil {
		return
	}
	return
}
