package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"gnet_redis/database"
	"gnet_redis/model"
	"io"
)

type Handler struct {
	dbTrigger *database.Trigger
}

func (handler *Handler) HandleDroplet(conn gnet.Conn, droplets []*model.Droplet) error {
	for _, droplet := range droplets {
		if droplet.Err != nil {
			// todo EOF
			if errors.Is(droplet.Err, io.EOF) {
				return nil
			}
			_, _ = conn.Write(droplet.Reply.ToBytes())
			_ = conn.Flush()
			_ = conn.Close()
			return errors.New(fmt.Sprintf("[handler]conn request, err: %s", droplet.Err.Error()))
		}
		if droplet.Reply == nil {
			return errors.New("[handler]conn empty request")
		}
		// 请求参数必须为 multiBulkReply 类型
		multiReply, ok := droplet.Reply.(model.MultiReply)
		if !ok {
			return errors.New(fmt.Sprintf("[handler]conn invalid request: %s", droplet.Reply.ToBytes()))
		}
		receiver := make(database.CmdReceiver, 1)
		command := database.NewCommand(
			database.CmdType(multiReply.Args()[0]),
			context.Background(),
			multiReply.Args(),
			receiver,
		)
		handler.dbTrigger.EnCommand(command)
		reply := <-command.Receiver()
		_, _ = conn.Write(reply.ToBytes())
	}
	return nil
}

func NewHandler() *Handler {
	trigger := database.NewTrigger(context.Background())
	go trigger.Run()
	return &Handler{dbTrigger: trigger}
}

func (handler *Handler) Run() {
	handler.dbTrigger.Run()
}
