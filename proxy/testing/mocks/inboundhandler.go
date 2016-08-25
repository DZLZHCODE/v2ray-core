package mocks

import (
	"io"
	"sync"

	"v2ray.com/core/app/dispatcher"
	v2io "v2ray.com/core/common/io"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type InboundConnectionHandler struct {
	ListeningPort    v2net.Port
	ListeningAddress v2net.Address
	PacketDispatcher dispatcher.PacketDispatcher
	ConnInput        io.Reader
	ConnOutput       io.Writer
}

func (this *InboundConnectionHandler) Start() error {
	return nil
}

func (this *InboundConnectionHandler) Port() v2net.Port {
	return this.ListeningPort
}

func (this *InboundConnectionHandler) Close() {

}

func (this *InboundConnectionHandler) Communicate(destination v2net.Destination) error {
	ray := this.PacketDispatcher.DispatchToOutbound(&proxy.InboundHandlerMeta{
		AllowPassiveConnection: false,
	}, &proxy.SessionInfo{
		Source:      v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(0)),
		Destination: destination,
	})

	input := ray.InboundInput()
	output := ray.InboundOutput()

	readFinish := &sync.Mutex{}
	writeFinish := &sync.Mutex{}

	readFinish.Lock()
	writeFinish.Lock()

	go func() {
		v2reader := v2io.NewAdaptiveReader(this.ConnInput)
		defer v2reader.Release()

		v2io.Pipe(v2reader, input)
		input.Close()
		readFinish.Unlock()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(this.ConnOutput)
		defer v2writer.Release()

		v2io.Pipe(output, v2writer)
		output.Release()
		writeFinish.Unlock()
	}()

	readFinish.Lock()
	writeFinish.Lock()
	return nil
}
