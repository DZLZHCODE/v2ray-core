package ws_test

import (
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/ws"
)

func Test_Connect_ws(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "ws", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "wss", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_1_nil(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "wss", Path: ""}).Apply()
	conn, err := Dial(nil, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_ws_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess_fail(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	_, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("static.kkdev.org"), 443))
	assert.Error(err).IsNotNil()
}

func Test_Connect_wss_guess_fail_port(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	_, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("static.kkdev.org"), 179))
	assert.Error(err).IsNotNil()
}

func Test_Connect_wss_guess_reuse(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: "", ConnectionReuse: true}).Apply()
	i := 3
	for i != 0 {
		conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
		assert.Error(err).IsNil()
		conn.Write([]byte("echo"))
		s := make(chan int)
		go func() {
			buf := make([]byte, 4)
			conn.Read(buf)
			str := string(buf)
			if str != "echo" {
				assert.Fail("Data mismatch")
			}
			s <- 0
		}()
		<-s
		if i == 0 {
			conn.SetDeadline(time.Now())
			conn.SetReadDeadline(time.Now())
			conn.SetWriteDeadline(time.Now())
			conn.SetReusable(false)
		}
		conn.Close()
		i--
	}
}

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "ws", Path: "ws"}).Apply()
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13142)
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142))
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 5)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142))
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 15)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142))
	assert.Error(err).IsNil()
	conn.Close()
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	assert := assert.On(t)
	go func() {
		<-time.After(time.Second * 5)
		assert.Fail("Too slow")
	}()
	(&Config{Pto: "wss", Path: "wss", ConnectionReuse: true, DeveloperInsecureSkipVerify: true, PrivKey: "./../../../testing/tls/key.pem", Cert: "./../../../testing/tls/cert.pem"}).Apply()
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13143)
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13143))
	assert.Error(err).IsNil()
	conn.Close()
}
