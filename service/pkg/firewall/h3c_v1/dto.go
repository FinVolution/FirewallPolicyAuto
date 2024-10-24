package h3c_v1

// 模拟的策略结构体
type h3cPolicy struct {
	ID                 int64               `json:"id"`
	Name               string              `json:"name"`
	Action             int64               `json:"action"`
	Enable             bool                `json:"enable"`
	SrcZoneList        map[string][]string `json:"srcZoneList"`
	DestZoneList       map[string][]string `json:"destZoneList"`
	ServGrpList        map[string][]string `json:"servGrpList"`
	SrcSimpleAddrList  map[string][]string `json:"srcSimpleAddrList"`
	DestSimpleAddrList map[string][]string `json:"destSimpleAddrList"`
	ServObjList        map[string][]string `json:"servObjList"`
	SrcAddrList        map[string][]string `json:"srcAddrList"`
	DestAddrList       map[string][]string `json:"destAddrList"`
}

// IPv4地址组结构体
type iPv4Objs struct {
	Group             string `json:"Group"`
	Type              int    `json:"Type"`
	StartIPv4Address  string `json:"StartIPv4Address"`
	EndIPv4Address    string `json:"EndIPv4Address"`
	HostIPv4Address   string `json:"HostIPv4Address"`
	SubnetIPv4Address string `json:"SubnetIPv4Address"`
	IPv4Mask          string `json:"IPv4Mask"`
	NestedGroup       string `json:"NestedGroup"`
}

// 服务端口组结构体
type serviceObjs struct {
	Group         string `json:"Group"`
	Type          int    `json:"Type"`
	StartSrcPort  int    `json:"StartSrcPort"`
	StartDestPort int    `json:"StartDestPort"`
	EndSrcPort    int    `json:"EndSrcPort"`
	EndDestPort   int    `json:"EndDestPort"`
	NestedGroup   string `json:"NestedGroup"`
}

type iPv4ServObj struct {
	ID          int64 `json:"ID"`
	ServObjList struct {
		ServObjItem []serviceObjsCopy `json:"ServObjItem"`
	} `json:"ServObjList"`
}

type serviceObjsCopy struct {
	Group         string `json:"Group"`
	Type          string `json:"Type"`
	StartSrcPort  string `json:"StartSrcPort"`
	StartDestPort string `json:"StartDestPort"`
	EndSrcPort    string `json:"EndSrcPort"`
	EndDestPort   string `json:"EndDestPort"`
	NestedGroup   string `json:"NestedGroup"`
}

type simpleAddr struct {
	ID             int64 `json:"ID"`
	SimpleAddrList struct {
		SimpleAddrItem []string `json:"SimpleAddrItem"`
	} `json:"SimpleAddrList"`
}

type policyGroup struct {
	ID       int64 `json:"ID"`
	NameList struct {
		NameItem []string `json:"NameItem"`
	} `json:"NameList"`
}

type tokenResp struct {
	TokenID    string `json:"token-id"`
	Link       string `json:"link"`
	ExpiryTime string `json:"expiry-time"`
}

type rulesResp struct {
	GetRules []h3cPolicy `json:"GetRules"`
}

type policyGroupAddressResp struct {
	IPv4Objs []iPv4Objs `json:"IPv4Objs"`
}

type policyGroupServicePortResp struct {
	ServObjs []serviceObjs `json:"ServObjs"`
}

type srcPolicySimpleAddrResp struct {
	IPv4SrcSimpleAddr []simpleAddr `json:"IPv4SrcSimpleAddr"`
}

type destPolicySimpleAddrResp struct {
	IPv4DestSimpleAddr []simpleAddr `json:"IPv4DestSimpleAddr"`
}

type servicePortResp struct {
	IPv4ServObj []iPv4ServObj `json:"IPv4ServObj"`
}

type srcPolicyIDGroupResp struct {
	IPv4SrcAddr []policyGroup `json:"IPv4SrcAddr"`
}

type destPolicyIDGroupResp struct {
	IPv4DestAddr []policyGroup `json:"IPv4DestAddr"`
}

type servicePolicyIDGroupResp struct {
	IPv4ServGrp []policyGroup `json:"IPv4ServGrp"`
}
