// +build json

package outbound_test

import (
	"encoding/json"
	"testing"

	//"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/vmess/outbound"
	"v2ray.com/core/testing/assert"
)

func TestConfigTargetParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "vnext": [{
      "address": "127.0.0.1",
      "port": 80,
      "users": [
        {
          "id": "e641f5ad-9397-41e3-bf1a-e8740dfed019",
          "email": "love@v2ray.com",
          "level": 255
        }
      ]
    }]
  }`

	config := new(Config)
	err := json.Unmarshal([]byte(rawJson), &config)
	assert.Error(err).IsNil()
	assert.Destination(config.Receivers[0].Destination()).EqualsString("tcp:127.0.0.1:80")
}
