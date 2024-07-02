package codec

import (
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"io"
)

var ErrIncompletePacket = errors.New("incomplete packet")

type headerLenDecoder struct {
	headerLen int // TCP包的头部长度，用来描述这个包的字节长度
}

type ProtocolData struct {
	Version    uint16 // 协议版本标识
	ActionType uint16 // 行为定义
	DataLength uint32
	Data       []byte

	// headDecode bool
	// Lock       sync.Mutex
}

// NewHeaderLenDecoder 创建基于头部长度的解码器
// headerLen TCP包的头部内容，用来描述这个包的字节长度
// readMaxLen 所读取的客户端包的最大长度，客户端发送的包不能超过这个长度
func NewHeaderLenDecoder(headerLen int) StreamFrameCodec {
	if headerLen <= 0 {
		panic("headerLen or readMaxLen must must greater than 0")
	}

	return &headerLenDecoder{
		headerLen: headerLen,
	}
}

// Decode 解码 固定headerLen字节的协议
func (d *headerLenDecoder) Decode(conn gnet.Conn) ([]byte, error) {
	if conn.InboundBuffered() < d.headerLen {
		return nil, ErrIncompletePacket
	}
	headLen, _ := conn.Next(d.headerLen)
	bodySize := binary.BigEndian.Uint32(headLen)
	buf := make([]byte, bodySize)

	// 接下来尝试读取body

	// 第一种方法，用io.ReadFull 但是gnet.Conn 是不阻塞的,也就是下面buf不一定填充满就会返回 这个方式行不通
	io.ReadFull(conn, buf)

	// 第二种方法,但是如果bodySize特别大，超过了缓冲区的大小,那这个frame永远也解不出来，因为这个frame的长度超过了缓冲区的大小
	if conn.InboundBuffered() < int(bodySize) {
		return nil, ErrIncompletePacket
	}
	body, _ := conn.Next(int(bodySize))

	reader.Peek(d.headerLen + bodySize)
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
