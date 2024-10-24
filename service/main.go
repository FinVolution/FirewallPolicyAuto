//go:build go1.22

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FinVolution/FirewallPolicyAuto/service/config"
	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/logger"
	"github.com/FinVolution/FirewallPolicyAuto/service/router"
	"github.com/FinVolution/FirewallPolicyAuto/service/version"

	"github.com/kataras/iris/v12"
)

var (
	cmdargs      config.CmdArgs
	debugMode    bool
	printVersion bool
)

func init() {
	flag.StringVar(&cmdargs.ConfigFile, "c", "", "config file")
	flag.BoolVar(&debugMode, "d", false, "enable debug mode")
	flag.BoolVar(&printVersion, "v", false, "show version and exit")
	flag.Parse()
}

func main() {
	if printVersion {
		fmt.Println(version.VERSION)
	}

	// 初始化配置
	config.Init(cmdargs)

	// 初始化日志配置
	logger.InitLogger(
		config.Config().LogConfig.Path,
		config.Config().LogConfig.Level,
		config.Config().LogConfig.MaxSize,
		config.Config().LogConfig.Backups,
		config.Config().LogConfig.MaxAge,
	)

	// 初始化路由
	app := router.InitRouter(debugMode)

	if err := app.Run(iris.Addr(config.Config().ListenAddr), iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		logger.Errorf("server encounter an error: %s", err.Error())
		os.Exit(1)
	}
}
