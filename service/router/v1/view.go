package v1

import (
	"fmt"

	"github.com/FinVolution/FirewallPolicyAuto/service/config"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/firewall/dto"
	"github.com/FinVolution/FirewallPolicyAuto/service/utils"

	"github.com/kataras/iris/v12"
)

// ListPolicy
func ListPolicy(ctx iris.Context) {
	filterParams := ctx.URLParams()
	firewallAddress := filterParams["address"]
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pageSize", 10)
	if firewallAddress == "" {
		ctx.JSON(utils.ResponseSuccess(map[string]interface{}{"total": 0, "policyList": []dto.Policy{}}))
		return
	}
	virtualZone := filterParams["virtualZone"]
	firewallInfo := filterFirewallByAddress(firewallAddress, virtualZone)
	if firewallInfo == nil {
		ctx.JSON(utils.ResponseError(1004, fmt.Sprintf("Firewall `address=%s` and `virtualZone=%s` not found", firewallAddress, virtualZone)))
		return
	}
	// 限制分页数量
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	firewallClient, err := firewall.NewFirewallClient(firewallInfo.Brand, firewallInfo.Version, firewallInfo.Name, firewallInfo.Address,
		firewallInfo.Protocol, firewallInfo.Username, firewallInfo.Password, firewallInfo.Token, virtualZone)
	if err != nil {
		ctx.JSON(utils.ResponseError(1005, fmt.Sprintf("Get firewall client error:%s", err)))
		return
	}
	// 查询策略
	policyList, err := firewallClient.ListPolicy(filterParams)
	if err != nil {
		ctx.JSON(utils.ResponseError(1005, fmt.Sprintf("Get policy list error:%s", err)))
		return
	}

	// 数据分页
	start := pageSize * (page - 1)
	end := start + pageSize
	total := int64(len(policyList))
	data := map[string]interface{}{"total": total}
	if start >= total {
		data["policyList"] = []dto.Policy{}
	} else if end >= total && total > start {
		data["policyList"] = policyList[start:]
	} else {
		data["policyList"] = policyList[start:end]
	}

	ctx.JSON(utils.ResponseSuccess(data))
	return
}

func CreatePolicy(ctx iris.Context) {
	// 开通的逻辑是 创建跟name相关的源地址组、目的地址组和服务组，然后添加相关的策略
	// ip: 10.2.3.2, 网段10.2.3.0/24, 范围：10.2.3.6-10.3.6.9
	var params dto.CreatePolicyParams
	if err := ctx.ReadJSON(&params); err != nil {
		ctx.JSON(utils.ResponseError(1004, fmt.Sprintf("Params format incorrect: %s", err.Error())))
		return
	}
	// 参数校验
	if err := utils.IniValidator().Check(params); err != nil {
		ctx.JSON(utils.ResponseError(1004, fmt.Sprintf("Params validated failed: %s", err.Error())))
		return
	}

	// 查询负责上网权限的防火墙  可能有多个
	firewallInfo := filterFirewallByAddress(params.FirewallAddress, params.VirtualZone)
	if firewallInfo == nil {
		ctx.JSON(utils.ResponseError(1004, fmt.Sprintf("Firewall `address=%s` and `virtualZone=%s` not found", params.FirewallAddress, params.VirtualZone)))
		return
	}
	firewallClient, err := firewall.NewFirewallClient(firewallInfo.Brand, firewallInfo.Version, firewallInfo.Name, firewallInfo.Address,
		firewallInfo.Protocol, firewallInfo.Username, firewallInfo.Password, firewallInfo.Token, params.VirtualZone)
	if err != nil {
		ctx.JSON(utils.ResponseError(1004, fmt.Sprintf("Get firewall client error:%s", err)))
		return
	}
	if err := firewallClient.CreatePolicy(params); err != nil {
		ctx.JSON(utils.ResponseError(1005, fmt.Sprintf("Create policy error:%s", err)))
		return
	}
	ctx.JSON(utils.ResponseSuccess(nil))
	return
}

// ListFirewalls 查询所有防火墙
func ListFirewalls(ctx iris.Context) {
	data := []map[string]interface{}{}
	for _, v := range config.Config().FirewallConfig {
		vz := []map[string]string{}
		for _, vv := range v.VirtualZone {
			vz = append(vz, map[string]string{"name": vv.Name, "code": vv.Code})
		}
		data = append(data, map[string]interface{}{
			"address":     v.Address,
			"brand":       v.Brand,
			"name":        v.Name,
			"virtualZone": vz,
		})
	}
	ctx.JSON(utils.ResponseSuccess(data))
	return
}

// filterFirewallByAddress 根据防火墙地址和虚拟zone查询防火墙
func filterFirewallByAddress(address, virtualZone string) *config.FirewallConfig {
	for _, v := range config.Config().FirewallConfig {
		if v.Address == address {
			if virtualZone == "" {
				return &v
			}
			for _, vv := range v.VirtualZone {
				if vv.Code == virtualZone {
					return &v
				}
			}
		}
	}
	return nil
}
