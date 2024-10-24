package config

var defaultConfig = serviceConfig{
	ListenAddr: ":8080",
	LogConfig: logConfig{
		Path:    "./logs/api.log",
		Level:   "DEBUG",
		MaxSize: 100,
		Backups: 10,
		MaxAge:  30,
	},
	FirewallConfig: []FirewallConfig{
		{
			Address:  "172.xxx.xxx.xxx",
			Brand:    "fortinet",
			Name:     "飞塔防火墙",
			Protocol: "https",
			Token:    "fortinet_rest_api_token",
			Version:  "v1",
			VirtualZone: []firewallVirtualZoneConfig{
				{
					Name: "虚拟防火墙1",
					Code: "code1",
				},
				{
					Name: "虚拟防火墙2",
					Code: "code2",
				},
			},
		},
		{
			Address:     "172.xxx.xxx.xxx",
			Brand:       "h3c",
			Name:        "H3C防火墙",
			Username:    "username",
			Password:    "password",
			Protocol:    "https",
			Version:     "v1",
			VirtualZone: []firewallVirtualZoneConfig{},
		},
	},
}
