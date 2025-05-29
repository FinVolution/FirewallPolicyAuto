package cisco_v1

import (
	"encoding/json"
	"fmt"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/requests"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/logger"
	"net/http"
	"strings"
)

// createURL generates a complete URL based on the given path.
func (firewall *FirewallCiscoV1) createURL(path string) string {
	return fmt.Sprintf("%s://%s/api/objects%s", firewall.Protocol, firewall.Address, path)
}

// createRequestParams creates a request object with common headers.
func (firewall *FirewallCiscoV1) createRequestParams(tokenID string, method, url string, body interface{}) *requests.RequestParams {
	requestParams := &requests.RequestParams{
		URL:    url,
		Method: method,
		BasicAuth: struct {
			Username string
			Password string
		}{
			Username: firewall.Username,
			Password: firewall.Password,
		},
		Body: body,
	}

	return requestParams
}

func (firewall *FirewallCiscoV1) doRequest(path, method string, body map[string]interface{}) (*requests.Response, error) {
	// create base url
	url := firewall.createURL(path)
	// create request params
	requestParams := firewall.createRequestParams(firewall.TokenID, method, url, body)
	logger.Infof("Fortinet requestParams: %+v", requestParams)
	// create http client
	client := requests.NewHTTPClient(true)
	// request and receive response
	resp, err := client.Request(requestParams)
	if err != nil {
		logger.Infof("Fortinet request error: %s", err.Error())
		return nil, err
	}
	logger.Infof("Fortinet request statusCode: %d, responseBody: %s", resp.StatusCode, string(resp.Body))
	return resp, err
}

func (firewall *FirewallCiscoV1) GetTokenID() (tokenID string, err error) {
	return
}

// DealACLSourceOrDestAddr 接收地址对象，返回ip地址列表和错误信息
func (firewall *FirewallCiscoV1) DealACLSourceOrDestAddr(addrObject_ interface{}) ([]string, error) {
	var addrObjectList []map[string]interface{}
	switch addrObject_.(type) {
	case []interface{}:
		for _, v1 := range addrObject_.([]interface{}) {
			addrObjectList = append(addrObjectList, v1.(map[string]interface{}))
		}
	case map[string]interface{}:
		addrObjectList = append(addrObjectList, addrObject_.(map[string]interface{}))
	default:
		return nil, fmt.Errorf("unexpected type for addrObject_: %T", addrObject_)
	}

	var addrList []string
	for _, addrObject := range addrObjectList {
		kindType := addrObject["kind"].(string)
		switch kindType {
		case "objectRef#NetworkObj":
			addr, err := firewall.GetNetworkObjectAddress(addrObject["objectId"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to get object address: %v", err)
			}
			addrList = append(addrList, addr)
		case "objectRef#NetworkObjGroup":
			addrs, err := firewall.GetNetWorkObjectGroupAddress(addrObject["objectId"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to get object group address: %v", err)
			}
			addrList = append(addrList, addrs...)
		default:
			addrList = append(addrList, addrObject["value"].(string))
		}
	}
	return addrList, nil
}

// DealServicePort 接收服务对象，返回服务端口列表和错误信息
func (firewall *FirewallCiscoV1) DealServicePort(serviceObject interface{}) (serviceList []string, err error) {
	var serviceObjectList []map[string]interface{}
	switch serviceObject.(type) {
	case []interface{}:
		for _, service := range serviceObject.([]interface{}) {
			serviceObjectList = append(serviceObjectList, service.(map[string]interface{}))
		}
	case map[string]interface{}:
		serviceObjectList = append(serviceObjectList, serviceObject.(map[string]interface{}))
	default:
		return nil, fmt.Errorf("unexpected type for serviceObject: %T", serviceObject)
	}

	for _, item := range serviceObjectList {
		kindType := item["kind"].(string)
		switch kindType {
		case "objectRef#NetworkServiceGroup":
			var serviceListItem []string
			serviceListItem, err = firewall.GetNetworkServiceGroupsPort(item["objectId"].(string))
			if err != nil {
				return
			}
			serviceList = append(serviceList, serviceListItem...)
		case "objectRef#NetworkServiceObject":
			var service string
			service, err = firewall.GetNetworkServicePort(item["objectId"].(string))
			if err != nil {
				return
			}
			serviceList = append(serviceList, service)
		default:
			serviceList = append(serviceList, item["value"].(string))
		}
	}
	return
}

// GetFirewallAclItems 获取防火墙所有的ACL
func (firewall *FirewallCiscoV1) GetFirewallAclItems(aclName string, offset int) (items []ciscoPolicy, err error) {
	for {
		// 构造返回结构体
		var ciscoPolicyResponse struct {
			RangeInfo struct {
				Offset int `json:"offset"`
				Limit  int `json:"limit"`
				Total  int `json:"total"`
			} `json:"rangeInfo"`
			Items []ciscoPolicy `json:"items"`
		}
		// 执行请求
		resp, err := firewall.doRequest(fmt.Sprintf("/extendedacls/%s/aces?offset=%d", aclName, offset), http.MethodGet, map[string]interface{}{})
		if err != nil {
			return
		}
		err = json.Unmarshal(resp.Body, &ciscoPolicyResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		// 累积返回的 items
		items = append(items, ciscoPolicyResponse.Items...)
		// 判断是否需要继续
		if offset+ciscoPolicyResponse.RangeInfo.Limit > ciscoPolicyResponse.RangeInfo.Total {
			break
		}
		// 更新 offset
		offset += ciscoPolicyResponse.RangeInfo.Limit

	}

	return
}

// GetNetworkServiceGroupsPort 查询网络服务组所有端口
func (firewall *FirewallCiscoV1) GetNetworkServiceGroupsPort(objectID string) (serviceList []string, err error) {
	var response struct {
		ItemResp
		Members []struct {
			Kind     string `json:"kind"`
			Value    string `json:"value"`
			ObjectId string `json:"objectId"`
		} `json:"members"`
	}
	resp, err := firewall.doRequest(fmt.Sprintf("/networkservicegroups/%s", objectID), http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return
	}
	for _, item := range response.Members {
		if strings.Contains(item.Kind, "NetworkObj") {
			var addr string
			addr, err = firewall.GetNetworkObjectAddress(item.Value)
			if err != nil {
				return
			}
			serviceList = append(serviceList, addr)
		} else if strings.Contains(item.Kind, "NetworkServiceGroups") {
			var addr_List []string
			addr_List, err = firewall.GetNetworkServiceGroupsPort(item.ObjectId)
			if err != nil {
				return
			}
			serviceList = append(serviceList, addr_List...)
		} else {
			serviceList = append(serviceList, item.Value)

		}
	}
	return

}

// GetNetworkServicePort 获取网络服务所有端口
func (firewall *FirewallCiscoV1) GetNetworkServicePort(objectId string) (res string, err error) {
	var response struct {
		ItemResp
	}

	resp, err := firewall.doRequest(fmt.Sprintf("/networkservices/%s", objectId), http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resp.Body), &response)
	if err != nil {
		return
	}
	res = response.Value
	return

}

// GetNetworkObjectAddress 查询某个网络对象包含的所有IP地址， 返回一个IP地址字符串
func (firewall *FirewallCiscoV1) GetNetworkObjectAddress(objectId string) (res string, err error) {
	var response struct {
		ItemResp
		Host struct {
			Kind  string `json:"kind"`
			Value string `json:"value"`
		} `json:"host"`
	}
	resp, err := firewall.doRequest(fmt.Sprintf("/networkobjects/%s", objectId), http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return
	}
	res = response.Host.Value
	return
}

// GetNetWorkObjectGroupAddress 查询某个网络对象组包含的所有IP地址， 返回一个IP地址数组
func (firewall *FirewallCiscoV1) GetNetWorkObjectGroupAddress(objectName string) (addressList []string, err error) {
	var response struct {
		ItemResp
		Members []struct {
			Kind     string `json:"kind"`
			Value    string `json:"value"`
			ObjectId string `json:"objectId"`
		} `json:"members"`
	}
	resp, err := firewall.doRequest(fmt.Sprintf("/networkobjectgroups/%s", objectName), http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return
	}
	for _, item := range response.Members {
		if item.Kind == "objectRef#NetworkObj" {
			var addr string
			addr, err = firewall.GetNetworkObjectAddress(item.ObjectId)
			if err != nil {
				return
			}
			addressList = append(addressList, addr)
		} else if item.Kind == "objectRef#NetworkObjGroup" {
			var addrList []string
			addrList, err = firewall.GetNetWorkObjectGroupAddress(item.ObjectId)
			if err != nil {
				return
			}
			addressList = append(addressList, addrList...)
		} else {
			addressList = append(addressList, item.Value)
		}
	}
	return

}

// GetFirewallNetworkObjectGroupAll 查询所有的网络对象组，返回一个数组
func (firewall *FirewallCiscoV1) GetFirewallNetworkObjectGroupAll(offset int) (items []NetworkObjGroupStruct, err error) {
	for {
		// 构造返回结构体
		var ResponseStruct struct {
			Kind      string `json:"kind"`
			RangeInfo struct {
				Offset int `json:"offset"`
				Limit  int `json:"limit"`
				Total  int `json:"total"`
			} `json:"rangeInfo"`
			Items []NetworkObjGroupStruct `json:"items"`
		}

		// 执行请求
		resp, err := firewall.doRequest(fmt.Sprintf("/networkobjectgroups/?offset=%d", offset), http.MethodGet, map[string]interface{}{})
		if err != nil {
			return nil, fmt.Errorf("failed to request URL %s: %w", fmt.Sprintf("/networkobjectgroups/?offset=%d", offset), err)
		}

		// 解析 JSON 响应
		err = json.Unmarshal(resp.Body, &ResponseStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		// 累积返回的 items
		items = append(items, ResponseStruct.Items...)

		// 判断是否需要继续
		if offset+ResponseStruct.RangeInfo.Limit > ResponseStruct.RangeInfo.Total {
			break
		}
		// 更新 offset
		offset += ResponseStruct.RangeInfo.Limit
	}
	return
}

// GetFirewallNetworkObjectAll 查询所有的网络对象，返回一个数组
func (firewall *FirewallCiscoV1) GetFirewallNetworkObjectAll(offset int) (items []NetworkObjStruct, err error) {
	for {
		// 构造返回结构体
		var ResponseStruct struct {
			Kind      string `json:"kind"`
			RangeInfo struct {
				Offset int `json:"offset"`
				Limit  int `json:"limit"`
				Total  int `json:"total"`
			} `json:"rangeInfo"`
			Items []NetworkObjStruct `json:"items"`
		}

		// 执行请求
		resp, err := firewall.doRequest(fmt.Sprintf("/networkobjects/?offset=%d", offset), http.MethodGet, map[string]interface{}{})
		if err != nil {
			return nil, fmt.Errorf("failed to request URL %s: %w", fmt.Sprintf("/networkobjects/?offset=%d", offset), err)
		}

		// 解析 JSON 响应
		err = json.Unmarshal(resp.Body, &ResponseStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		// 累积返回的 items
		items = append(items, ResponseStruct.Items...)

		// 判断是否需要继续
		if offset+ResponseStruct.RangeInfo.Limit > ResponseStruct.RangeInfo.Total {
			break
		}
		// 更新 offset
		offset += ResponseStruct.RangeInfo.Limit
	}
	return
}

// GetPolicyACLName 获取策略区域列表，返回策略区域名称列表
func (firewall *FirewallCiscoV1) GetPolicyACLName() (res []string, err error) {
	resp, err := firewall.doRequest("/extendedacls", http.MethodGet, map[string]interface{}{})
	if err != nil {
		return
	}
	var response struct {
		Items []struct {
			Name     string `json:"name"`
			ObjectId string `json:"objectId"`
		} `json:"items"`
	}
	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return
	}
	for _, item := range response.Items {
		res = append(res, item.Name)
	}
	return
}
