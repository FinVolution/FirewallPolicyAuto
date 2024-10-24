package fortinet_v1

import (
	"strings"
	"testing"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"

	"github.com/stretchr/testify/assert"
)

func TestListPolicy(t *testing.T) {

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

	policyResults, err := firewall.ListPolicy(map[string]string{})

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, "TestPolicy", policyResults[0].Name)
}

func TestListPolicyWithErr(t *testing.T) {

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

	policyResults, err := firewall.ListPolicy(map[string]string{})

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, 0, len(policyResults))
}

func TestCreatePolicy(t *testing.T) {
	mockServer := setupSuccessMockServer()
	defer mockServer.Close()

	params := dto.CreatePolicyParams{
		Title:    "test",
		SrcAddr:  []string{"192.168.1.1", "192.168.1.2"},
		DestAddr: []string{"192.168.1.3", "192.168.1.4"},
		Service:  []string{"tcp:TestService", "udp:TestService", "TestService"},
		Action:   1,
		SrcZone:  "SrcZone",
		DestZone: "DestZone",
	}

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	err := firewall.CreatePolicy(params)

	// Assert expectations
	assert.NoError(t, err)
}

func TestCreatePolicyWithErr(t *testing.T) {
	mockServer := setupFailedMockServer()
	defer mockServer.Close()

	params := dto.CreatePolicyParams{
		Title:    "test",
		SrcAddr:  []string{"192.168.1.1", "192.168.1.2"},
		DestAddr: []string{"192.168.1.3", "192.168.1.4"},
		Service:  []string{"tcp:TestService", "udp:TestService", "TestService"},
		Action:   1,
		SrcZone:  "SrcZone",
		DestZone: "DestZone",
	}

	// Create an instance of FirewallFortinetV1
	firewall := FirewallFortinetV1{
		Name:        "testFirewall",
		Address:     strings.Split(mockServer.URL, "/")[2],
		Protocol:    strings.ReplaceAll(strings.Split(mockServer.URL, "//")[0], ":", ""),
		TokenID:     "fakeTokenID",
		VirtualZone: "",
	}

	err := firewall.CreatePolicy(params)

	// Assert expectations
	assert.Error(t, err)
}
