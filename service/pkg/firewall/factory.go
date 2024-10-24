package firewall

import (
	"errors"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/fortinet_v1"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/h3c_v1"
)

type FirewallClient interface {
	ListPolicy(filters map[string]string) ([]dto.Policy, error)
	CreatePolicy(params dto.CreatePolicyParams) error
}

// NewFirewallClient 工厂函数  返回对应品牌的对象
func NewFirewallClient(brand, version, name, address, protocol, username, password, tokenID, virtualZone string) (FirewallClient, error) {
	switch brand {
	case "h3c":
		switch version {
		case "v1":
			return &h3c_v1.FirewallH3CV1{
				Username: username,
				Password: password,
				Address:  address,
				Protocol: protocol,
				Name:     name,
			}, nil
		default:
			return nil, errors.New("unsupported H3C version")
		}
	case "fortinet":
		switch version {
		case "v1":
			return &fortinet_v1.FirewallFortinetV1{
				Name:        name,
				Address:     address,
				Protocol:    protocol,
				TokenID:     tokenID,
				VirtualZone: virtualZone,
			}, nil
		default:
			return nil, errors.New("unsupported Fortinet version")
		}
	default:
		return nil, errors.New("unsupported firewall brand")
	}
}
