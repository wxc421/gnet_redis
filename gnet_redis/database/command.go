package database

import (
	"context"
	"gnet_redis/model"
	"gnet_redis/utils/poolx"
	"strings"
)

type CmdType string

type CmdReceiver chan model.Reply

var commandPool = poolx.NewPoolNormal[*Command](func() *Command {
	return &Command{}
})

func GetCommand() *Command {
	command := commandPool.Get()
	return command
}

func PutCommand(command *Command) {
	command.ctx = nil
	command.receiver = nil
	command.args = nil
	commandPool.Put(command)
}

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

func (c *Command) SetCmdType(cmdType CmdType) *Command {
	c.cmdType = cmdType
	return c
}

func (c *Command) SetCtx(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}

func (c *Command) SetArgs(args [][]byte) *Command {
	c.args = args
	return c
}

func (c *Command) SetReceiver(receiver CmdReceiver) *Command {
	c.receiver = receiver
	return c
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
