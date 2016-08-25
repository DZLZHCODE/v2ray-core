// +build linux

package udp_test

import (
	"os"
	"syscall"
	"testing"

	"v2ray.com/core/common/alloc"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet/internal"
	. "v2ray.com/core/transport/internet/udp"
)

func TestHubSocksOption(t *testing.T) {
	assert := assert.On(t)
	if os.Geteuid() != 0 {
		// This test case requires root permission.
		return
	}

	hub, err := ListenUDP(v2net.LocalHostIP, v2net.Port(0), ListenOption{
		Callback:            func(*alloc.Buffer, *proxy.SessionInfo) {},
		ReceiveOriginalDest: true,
	})
	assert.Error(err).IsNil()
	conn := hub.Connection()

	fd, err := internal.GetSysFd(conn)
	assert.Error(err).IsNil()

	v, err := syscall.GetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT)
	assert.Error(err).IsNil()
	assert.Int(v).Equals(1)

	v, err = syscall.GetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR)
	assert.Error(err).IsNil()
	assert.Int(v).Equals(1)
}
