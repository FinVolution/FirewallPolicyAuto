package fortinet_v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoTokenIdErr(t *testing.T) {
	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     "127.0.0.1",
		Protocol:    "http",
		TokenID:     "",
		VirtualZone: "",
	}

	// Call the method to test
	_, err := firewall.getPolicy()

	// Assert expectations
	assert.Equal(t, "FirewallFortinetV1 tokenID can not be empty!!", err.Error())
}

func TestDoRequestErr(t *testing.T) {
	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     "127.0.0.1",
		Protocol:    "http",
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	path := "/not-found"
	methodGet := http.MethodGet
	methodPost := http.MethodPost
	body := map[string]interface{}{}

	// Call the method to test
	_, errGet := firewall.doRequest(path, methodGet, body)
	_, errPost := firewall.doRequest(path, methodPost, body)

	// Assert expectations
	assert.Error(t, errGet)
	assert.Error(t, errPost)
}

// Mock HTTP success Server
func setupSuccessMockServer() *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			switch r.URL.Path {
			case "/api/v2/cmdb/firewall/policy/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"http_method": "GET",
						"http_status": 200,
						"status": "success",
						"results": [
							{
								"action": "accept",
								"dstaddr": [
									{
									"name": "dstaddr1",
									"q_origin_key": "dstaddr1"
									}
								],
								"dstintf": [
									{
									"name": "dstintf1",
									"q_origin_key": "dstintf1"
									}
								],
								"name": "TestPolicy",
								"policyid": 1,
								"service": [
									{
									"name": "service1",
									"q_origin_key": "service1"
									}
								],
								"srcaddr": [
									{
									"name": "srcaddr1",
									"q_origin_key": "srcaddr1"
									}
								],
								"srcintf": [
									{
									"name": "srcIntf1",
									"q_origin_key": "srcIntf1"
									}
								],
								"uuid": "uuid"
							}
						],
						"vdom": "default"
					}`),
				)
			case "/api/v2/cmdb/firewall/address/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"http_method": "GET",
						"http_status": 200,
						"status": "success",
						"results": [{"end-ip": "192.168.1.10","name": "TestName","start-ip": "192.168.1.1","subnet": "192.168.1.0/24","type": "ipmask","uuid": "uuid"}],
						"vdom": "default"
					}`),
				)
			case "/api/v2/cmdb/firewall/addrgrp/", "/api/v2/cmdb/firewall.service/group/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"http_method": "GET",
						"http_status": 200,
						"status": "success",
						"results": [
							{
								"name": "TestGroup",
								"uuid": "uuid",
								"member": [
									{
										"name": "member1"
									},
									{
										"name": "member2"
									},
									{
										"name": "member3"
									}
								]
							}
						],
						"vdom": "default"
					}`),
				)
			case "/api/v2/cmdb/firewall.service/custom/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"http_method": "GET",
						"http_status": 200,
						"status": "success",
						"results": [
							{
								"name": "TestService",
								"tcp-portrange": "80-100",
								"udp-portrange": "101-200",
								"sctp-portrange": "201-300",
								"protocol": "IP",
								"iprange": "192.168.1.0/24"
							}
						],
						"vdom": "default"
					}`),
				)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case "POST":
			switch r.URL.Path {
			case "/api/v2/cmdb/firewall/address/", "/api/v2/cmdb/firewall/addrgrp/", "/api/v2/cmdb/firewall.service/custom/", "/api/v2/cmdb/firewall.service/group/", "/api/v2/cmdb/firewall/policy/":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					{
						"http_method": "GET",
						"http_status": 200,
						"status": "success",
						"vdom": "default"
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
		w.WriteHeader(http.StatusForbidden)
		if r.Method == "GET" {
			w.Write([]byte(`{"http_method":"GET","status":"error","http_status":403,"vdom":"vdom1"}"}`))
		} else if r.Method == "POST" {
			w.Write([]byte(`{"http_method":"POST","status":"error","http_status":403,"vdom":"vdom1"}"}`))
		}
	}))
	return mockServer
}

func TestGetPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	policyResult, err := firewall.getPolicy()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, "TestPolicy", policyResult.Results[0].Name)
}

func TestGetPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	policyResult, err := firewall.getPolicy()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(policyResult.Results))
}

func TestGetPolicyNameAddressMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyNameAddressMap()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, []string{"192.168.1.0/24"}, resMap["TestName"])
}

func TestGetPolicyNameAddressMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyNameAddressMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(resMap))
}

func TestGetPolicyGroupNameAddressMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyGroupNameAddressMap()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, []string{"member1", "member2", "member3"}, resMap["TestGroup"])
}

func TestGetPolicyGroupNameAddressMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyGroupNameAddressMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(resMap))
}

func TestGetPolicyNameServicePortMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyNameServicePortMap()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, []string{"ip:192.168.1.0/24"}, resMap["TestService"])
}

func TestGetPolicyNameServicePortMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyNameServicePortMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(resMap))
}

func TestGetPolicyGroupNameServiceMap(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyGroupNameServiceMap()

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, []string{"member1", "member2", "member3"}, resMap["TestGroup"])
}

func TestGetPolicyGroupNameServiceMapWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	// Call the method to test
	resMap, err := firewall.getPolicyGroupNameServiceMap()

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(resMap))
}

func TestCreatePolicyAddress(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	addr := "192.168.1.0/24"
	addr2 := "192.168.1.1-192.168.1.10"

	// Call the method to test
	err := firewall.createPolicyAddress(addr)
	err2 := firewall.createPolicyAddress(addr2)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, err2)
}

func TestCreatePolicyAddressWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	addr := "192.168.1.0/24"

	// Call the method to test
	err := firewall.createPolicyAddress(addr)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyAddrToGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	addrList := []string{"192.168.1.1", "192.168.1.2"}

	// Call the method to test
	err := firewall.addPolicyAddrToGroup(groupName, addrList)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyAddrToGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	addrList := []string{"192.168.1.1", "192.168.1.2"}

	// Call the method to test
	err := firewall.addPolicyAddrToGroup(groupName, addrList)

	// Assert expectations
	assert.Error(t, err)
}

func TestCreatePolicyService(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	serviceTcp := "tcp:TestService"
	serviceUdp := "udp:TestService"
	serviceIp := "TestService"

	// Call the method to test
	errTcp := firewall.createPolicyService(serviceTcp)
	errUdp := firewall.createPolicyService(serviceUdp)
	errIp := firewall.createPolicyService(serviceIp)

	// Assert expectations
	assert.NoError(t, errTcp)
	assert.NoError(t, errUdp)
	assert.NoError(t, errIp)
}

func TestCreatePolicyServiceWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	service := "TestGroup"

	// Call the method to test
	err := firewall.createPolicyService(service)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicyServiceToGroup(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	serviceList := []string{"tcp:TestService", "udp:TestService", "TestService"}

	// Call the method to test
	err := firewall.addPolicyServiceToGroup(groupName, serviceList)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyServiceToGroupWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	groupName := "TestGroup"
	serviceList := []string{"tcp:TestService", "udp:TestService", "TestService"}

	// Call the method to test
	err := firewall.addPolicyServiceToGroup(groupName, serviceList)

	// Assert expectations
	assert.Error(t, err)
}

func TestAddPolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	policyName := "TestPolicy"
	action := "accept"
	srcZone := "srcZone"
	destZone := "destZone"
	srcAddrGroupName := "srcAddrGroupName"
	destAddrGroupName := "destAddrGroupName"
	serviceGroupName := "serviceGroupName"

	// Call the method to test
	err := firewall.addPolicy(policyName, action, srcZone, destZone, srcAddrGroupName, destAddrGroupName, serviceGroupName)

	// Assert expectations
	assert.NoError(t, err)
}

func TestAddPolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	policyName := "TestPolicy"
	action := "accept"
	srcZone := "srcZone"
	destZone := "destZone"
	srcAddrGroupName := "srcAddrGroupName"
	destAddrGroupName := "destAddrGroupName"
	serviceGroupName := "serviceGroupName"

	// Call the method to test
	err := firewall.addPolicy(policyName, action, srcZone, destZone, srcAddrGroupName, destAddrGroupName, serviceGroupName)

	// Assert expectations
	assert.Error(t, err)
}
