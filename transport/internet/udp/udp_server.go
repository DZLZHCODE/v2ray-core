package udp

import (
	"sync"
	"time"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type UDPResponseCallback func(destination v2net.Destination, payload *alloc.Buffer)

type TimedInboundRay struct {
	name       string
	inboundRay ray.InboundRay
	accessed   chan bool
	server     *UDPServer
	sync.RWMutex
}

func NewTimedInboundRay(name string, inboundRay ray.InboundRay, server *UDPServer) *TimedInboundRay {
	r := &TimedInboundRay{
		name:       name,
		inboundRay: inboundRay,
		accessed:   make(chan bool, 1),
		server:     server,
	}
	go r.Monitor()
	return r
}

func (this *TimedInboundRay) Monitor() {
	for {
		time.Sleep(time.Second * 16)
		select {
		case <-this.accessed:
		default:
			// Ray not accessed for a while, assuming communication is dead.
			this.RLock()
			if this.server == nil {
				this.RUnlock()
				return
			}
			this.server.RemoveRay(this.name)
			this.RUnlock()
			this.Release()
			return
		}
	}
}

func (this *TimedInboundRay) InboundInput() ray.OutputStream {
	this.RLock()
	defer this.RUnlock()
	if this.inboundRay == nil {
		return nil
	}
	select {
	case this.accessed <- true:
	default:
	}
	return this.inboundRay.InboundInput()
}

func (this *TimedInboundRay) InboundOutput() ray.InputStream {
	this.RLock()
	defer this.RUnlock()
	if this.inboundRay == nil {
		return nil
	}
	select {
	case this.accessed <- true:
	default:
	}
	return this.inboundRay.InboundOutput()
}

func (this *TimedInboundRay) Release() {
	log.Debug("UDP Server: Releasing TimedInboundRay: ", this.name)
	this.Lock()
	defer this.Unlock()
	if this.server == nil {
		return
	}
	this.server = nil
	this.inboundRay.InboundInput().Close()
	this.inboundRay.InboundOutput().Release()
	this.inboundRay = nil
}

type UDPServer struct {
	sync.RWMutex
	conns            map[string]*TimedInboundRay
	packetDispatcher dispatcher.PacketDispatcher
	meta             *proxy.InboundHandlerMeta
}

func NewUDPServer(meta *proxy.InboundHandlerMeta, packetDispatcher dispatcher.PacketDispatcher) *UDPServer {
	return &UDPServer{
		conns:            make(map[string]*TimedInboundRay),
		packetDispatcher: packetDispatcher,
		meta:             meta,
	}
}

func (this *UDPServer) RemoveRay(name string) {
	this.Lock()
	defer this.Unlock()
	delete(this.conns, name)
}

func (this *UDPServer) locateExistingAndDispatch(name string, payload *alloc.Buffer) bool {
	log.Debug("UDP Server: Locating existing connection for ", name)
	this.RLock()
	defer this.RUnlock()
	if entry, found := this.conns[name]; found {
		outputStream := entry.InboundInput()
		if outputStream == nil {
			return false
		}
		err := outputStream.Write(payload)
		if err != nil {
			go entry.Release()
			return false
		}
		return true
	}
	return false
}

func (this *UDPServer) Dispatch(session *proxy.SessionInfo, payload *alloc.Buffer, callback UDPResponseCallback) {
	source := session.Source
	destination := session.Destination

	// TODO: Add user to destString
	destString := source.String() + "-" + destination.String()
	log.Debug("UDP Server: Dispatch request: ", destString)
	if this.locateExistingAndDispatch(destString, payload) {
		return
	}

	log.Info("UDP Server: establishing new connection for ", destString)
	inboundRay := this.packetDispatcher.DispatchToOutbound(this.meta, session)
	timedInboundRay := NewTimedInboundRay(destString, inboundRay, this)
	outputStream := timedInboundRay.InboundInput()
	if outputStream != nil {
		outputStream.Write(payload)
	}

	this.Lock()
	this.conns[destString] = timedInboundRay
	this.Unlock()
	go this.handleConnection(timedInboundRay, source, callback)
}

func (this *UDPServer) handleConnection(inboundRay *TimedInboundRay, source v2net.Destination, callback UDPResponseCallback) {
	for {
		inputStream := inboundRay.InboundOutput()
		if inputStream == nil {
			break
		}
		data, err := inboundRay.InboundOutput().Read()
		if err != nil {
			break
		}
		callback(source, data)
	}
	inboundRay.Release()
}
