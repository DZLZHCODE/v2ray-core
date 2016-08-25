package kcp_test

import (
	"crypto/rand"
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestSimpleAuthenticator(t *testing.T) {
	assert := assert.On(t)

	buffer := alloc.NewLocalBuffer(512).Clear()
	buffer.AppendBytes('a', 'b', 'c', 'd', 'e', 'f', 'g')

	auth := NewSimpleAuthenticator()
	auth.Seal(buffer)

	assert.Bool(auth.Open(buffer)).IsTrue()
	assert.Bytes(buffer.Value).Equals([]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'})
}

func TestSimpleAuthenticator2(t *testing.T) {
	assert := assert.On(t)

	buffer := alloc.NewLocalBuffer(512).Clear()
	buffer.AppendBytes('1', '2')

	auth := NewSimpleAuthenticator()
	auth.Seal(buffer)

	assert.Bool(auth.Open(buffer)).IsTrue()
	assert.Bytes(buffer.Value).Equals([]byte{'1', '2'})
}

func BenchmarkSimpleAuthenticator(b *testing.B) {
	buffer := alloc.NewLocalBuffer(2048).Clear()
	buffer.Slice(0, 1024)
	rand.Read(buffer.Value)

	auth := NewSimpleAuthenticator()
	b.SetBytes(int64(buffer.Len()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.Seal(buffer)
		auth.Open(buffer)
	}
}
