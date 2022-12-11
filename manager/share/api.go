package share

import "github.com/rayguo17/go-socks/util"

type LightConfig struct {
	RemoteAddr util.Address
	PublicKey  string
	NodeId     string
}

func NewLightConfig(addr util.Address, pk string, nd string) *LightConfig {
	return &LightConfig{
		RemoteAddr: addr,
		PublicKey:  pk,
		NodeId:     nd,
	}
}
