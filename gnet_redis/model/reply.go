package model

import (
	"strconv"
	"strings"
)

// CRLF 是 redis 统一的行分隔符协议
const CRLF = "\r\n"
const OK = "OK"

var (
	nilBulkBytes = []byte("$-1" + CRLF)
	pongBytes    = []byte("+PONG" + CRLF)
)

type Reply interface {
	ToBytes() []byte
}

type MultiReply interface {
	Reply
	Args() [][]byte
}

type Droplet struct {
	Reply Reply
	Err   error
}

// SimpleStringReply 简单字符串类型. 协议为 【+】【string】【CRLF】
type SimpleStringReply struct {
	Str string
}

func NewSimpleStringReply(str string) *SimpleStringReply {
	return &SimpleStringReply{
		Str: str,
	}
}

func (s *SimpleStringReply) ToBytes() []byte {
	return []byte("+" + s.Str + CRLF)
}

// IntReply 简单数字类型. 协议为 【:】【int】【CRLF】
type IntReply struct {
	Code int64
}

func NewIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

// ErrReply 错误类型. 协议为 【-】【err】【CRLF】
type ErrReply struct {
	ErrStr string
}

func NewErrReply(errStr string) *ErrReply {
	return &ErrReply{
		ErrStr: errStr,
	}
}

func (e *ErrReply) ToBytes() []byte {
	return []byte("-" + e.ErrStr + CRLF)
}

// SuccessReply success Reply
type SuccessReply struct {
}

func NewSuccessReply() *SuccessReply {
	return &SuccessReply{}
}

func (e *SuccessReply) ToBytes() []byte {
	return []byte("+" + OK + CRLF)
}

// 定长字符串类型，协议固定为 【$】【length】【CRLF】【content】【CRLF】
type BulkReply struct {
	Arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

func (b *BulkReply) ToBytes() []byte {
	if b.Arg == nil {
		return nilBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

// MultiBulkReply 数组类型. 协议固定为 【*】【arr.length】【CRLF】+ arr.length * (【$】【length】【CRLF】【content】【CRLF】)
type MultiBulkReply struct {
	args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		args: args,
	}
}

func (m *MultiBulkReply) Args() [][]byte {
	return m.args
}

func (m *MultiBulkReply) ToBytes() []byte {
	var strBuf strings.Builder
	strBuf.WriteString("*" + strconv.Itoa(len(m.args)) + CRLF)
	for _, arg := range m.args {
		if arg == nil {
			strBuf.WriteString(string(nilBulkBytes))
			continue
		}
		strBuf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
	}
	return []byte(strBuf.String())
}

var emptyMultiBulkBytes = []byte("*0\r\n")

// EmptyMultiBulkReply 空数组类型. 采用单例，协议固定为【*】【0】【CRLF】
type EmptyMultiBulkReply struct{}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

// PongReply is +PONG
type PongReply struct{}

// ToBytes marshal redis.Reply
func (r *PongReply) ToBytes() []byte {
	return pongBytes
}
