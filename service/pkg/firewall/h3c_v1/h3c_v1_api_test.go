package h3c_v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock HTTP success Server
func setupSuccessMockServer() *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			switch r.URL.Path {
			case "/api/v1/SecurityPolicies/GetRules":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"GetRules": [
							{
								"id": 1,
								"name": "TestRule",
								"action": 1,
								"enable": true,
								"srcZoneList": {
									"srczone1": [
									"srcz1",
									"srcz2"
									]
								},
								"destZoneList": {
									"destzone2": [
									"destz3",
									"destz4"
									]
								},
								"servGrpList": {
									"servgrp1": [
									"servg1",
									"servg2"
									]
								},
								"srcSimpleAddrList": {
									"srcSimpleaddr1": [
									"srcSimplea1",
									"srcSimplea2"
									]
								},
								"destSimpleAddrList": {
									"destSimpleaddr2": [
									"destSimplea3",
									"destSimplea4"
									]
								},
								"servObjList": {
									"servobj1": [
									"servo1",
									"servo2"
									]
								},
								"srcAddrList": {
									"srcaddr3": [
									"srca5",
									"srca6"
									]
								},
								"destAddrList": {
									"destaddr4": [
									"desta7",
									"desta8"
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/OMS/IPv4Objs":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4Objs": [
							{
								"Group": "TestGroup1",
								"Type": 0,
								"StartIPv4Address": "",
								"EndIPv4Address": "",
								"HostIPv4Address": "",
								"SubnetIPv4Address": "",
								"IPv4Mask": "",
								"NestedGroup": "NestedGroup"
							},
							{
								"Group": "TestGroup2",
								"Type": 1,
								"StartIPv4Address": "",
								"EndIPv4Address": "",
								"HostIPv4Address": "",
								"SubnetIPv4Address": "192.168.1.0",
								"IPv4Mask": "32",
								"NestedGroup": ""
							},
							{
								"Group": "TestGroup3",
								"Type": 2,
								"StartIPv4Address": "192.168.1.1",
								"EndIPv4Address": "192.168.1.10",
								"HostIPv4Address": "",
								"SubnetIPv4Address": "",
								"IPv4Mask": "",
								"NestedGroup": ""
							},
							{
								"Group": "TestGroup4",
								"Type": 3,
								"StartIPv4Address": "",
								"EndIPv4Address": "",
								"HostIPv4Address": "192.168.1.10",
								"SubnetIPv4Address": "",
								"IPv4Mask": "",
								"NestedGroup": ""
							}
						]
					}`),
				)
			case "/api/v1/OMS/ServObjs":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"ServObjs": [
							{
								"Group": "TestGroup1",
								"Type": 0,
								"StartSrcPort": 0,
								"StartDestPort": 0,
								"EndSrcPort": 0,
								"EndDestPort": 0,
								"NestedGroup": "NestedGroup1"
							},
							{
								"Group": "TestGroup2",
								"Type": 3,
								"StartSrcPort": 80,
								"StartDestPort": 80,
								"EndSrcPort": 65535,
								"EndDestPort": 65535,
								"NestedGroup": ""
							},
							{
								"Group": "TestGroup2",
								"Type": 3,
								"StartSrcPort": 443,
								"StartDestPort": 443,
								"EndSrcPort": 65535,
								"EndDestPort": 65535,
								"NestedGroup": ""
							},
							{
								"Group": "TestGroup3",
								"Type": 4,
								"StartSrcPort": 11024,
								"StartDestPort": 11024,
								"EndSrcPort": 65535,
								"EndDestPort": 65535,
								"NestedGroup": ""
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4SrcSimpleAddr":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4SrcSimpleAddr": [
							{
								"ID": 1,
								"SimpleAddrList": {
									"SimpleAddrItem": [
										"192.168.1.1",
										"192.168.1.2"
									]
								}
							},
							{
								"ID": 2,
								"SimpleAddrList": {
									"SimpleAddrItem": [
										"192.168.2.1",
										"192.168.2.2"
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4DestSimpleAddr":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4DestSimpleAddr": [
							{
								"ID": 1,
								"SimpleAddrList": {
									"SimpleAddrItem": [
										"192.168.1.1",
										"192.168.1.2"
									]
								}
							},
							{
								"ID": 2,
								"SimpleAddrList": {
									"SimpleAddrItem": [
										"192.168.2.1",
										"192.168.2.2"
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4ServObj":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4ServObj": [
							{
								"ID": 1,
								"ServObjList": {
									"ServObjItem": [
										{
											"Group": "TestGroup1",
											"Type": "0",
											"StartSrcPort": "1024",
											"StartDestPort": "80",
											"EndSrcPort": "65535",
											"EndDestPort": "65535",
											"NestedGroup": ""
										},
										{
											"Group": "TestGroup2",
											"Type": "1",
											"StartSrcPort": "1024",
											"StartDestPort": "11024",
											"EndSrcPort": "65535",
											"EndDestPort": "65535",
											"NestedGroup": ""
										},
										{
											"Group": "TestGroup3",
											"Type": "2",
											"StartSrcPort": "",
											"StartDestPort": "",
											"EndSrcPort": "",
											"EndDestPort": "",
											"NestedGroup": ""
										}
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4SrcAddr":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4SrcAddr": [
							{
								"ID": 1,
								"NameList": {
									"NameItem": [
										"TestGroup1",
										"TestGroup2"
									]
								}
							},
							{
								"ID": 2,
								"NameList": {
									"NameItem": [
										"TestGroup3",
										"TestGroup2"
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4DestAddr":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4DestAddr": [
							{
								"ID": 1,
								"NameList": {
									"NameItem": [
										"TestGroup1",
										"TestGroup2"
									]
								}
							},
							{
								"ID": 2,
								"NameList": {
									"NameItem": [
										"TestGroup3",
										"TestGroup2"
									]
								}
							}
						]
					}`),
				)
			case "/api/v1/SecurityPolicies/IPv4ServGrp":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"IPv4ServGrp": [
							{
								"ID": 1,
								"NameList": {
									"NameItem": [
										"TestGroup1",
										"TestGroup2"
									]
								}
							},
							{
								"ID": 2,
								"NameList": {
									"NameItem": [
										"TestGroup3",
										"TestGroup2"
									]
								}
							}
						]
					}`),
				)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case "POST":
			switch r.URL.Path {
			case "/api/v1/tokens":
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`
					{
						"token-id": "Token Test",
						"link": "https://example.com/api/v1/tokens",
						"expiry-time": "1759907493000"
					}`),
				)
			case "/api/v1/OMS/IPv4Groups", "/api/v1/OMS/IPv4Objs", "/api/v1/OMS/ServGroups", "/api/v1/OMS/ServObjs", "/api/v1/SecurityPolicies/IPv4Rules", "/api/v1/SecurityPolicies/IPv4ServGrp", "/api/v1/SecurityPolicies/IPv4SrcAddr", "/api/v1/SecurityPolicies/IPv4DestAddr", "/api/v1/SecurityPolicies/IPv4SrcSecZone", "/api/v1/SecurityPolicies/IPv4DestSecZone":
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`
					{
						"http_method": "POST",
						"http_status": 201,
						"status": "success",
					}`),
				)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	return mockServer
}

// Mock HTTP failed Server
func setupFailedMockServer() *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"http_method":"GET","status":"error","http_status":403}"}`))
		} else if r.Method == "POST" {
			if r.URL.Path == "/api/v1/tokens" {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`
					{
						"token-id": "Token Test",
						"link": "https://example.com/api/v1/tokens",
						"expiry-time": "1759907493000"
					}`),
				)
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"http_method":"POST","status":"error","http_status":403}"}`))
			}
		}
	}))
	return mockServer
}

func TestCreateTokenID(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	err := firewall.createTokenID()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, "Token Test", firewall.TokenID)
}

func TestCreateTokenIDWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	err := firewall.createTokenID()

	// Assert expectations
	assert.Error(t, err)
}

func TestDoRequest(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	path := "/not-found"
	methodGet := http.MethodGet
	methodPost := http.MethodPost
	body := map[string]interface{}{}

	// Call the method to test
	respGet, errGet := firewall.doRequest(path, methodGet, body)
	respPost, errPost := firewall.doRequest(path, methodPost, body)

	// Assert expectations
	assert.Nil(t, errGet)
	assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	assert.Nil(t, errPost)
	assert.Equal(t, http.StatusNotFound, respPost.StatusCode)
}

func TestDoRequestErr(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	path := "/not-found"
	methodGet := http.MethodGet
	methodPost := http.MethodPost
	body := map[string]interface{}{}

	// Call the method to test
	respGet, errGet := firewall.doRequest(path, methodGet, body)
	respPost, errPost := firewall.doRequest(path, methodPost, body)

	// Assert expectations
	assert.Nil(t, errGet)
	assert.Equal(t, http.StatusNotFound, respGet.StatusCode)
	assert.Nil(t, errPost)
	assert.Equal(t, http.StatusNotFound, respPost.StatusCode)
}

func TestGetRules(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	ruleResult, err := firewall.getRules()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, "TestRule", ruleResult.GetRules[0].Name)
}

func TestGetRulesWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	ruleResult, err := firewall.getRules()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(ruleResult.GetRules))
}

func TestGetPolicyGroupAddress(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getPolicyGroupAddress()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string(nil), respMap["TestGroup1"])
	assert.Equal(t, []string{"192.168.1.0/32"}, respMap["TestGroup2"])
	assert.Equal(t, []string{"192.168.1.1-192.168.1.10"}, respMap["TestGroup3"])
	assert.Equal(t, []string{"192.168.1.10"}, respMap["TestGroup4"])
}

func TestGetPolicyGroupAddressWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getPolicyGroupAddress()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetPolicyGroupServicePort(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getPolicyGroupServicePort()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string(nil), respMap["TestGroup1"])
	assert.Equal(t, []string{"TCP:80/65535", "TCP:443/65535"}, respMap["TestGroup2"])
	assert.Equal(t, []string{"UDP:11024/65535"}, respMap["TestGroup3"])
}

func TestGetPolicyGroupServicePortWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getPolicyGroupServicePort()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetSrcPolicySimpleAddrMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getSrcPolicySimpleAddrMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"192.168.1.1", "192.168.1.2"}, respMap[1])
	assert.Equal(t, []string{"192.168.2.1", "192.168.2.2"}, respMap[2])
}

func TestGetSrcPolicySimpleAddrMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getSrcPolicySimpleAddrMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetDestPolicySimpleAddrMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getDestPolicySimpleAddrMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"192.168.1.1", "192.168.1.2"}, respMap[1])
	assert.Equal(t, []string{"192.168.2.1", "192.168.2.2"}, respMap[2])
}

func TestGetDestPolicySimpleAddrMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getDestPolicySimpleAddrMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetServPolicyPortMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getServPolicyPortMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"TCP:80/65535", "UDP:11024/65535", "ICMP"}, respMap[1])
}

func TestGetServPolicyPortMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getServPolicyPortMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetSrcPolicyIDGroupMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getSrcPolicyIDGroupMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"TestGroup1", "TestGroup2"}, respMap[1])
	assert.Equal(t, []string{"TestGroup3", "TestGroup2"}, respMap[2])
}

func TestGetSrcPolicyIDGroupMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getSrcPolicyIDGroupMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetDestPolicyIDGroupMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getDestPolicyIDGroupMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"TestGroup1", "TestGroup2"}, respMap[1])
	assert.Equal(t, []string{"TestGroup3", "TestGroup2"}, respMap[2])
}

func TestGetDestPolicyIDGroupMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getDestPolicyIDGroupMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestGetServicePolicyIDGroupMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getServicePolicyIDGroupMap()

	// Assert expectations
	assert.Nil(t, err)
	assert.Equal(t, []string{"TestGroup1", "TestGroup2"}, respMap[1])
	assert.Equal(t, []string{"TestGroup3", "TestGroup2"}, respMap[2])
}

func TestTestGetServicePolicyIDGroupMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()
	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	// Call the method to test
	respMap, err := firewall.getServicePolicyIDGroupMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(respMap))
}

func TestCreatePolicyAddressGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	title := "TestTitle"

	// Call the method to test
	err := firewall.createPolicyAddressGroup(title)

	// Assert expectations
	assert.NoError(t, err)
}

func TestCreatePolicyAddressGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	title := "TestTitle"

	// Call the method to test
	err := firewall.createPolicyAddressGroup(title)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyAddressGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	addr1 := "192.168.1.0/32"
	addr2 := "192.168.1.1-192.168.1.10"

	// Call the method to test
	err1 := firewall.addPolicyAddressGroup(groupName, addr1)
	err2 := firewall.addPolicyAddressGroup(groupName, addr2)

	// Assert expectations
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestAddPolicyAddressGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	addr1 := "192.168.1.0/32"
	addr2 := "192.168.1.1-192.168.1.10"

	// Call the method to test
	err1 := firewall.addPolicyAddressGroup(groupName, addr1)
	err2 := firewall.addPolicyAddressGroup(groupName, addr2)

	// Assert expectations
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func TestCreatePolicyRules(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	policyName := "TestPolicy"
	var action int64 = 1

	// Call the method to test
	policyID, err := firewall.createPolicyRules(policyName, action)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, int64(1002), policyID)
}

func TestCreatePolicyRulesWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	policyName := "TestPolicy"
	var action int64 = 1

	// Call the method to test
	_, err := firewall.createPolicyRules(policyName, action)

	// Assert expectations
	assert.Error(t, err)
}

func TestCreatePolicyServiceGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"

	// Call the method to test
	err := firewall.createPolicyServiceGroup(groupName)

	// Assert expectations
	assert.NoError(t, err)
}

func TestCreatePolicyServiceGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"

	// Call the method to test
	err := firewall.createPolicyServiceGroup(groupName)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyServiceGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	service1 := "TCP:80/65535"
	service2 := "UDP:11024/65535"

	// Call the method to test
	err1 := firewall.addPolicyServiceGroup(groupName, service1)
	err2 := firewall.addPolicyServiceGroup(groupName, service2)

	// Assert expectations
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestAddPolicyServiceGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	service1 := "TCP:80/65535"
	service2 := "UDP:11024/65535"

	// Call the method to test
	err1 := firewall.addPolicyServiceGroup(groupName, service1)
	err2 := firewall.addPolicyServiceGroup(groupName, service2)

	// Assert expectations
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func TestAddPolicyServiceToPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyServiceToPolicy(groupName, policyID)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyServiceToPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyServiceToPolicy(groupName, policyID)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicySrcGroupToPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicySrcGroupToPolicy(groupName, policyID)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicySrcGroupToPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicySrcGroupToPolicy(groupName, policyID)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyDestGroupToPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyDestGroupToPolicy(groupName, policyID)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyDestGroupToPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyDestGroupToPolicy(groupName, policyID)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicySrcZoneToPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicySrcZoneToPolicy(groupName, policyID)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicySrcZoneToPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicySrcZoneToPolicy(groupName, policyID)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyDestZoneToPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyDestZoneToPolicy(groupName, policyID)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyDestZoneToPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallH3CV1
	firewall := FirewallH3CV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		Username:    "username",
		Password:    "password",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	var policyID int64 = 1002

	// Call the method to test
	err := firewall.addPolicyDestZoneToPolicy(groupName, policyID)

	// Assert expectations
	assert.Error(t, err)
}
