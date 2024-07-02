package database

import (
	"context"
	tb "gnet_redis/internal/queue"
	"gnet_redis/internal/queue/lock_free_queue"
	"gnet_redis/model"
	"runtime"
)

type CmdHandler func(*Command) model.Reply

type Trigger struct {
	ctx      context.Context
	queue    tb.Queue[*Command]
	executor *Executor
}

func NewTrigger(ctx context.Context) *Trigger {
	return &Trigger{
		ctx:      ctx,
		queue:    lock_free_queue.NewLKQueue[*Command](),
		executor: NewExecutor(),
	}
}

func (trigger *Trigger) EnCommand(command *Command) {
	trigger.queue.Enqueue(command)
}

func (trigger *Trigger) Run() {
	for {
		select {
		case <-trigger.ctx.Done():
			return
		default:
			command := trigger.queue.Dequeue()
			if command == nil {
				runtime.Gosched()
				continue
			}
			trigger.executor.Do(command)
		}
	}
}
