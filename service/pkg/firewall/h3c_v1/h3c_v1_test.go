package h3c_v1

import (
	"strings"
	"testing"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"

	"github.com/stretchr/testify/assert"
)

func TestListPolicy(t *testing.T) {
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

	filters := map[string]string{}

	// Call the method to test
	policyList, err := firewall.ListPolicy(filters)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, "TestRule", policyList[0].Name)
	assert.Equal(t, int64(1), policyList[0].Action)
	assert.Equal(t, true, policyList[0].Enable)
	assert.Equal(t, []string{"any"}, policyList[0].SrcZone)
}

func TestListPolicyWithErr(t *testing.T) {
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

	filters := map[string]string{}

	// Call the method to test
	policyList, err := firewall.ListPolicy(filters)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, []dto.Policy(nil), policyList)
}

func TestCreatePolicy(t *testing.T) {
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

	params := dto.CreatePolicyParams{
		Title:    "test",
		SrcAddr:  []string{"192.168.1.1", "192.168.1.2"},
		DestAddr: []string{"192.168.1.3", "192.168.1.4"},
		Service:  []string{"TCP:80/65535", "UDP:11024/65535"},
		Action:   1,
		SrcZone:  "SrcZone",
		DestZone: "DestZone",
	}

	// Call the method to test
	err := firewall.CreatePolicy(params)

	// Assert expectations
	assert.NoError(t, err)
}

func TestCreatePolicyWithErr(t *testing.T) {
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

	params := dto.CreatePolicyParams{
		Title:    "test",
		SrcAddr:  []string{"192.168.1.1", "192.168.1.2"},
		DestAddr: []string{"192.168.1.3", "192.168.1.4"},
		Service:  []string{"TCP:80/65535", "UDP:11024/65535"},
		Action:   1,
		SrcZone:  "SrcZone",
		DestZone: "DestZone",
	}

	// Call the method to test
	err := firewall.CreatePolicy(params)

	// Assert expectations
	assert.Error(t, err)
}
