package dto

type Policy struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Action       int64    `json:"action"`
	Enable       bool     `json:"enable"`
	SrcZone      []string `json:"srcZone"`
	DestZone     []string `json:"destZone"`
	SrcAddress   []string `json:"srcAddress"`
	DestAddress  []string `json:"destAddress"`
	ServicePort  []string `json:"servicePort"`
	FirewallName string   `json:"firewallName"`
}

type CreatePolicyParams struct {
	Title           string   `json:"title" validate:"required"`
	Action          int64    `json:"action" validate:"required,oneof=1 2"`
	SrcZone         string   `json:"srcZone" validate:"omitempty"`
	DestZone        string   `json:"destZone" validate:"omitempty"`
	SrcAddr         []string `json:"srcAddr" validate:"required"`
	DestAddr        []string `json:"destAddr" validate:"required"`
	Service         []string `json:"service" validate:"required"`
	FirewallAddress string   `json:"firewallAddress" validate:"required"`
	VirtualZone     string   `json:"virtualZone" validate:"omitempty"`
}
