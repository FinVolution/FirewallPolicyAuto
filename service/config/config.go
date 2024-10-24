package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg       *serviceConfig
	once      sync.Once
	cfgLocker = new(sync.RWMutex)
)

type serviceConfig struct {
	ListenAddr     string           `mapstructure:"listen_addr"` // 监听地址
	LogConfig      logConfig        `mapstructure:"log"`
	FirewallConfig []FirewallConfig `mapstructure:"firewall"`
}

type FirewallConfig struct {
	Name        string                      `mapstructure:"name"`         // 防火墙名称
	Brand       string                      `mapstructure:"brand"`        // 品牌 当前支持可选[h3c fortinet]
	Version     string                      `mapstructure:"version"`      // API版本 当前仅支持[v1]
	Address     string                      `mapstructure:"address"`      // 地址
	Protocol    string                      `mapstructure:"protocol"`     // 请求协议 根据防火墙开启服务配置[http https]
	Username    string                      `mapstructure:"username"`     // 登录用户
	Password    string                      `mapstructure:"password"`     // 登录密码
	Token       string                      `mapstructure:"token"`        // 认证token
	VirtualZone []firewallVirtualZoneConfig `mapstructure:"virtual_zone"` // 虚拟防火墙
}

type firewallVirtualZoneConfig struct {
	Name string `mapstructure:"name"` // 虚拟防火墙名称
	Code string `mapstructure:"code"` // 虚拟防火墙编号
}

type logConfig struct {
	Path    string `mapstructure:"path"`     // 日志文件路径
	Level   string `mapstructure:"level"`    // 日志级别[DEBUG(default) INFO WARN ERROR FATAL]
	MaxSize int    `mapstructure:"max_size"` // 最大日志文件大小(MB)
	Backups int    `mapstructure:"backups"`  // 最大备份数
	MaxAge  int    `mapstructure:"max_age"`  // 最大保存天数
}

// CmdArgs 命令行参数
type CmdArgs struct {
	UseDefault bool
	ConfigFile string
}

var cmdargs *CmdArgs

// Config 配置
func Config() *serviceConfig {
	once.Do(loadConfig)
	cfgLocker.RLock()
	defer cfgLocker.RUnlock()
	return cfg
}

func configName(name string) string {
	if name != "" {
		return name
	}
	return "config"
}

func setup() {
	viper.SetConfigName(configName(cmdargs.ConfigFile))
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
}

func loadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("read config failed: %s", err.Error())
		}
	}

	var tempcfg = defaultConfig
	if err := viper.Unmarshal(&tempcfg); err != nil {
		log.Fatalf("parse config failed: %s", err.Error())
	}

	cfgLocker.Lock()
	cfg = &tempcfg
	cfgLocker.Unlock()
}

// Init 初始化
func Init(args CmdArgs) {
	cmdargs = &args

	once.Do(func() {
		setup()
		loadConfig()
	})
	// log.Printf("config info: %+v", cfg)
}
