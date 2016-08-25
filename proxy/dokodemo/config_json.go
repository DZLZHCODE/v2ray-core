// +build json

package dokodemo

import (
	"encoding/json"
	"errors"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy/registry"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type DokodemoConfig struct {
		Host         *v2net.AddressJson `json:"address"`
		PortValue    v2net.Port         `json:"port"`
		NetworkList  *v2net.NetworkList `json:"network"`
		TimeoutValue int                `json:"timeout"`
		Redirect     bool               `json:"followRedirect"`
	}
	rawConfig := new(DokodemoConfig)
	if err := json.Unmarshal(data, rawConfig); err != nil {
		return errors.New("Dokodemo: Failed to parse config: " + err.Error())
	}
	if rawConfig.Host != nil {
		this.Address = rawConfig.Host.Address
	}
	this.Port = rawConfig.PortValue
	this.Network = rawConfig.NetworkList
	this.Timeout = rawConfig.TimeoutValue
	this.FollowRedirect = rawConfig.Redirect
	return nil
}

func init() {
	registry.RegisterInboundConfig("dokodemo-door", func() interface{} { return new(Config) })
}
