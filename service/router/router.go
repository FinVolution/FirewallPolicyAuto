package router

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/FinVolution/FirewallPolicyAuto/service/router/v1"
	"github.com/FinVolution/FirewallPolicyAuto/service/version"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func InitRouter(enableLog bool) *iris.Application {
	app := iris.New()
	app.Use(recover.New(), cors)

	if enableLog {
		app.Logger().SetLevel("debug")
	}

	// Smoothly shut down the service
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
	})

	// health check api
	app.Get("/hs", func(ctx iris.Context) {
		ctx.WriteString("OK")
	})

	// version api
	app.Get("/version", func(ctx iris.Context) {
		t := time.Now().Format("2006-01-02 15:04:05")
		ctx.WriteString(fmt.Sprintf("(%s) %s", t, version.VERSION))
	})

	policyAPI := app.Party("/api/v1/policy")
	{
		policyAPI.Get("", v1.ListPolicy)    // 查询策略
		policyAPI.Post("", v1.CreatePolicy) // 创建策略
	}
	firewallAPI := app.Party("/api/v1/firewall")
	{
		firewallAPI.Get("", v1.ListFirewalls) // 查询所有防火墙
	}
	return app
}
