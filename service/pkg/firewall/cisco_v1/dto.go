package cisco_v1

type FirewallCiscoV1 struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	TokenID     string `json:"tokenID"`
	Name        string `json:"name"`
	VirtualZone string `json:"virtualZone"`
	Protocol    string `json:"protocol"`
}

type ItemResp struct {
	Kind     string `json:"kind"`
	Value    string `json:"value"`
	RefLink  string `json:"refLink"`
	ObjectId string `json:"objectId"`
}

type NetworkObjGroupStruct struct {
	Kind     string        `json:"kind"`
	Name     string        `json:"name"`
	Members  []interface{} `json:"members"`
	ObjectId string        `json:"objectId"`
}

type NetworkObjStruct struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	Host struct {
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"host"`
	ObjectId string `json:"objectId"`
}

type ciscoPolicy struct {
	ItemResp
	Permit             bool        `json:"permit"`             // true为permit，false为deny
	Kind               string      `json:"kind"`               // 类型
	SrcZone            []string    `json:"srcZone"`            // 源区域
	DestZone           []string    `json:"destZone"`           // 目标区域
	SourceAddress      interface{} `json:"sourceAddress"`      //源地址
	SrcAddress         []string    `json:"SrcAddress"`         // 源地址
	DestinationAddress interface{} `json:"destinationAddress"` //目的地址
	DestAddress        []string    `json:"DestAddress"`        //目的地址
	DestinationService interface{} `json:"destinationService"` //目标服务
	DestServicePort    []string    `json:"destServicePort"`    //目标服务端口
	Active             bool        `json:"active"`             // 是否启用
	Position           int         `json:"position"`           // 策略顺序
	FirewallName       string      `json:"FirewallName"`       // 防火墙名称
	ObjectId           string      `json:"objectId"`           // 策略ID

	SourceService interface{} `json:"sourceService"` // 源服务
	SelfLink      string      `json:"selfLink"`      // 策略链接
	IsAccessRule  bool        `json:"isAccessRule"`  // 是否为访问策略
	Remarks       interface{} `json:"remarks"`       // 备注
	RuleLogging   interface{} `json:"ruleLogging"`   // 是否记录日志

}
