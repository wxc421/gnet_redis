package database

import (
	"context"
	"gnet_redis/model"
	"strings"
)

type CmdType string

type CmdReceiver chan model.Reply

func (c CmdType) String() string {
	return strings.ToLower(string(c))
}

func (c CmdType) Lower() CmdType {
	return CmdType(strings.ToLower(string(c)))
}

type Command struct {
	cmdType  CmdType
	ctx      context.Context
	args     [][]byte
	receiver CmdReceiver
}

func NewCommand(cmdType CmdType, ctx context.Context, args [][]byte, receiver CmdReceiver) *Command {
	return &Command{cmdType: cmdType, ctx: ctx, receiver: receiver, args: args}
}

func (c *Command) Receiver() CmdReceiver {
	return c.receiver
}

const (
	CmdTypeExpire   CmdType = "expire"
	CmdTypeExpireAt CmdType = "expireat"

	/*
		redis string type
	*/
	CmdTypeGet  CmdType = "get"
	CmdTypeSet  CmdType = "set"
	CmdTypeMGet CmdType = "mget"
	CmdTypeMSet CmdType = "mset"

	/*
		redis list type
	*/
	CmdTypeLPush  CmdType = "lpush"
	CmdTypeLPop   CmdType = "lpop"
	CmdTypeRPush  CmdType = "rpush"
	CmdTypeRPop   CmdType = "rpop"
	CmdTypeLRange CmdType = "lrange"

	/*
		redis hash type
	*/
	CmdTypeHSet CmdType = "hset"
	CmdTypeHGet CmdType = "hget"
	CmdTypeHDel CmdType = "hdel"

	/*
		redis set type
	*/
	CmdTypeSAdd      CmdType = "sadd"
	CmdTypeSIsMember CmdType = "sismember"
	CmdTypeSRem      CmdType = "srem"

	/*
		redis sorted set type
	*/
	CmdTypeZAdd          CmdType = "zadd"
	CmdTypeZRangeByScore CmdType = "zrangebyscore"
	CmdTypeZRem          CmdType = "zrem"
)
