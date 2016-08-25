package blackhole

import (
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
)

type Config struct {
	Response Response
}

type Response interface {
	WriteTo(v2io.Writer)
}

type NoneResponse struct{}

func (this *NoneResponse) WriteTo(writer v2io.Writer) {}

type HTTPResponse struct {
}

const (
	http403response = `HTTP/1.1 403 Forbidden
Connection: close
Cache-Control: max-age=3600, public
Content-Length: 0


`
)

func (this *HTTPResponse) WriteTo(writer v2io.Writer) {
	writer.Write(alloc.NewLocalBuffer(512).Clear().AppendString(http403response))
}
