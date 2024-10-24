package h3c_v1

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

type FirewallH3CV1 struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	Protocol    string `json:"protocol"`
	TokenID     string `json:"tokenId"`
	Name        string `json:"name"`
	VirtualZone string `json:"virtualZone"`
}

// createTokenID 获取token
func (firewall *FirewallH3CV1) createTokenID() error {
	// create base url
	url := firewall.createURL("/tokens")
	// create request params
	requestParams := &requests.RequestParams{
		URL:    url,
		Method: http.MethodPost,
		BasicAuth: struct {
			Username string
			Password string
		}{
			Username: firewall.Username,
			Password: firewall.Password,
		},
	}
	// create http client
	client := requests.NewHTTPClient(true)
	// request and receive response
	resp, err := client.Request(requestParams)
	logger.Infof("H3C create token request statusCode: %d, responseBody: %+v", resp.StatusCode, string(resp.Body))
	if err != nil {
		logger.Infof("H3C create token request error: %s", err.Error())
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("H3C create token failed, statusCode: %d, responseBody: %+v", resp.StatusCode, string(resp.Body))
	}
	// unmarshal response
	var response tokenResp
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return fmt.Errorf("H3C create token failed, unmarshal response error: %s", err.Error())
	}
	firewall.TokenID = response.TokenID
	return nil
}

// createURL generates a complete URL based on the given path.
func (firewall *FirewallH3CV1) createURL(path string) string {
	return fmt.Sprintf("%s://%s/api/v1%s", firewall.Protocol, firewall.Address, path)
}

// createRequestParams creates a request object with common headers.
func (firewall *FirewallH3CV1) createRequestParams(tokenID string, method, url string, body interface{}) *requests.RequestParams {
	return &requests.RequestParams{
		URL:     url,
		Method:  method,
		Headers: map[string]string{"X-Auth-Token": tokenID},
		Body:    body,
	}
}

// doRequest sends a request to the Fortinet API.
func (firewall *FirewallH3CV1) doRequest(path, method string, body map[string]interface{}) (*requests.Response, error) {
	if firewall.TokenID == "" {
		err := firewall.createTokenID()
		if err != nil {
			return nil, fmt.Errorf("failed to create tokenID: %s", err.Error())
		}
	}
	// create base url
	url := firewall.createURL(path)
	// create request params
	requestParams := firewall.createRequestParams(firewall.TokenID, method, url, body)
	logger.Infof("H3C requestParams: %+v", requestParams)
	// create http client
	client := requests.NewHTTPClient(false)
	// request and receive response
	resp, err := client.Request(requestParams)
	logger.Infof("H3C request statusCode: %d, responseBody: %s", resp.StatusCode, string(resp.Body))
	return resp, err
}

// getRules
func (firewall *FirewallH3CV1) getRules() (ruleResult rulesResp, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/GetRules", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	ok, err := requests.RequestGetStatusCodeCheck("H3C getRules", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		if err = json.Unmarshal(resp.Body, &ruleResult); err != nil {
			return
		}
	}
	return
}

// getPolicyGroupAddress  地址组名：地址 （dip1：[10.2.3.3,10.2.3.0 255.255.255.0]）
func (firewall *FirewallH3CV1) getPolicyGroupAddress() (respMap map[string][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/OMS/IPv4Objs", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getPolicyGroupAddress", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := policyGroupAddressResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}

		GroupObjectMap := make(map[string][]string)
		respMap = make(map[string][]string)
		for _, v := range response.IPv4Objs {
			switch v.Type {
			case 0:
				GroupObjectMap[v.Group] = append(GroupObjectMap[v.Group], v.NestedGroup)
			case 1:
				respMap[v.Group] = append(respMap[v.Group], fmt.Sprintf("%s/%s", v.SubnetIPv4Address, v.IPv4Mask))
			case 2:
				respMap[v.Group] = append(respMap[v.Group], fmt.Sprintf("%s-%s", v.StartIPv4Address, v.EndIPv4Address))
			case 3:
				respMap[v.Group] = append(respMap[v.Group], v.HostIPv4Address)
			}
		}

		for group, valueList := range GroupObjectMap {
			for _, v := range valueList {
				respMap[group] = append(respMap[group], respMap[v]...)
			}
		}
	}
	return
}

// getPolicyGroupServicePort 获取服务组名：服务/端口 (ser1：[1TCP:12/23-12/12,TCP:0/65535-10000/20000])
func (firewall *FirewallH3CV1) getPolicyGroupServicePort() (resMap map[string][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/OMS/ServObjs", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getPolicyGroupServicePort", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := policyGroupServicePortResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}

		GroupObjectMap := make(map[string][]string)
		resMap = make(map[string][]string)
		for _, v := range response.ServObjs {
			switch v.Type {
			case 0:
				GroupObjectMap[v.Group] = append(GroupObjectMap[v.Group], v.NestedGroup)
			case 3:
				resMap[v.Group] = append(resMap[v.Group], fmt.Sprintf("%s:%d/%d", protocolTCP, v.StartDestPort, v.EndDestPort))
			case 4:
				resMap[v.Group] = append(resMap[v.Group], fmt.Sprintf("%s:%d/%d", protocolUDP, v.StartDestPort, v.EndDestPort))
			}
		}

		for n, vList := range GroupObjectMap {
			for _, v := range vList {
				resMap[n] = append(resMap[n], resMap[v]...)
			}
		}
	}
	return
}

// getSrcPolicySimpleAddrMap 获取策略ID：src_address (1：[77.7.7.0,1.2.3.0/24]）)
func (firewall *FirewallH3CV1) getSrcPolicySimpleAddrMap() (res map[int64][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4SrcSimpleAddr", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getPolicyGroupServicePort", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := srcPolicySimpleAddrResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}

		res = make(map[int64][]string)
		for _, v := range response.IPv4SrcSimpleAddr {
			res[v.ID] = append(res[v.ID], v.SimpleAddrList.SimpleAddrItem...)
		}
	}
	return
}

// getDestPolicySimpleAddrMap 获取策略ID：dest_address (1：[77.7.7.0,1.2.3.0/24]）)
func (firewall *FirewallH3CV1) getDestPolicySimpleAddrMap() (res map[int64][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4DestSimpleAddr", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getDestPolicySimpleAddrMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := destPolicySimpleAddrResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}
		res = make(map[int64][]string)
		for _, v := range response.IPv4DestSimpleAddr {
			res[v.ID] = append(res[v.ID], v.SimpleAddrList.SimpleAddrItem...)
		}
	}
	return
}

// getServPolicyPortMap 获取策略ID：Service/Port
func (firewall *FirewallH3CV1) getServPolicyPortMap() (res map[int64][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4ServObj", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getServPolicyPortMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := servicePortResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}

		res = make(map[int64][]string)
		for _, v := range response.IPv4ServObj {
			for _, serviceObj := range v.ServObjList.ServObjItem {
				var temp string
				switch serviceObj.Type {
				case "0":
					temp = fmt.Sprintf("%s:%s/%s", protocolTCP, serviceObj.StartDestPort, serviceObj.EndDestPort)
				case "1":
					temp = fmt.Sprintf("%s:%s/%s", protocolUDP, serviceObj.StartDestPort, serviceObj.EndDestPort)
				case "2":
					temp = protocolICMP
				}
				res[v.ID] = append(res[v.ID], temp)
			}
		}
	}
	return
}

// getSrcPolicyIDGroupMap 获取策略ID：源地址组
func (firewall *FirewallH3CV1) getSrcPolicyIDGroupMap() (res map[int64][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4SrcAddr", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getSrcPolicyIDGroupMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := srcPolicyIDGroupResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}
		res = make(map[int64][]string)
		for _, v := range response.IPv4SrcAddr {
			res[v.ID] = append(res[v.ID], v.NameList.NameItem...)
		}
	}
	return
}

// getDestPolicyIDGroupMap 获取策略ID：目的地址组
func (firewall *FirewallH3CV1) getDestPolicyIDGroupMap() (res map[int64][]string, err error) {

	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4DestAddr", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getDestPolicyIDGroupMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := destPolicyIDGroupResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}
		res = make(map[int64][]string)
		for _, v := range response.IPv4DestAddr {
			res[v.ID] = append(res[v.ID], v.NameList.NameItem...)
		}
	}
	return
}

// getServicePolicyIDGroupMap 获取策略ID：服务组
func (firewall *FirewallH3CV1) getServicePolicyIDGroupMap() (res map[int64][]string, err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4ServGrp", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}

	ok, err := requests.RequestGetStatusCodeCheck("H3C getServicePolicyIDGroupMap", resp.Body, resp.StatusCode)
	if err != nil {
		return
	}

	if ok {
		response := servicePolicyIDGroupResp{}
		if err = json.Unmarshal(resp.Body, &response); err != nil {
			return
		}
		res = make(map[int64][]string)
		for _, v := range response.IPv4ServGrp {
			res[v.ID] = append(res[v.ID], v.NameList.NameItem...)
		}
	}
	return
}

// createPolicyAddressGroup 创建地址组
func (firewall *FirewallH3CV1) createPolicyAddressGroup(title string) (err error) {
	body := map[string]interface{}{
		"Name": title,
	}
	// send request and receive response
	resp, err := firewall.doRequest("/OMS/IPv4Groups", http.MethodPost, body)
	if err != nil {
		err = fmt.Errorf("H3C createPolicyAddressGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}

	if err = requests.RequestPostStatusCodeCheck("H3C createPolicyAddressGroup", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicyAddressGroup 添加IP到地址组
func (firewall *FirewallH3CV1) addPolicyAddressGroup(groupName, addr string) (err error) {

	body := map[string]interface{}{
		"Group":           groupName,
		"ID":              int64(4294967295), //自动分配ID
		"Type":            3,
		"HostIPv4Address": addr,
		"ObjDescription":  "add address",
	}
	if strings.Contains(addr, "/") {
		addrSplit := strings.Split(addr, "/")
		SubnetIPv4Address := addrSplit[0]
		IPv4Mask, _ := strconv.Atoi(addrSplit[1])
		IPv4Mask_, _ := utils.CIDRToMask(IPv4Mask)
		body = map[string]interface{}{
			"Group":             groupName,
			"ID":                int64(4294967295), //自动分配ID
			"Type":              1,
			"SubnetIPv4Address": SubnetIPv4Address,
			"IPv4Mask":          IPv4Mask_,
			"ObjDescription":    "add address",
		}

	} else if strings.Contains(addr, "-") {
		addrSplit := strings.Split(addr, "-")
		StartIPv4Address := addrSplit[0]
		EndIPv4Address := addrSplit[1]
		body = map[string]interface{}{
			"Group":            groupName,
			"ID":               int64(4294967295), //自动分配ID
			"Type":             2,
			"StartIPv4Address": StartIPv4Address,
			"EndIPv4Address":   EndIPv4Address,
			"ObjDescription":   "add address",
		}

	}
	// send request and receive response
	resp, err := firewall.doRequest("/OMS/IPv4Objs", http.MethodPost, body)
	if err != nil {
		err = fmt.Errorf("H3C addPolicyAddressGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicyAddressGroup", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// createPolicyRules 创建策略规则
func (firewall *FirewallH3CV1) createPolicyRules(policyName string, action int64) (policyID int64, err error) {
	// 获取策略个数 计算ID
	rules, err := firewall.getRules()
	if err != nil {
		return
	}

	policyID = int64(len(rules.GetRules) + 1 + 1000)
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4Rules", http.MethodPost, map[string]interface{}{
		"ID":       policyID,
		"RuleName": policyName,
		"Action":   action,
	})
	if err != nil {
		err = fmt.Errorf("H3C createPolicyRules failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C createPolicyRules", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// createPolicyServiceGroup 创建服务组
func (firewall *FirewallH3CV1) createPolicyServiceGroup(groupName string) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/OMS/ServGroups", http.MethodPost, map[string]interface{}{
		"Name": groupName,
	})
	if err != nil {
		err = fmt.Errorf("H3C createPolicyServiceGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C createPolicyServiceGroup", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicyServiceGroup 添加服务到地址组
func (firewall *FirewallH3CV1) addPolicyServiceGroup(groupName, service string) (err error) {

	t := 4
	serviceSplit := strings.Split(service, ":")
	if serviceSplit[0] == strings.ToLower(protocolTCP) {
		t = 3
	}
	portSplit := strings.Split(serviceSplit[1], "/")
	StartDestPort, _ := strconv.Atoi(portSplit[0])
	EndDestPort, _ := strconv.Atoi(portSplit[1])
	body := map[string]interface{}{
		"Group":          groupName,
		"ID":             int64(4294967295), //自动分配ID
		"Type":           t,
		"StartDestPort":  StartDestPort,
		"EndDestPort":    EndDestPort,
		"ObjDescription": "add address",
	}
	// send request and receive response
	resp, err := firewall.doRequest("/OMS/ServObjs", http.MethodPost, body)
	if err != nil {
		err = fmt.Errorf("H3C addPolicyServiceGroup failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicyServiceGroup", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicyServiceToPolicy 添加服务组到策略
func (firewall *FirewallH3CV1) addPolicyServiceToPolicy(groupName string, policyID int64) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4ServGrp", http.MethodPost, map[string]interface{}{
		"ID":          policyID,
		"SeqNum":      1,
		"IsIncrement": false,
		"NameList": map[string][]string{
			"NameItem": {groupName},
		},
	})
	if err != nil {
		err = fmt.Errorf("H3C addPolicyServiceToPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicyServiceToPolicy", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicySrcGroupToPolicy 添加源地址组到策略
func (firewall *FirewallH3CV1) addPolicySrcGroupToPolicy(groupName string, policyID int64) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4SrcAddr", http.MethodPost, map[string]interface{}{
		"ID":          policyID,
		"SeqNum":      1,
		"IsIncrement": false,
		"NameList": map[string][]string{
			"NameItem": {groupName},
		},
	})
	if err != nil {
		err = fmt.Errorf("H3C addPolicySrcGroupToPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicySrcGroupToPolicy", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicyDestGroupToPolicy 添加目的地址组到策略
func (firewall *FirewallH3CV1) addPolicyDestGroupToPolicy(policyName string, policyID int64) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4DestAddr", http.MethodPost, map[string]interface{}{
		"ID":          policyID,
		"SeqNum":      1,
		"IsIncrement": false,
		"NameList": map[string][]string{
			"NameItem": {policyName},
		},
	})
	if err != nil {
		err = fmt.Errorf("H3C addPolicyDestGroupToPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicyDestGroupToPolicy", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicySrcZoneToPolicy 添加源区域到策略
func (firewall *FirewallH3CV1) addPolicySrcZoneToPolicy(srcName string, policyID int64) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4SrcSecZone", http.MethodPost, map[string]interface{}{
		"ID":          policyID,
		"SeqNum":      1,
		"IsIncrement": false,
		"NameList": map[string][]string{
			"NameItem": {srcName},
		},
	})
	if err != nil {
		err = fmt.Errorf("H3C addPolicySrcZoneToPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicySrcZoneToPolicy", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}

// addPolicyDestZoneToPolicy 添加源区域到策略
func (firewall *FirewallH3CV1) addPolicyDestZoneToPolicy(destName string, policyID int64) (err error) {
	// send request and receive response
	resp, err := firewall.doRequest("/SecurityPolicies/IPv4DestSecZone", http.MethodPost, map[string]interface{}{
		"ID":          policyID,
		"SeqNum":      1,
		"IsIncrement": false,
		"NameList": map[string][]string{
			"NameItem": {destName},
		},
	})
	if err != nil {
		err = fmt.Errorf("H3C addPolicyDestZoneToPolicy failed, status code: %d, error: %v", resp.StatusCode, err)
		return
	}
	if err = requests.RequestPostStatusCodeCheck("H3C addPolicyDestZoneToPolicy", resp.Body, resp.StatusCode, http.StatusCreated); err != nil {
		return
	}
	return
}
