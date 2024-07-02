package codec

import (
	"github.com/panjf2000/gnet/v2"
	"io"
)

type StreamFrameCodec interface {
	Encode(io.Writer) error                // data -> frame，并写入io.Writer
	Decode(conn gnet.Conn) ([]byte, error) // 从io.Reader中提取frame payload，并返回给上层
}
