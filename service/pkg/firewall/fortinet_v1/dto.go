package fortinet_v1

type fortinetPolicyItem struct {
	Name       string `json:"name"`
	QOriginKey string `json:"q_origin_key"`
}

type fortinetPolicy struct {
	ID      int64                `json:"policyid"`
	Name    string               `json:"name"`
	Action  string               `json:"action"`
	Uuid    string               `json:"uuid"`
	Status  string               `json:"status"`
	SrcIntf []fortinetPolicyItem `json:"srcintf"`
	DstIntf []fortinetPolicyItem `json:"dstintf"`
	SrcAddr []fortinetPolicyItem `json:"srcaddr"`
	DstAddr []fortinetPolicyItem `json:"dstaddr"`
	Service []fortinetPolicyItem `json:"service"`
}

type policyResp struct {
	Results []fortinetPolicy `json:"results"`
}

type nameIPv4 struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Type    string `json:"type"`
	Subnet  string `json:"subnet"`
	StartIp string `json:"start-ip"`
	EndIp   string `json:"end-ip"`
}

type vdomResp struct {
	Results []nameIPv4 `json:"results"`
	Vdom    string     `json:"vdom"`
}

type groupNameMember struct {
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Member []struct {
		Name string `json:"name"`
	} `json:"member"`
}

type groupResp struct {
	Results []groupNameMember `json:"results"`
	Vdom    string            `json:"vdom"`
}

type nameServicePort struct {
	Name          string `json:"name"`
	TcpPortRange  string `json:"tcp-portrange"`
	UdpPortRange  string `json:"udp-portrange"`
	SctpPortRange string `json:"sctp-portrange"`
	Protocol      string `json:"protocol"`
	IpRange       string `json:"iprange"`
}

type serviceResp struct {
	Results []nameServicePort `json:"results"`
	Vdom    string            `json:"vdom"`
}
