package persistence

import (
	"context"
	tb "gnet_redis/internal/queue"
	"gnet_redis/model"
	"os"
	"time"
)

// aof 持久化等级 always | everysec | no
type appendSyncStrategy string

const (
	alwaysAppendSyncStrategy   appendSyncStrategy = "always"   // 每条指令都进行持久化落盘
	everysecAppendSyncStrategy appendSyncStrategy = "everysec" // 每秒批量执行一次持久化落盘
	noAppendSyncStrategy       appendSyncStrategy = "no"       // 不主动进行指令的持久化落盘，由设备自行决定落盘节奏
)

// Persistence 持久化相关接口
type Persistence interface {
	PersistCmd(ctx context.Context, reply *model.MultiBulkReply) error
}

type AofPersistence struct {
	aofFile     *os.File
	aofFileName string
	ctx         context.Context
}

func (a *AofPersistence) write(reply *model.MultiBulkReply) error {
	bytes := reply.ToBytes()
	size := len(bytes)
	var nn = 0

	for {
		n, err := a.aofFile.Write(bytes)
		if err != nil {
			return err
		}
		nn += n
		if size <= nn {
			break
		}
	}
	return nil
}

func (a *AofPersistence) flush() error {
	return a.aofFile.Sync()
}

// AofPersistenceAlwaysAppend 每条指令都进行持久化落盘
type AofPersistenceAlwaysAppend struct {
	*AofPersistence
}

// AofPersistenceEverySecAppend 每秒批量执行一次持久化落盘
type AofPersistenceEverySecAppend struct {
	*AofPersistence
	queue tb.Queue[*model.MultiBulkReply]
}

func (a *AofPersistenceAlwaysAppend) PersistCmd(ctx context.Context, reply *model.MultiBulkReply) error {
	if err := a.write(reply); err != nil {
		return err
	}
	return a.aofFile.Sync()
}

func (a *AofPersistenceEverySecAppend) PersistCmd(ctx context.Context, reply *model.MultiBulkReply) error {
	a.queue.Enqueue(reply)
	return nil
}

func (a *AofPersistenceEverySecAppend) Run() {
	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			for {
				reply := a.queue.Dequeue()
				if reply == nil {
					break
				}
				if err := a.write(reply); err != nil {
					// log
					continue
				}
			}
			if err := a.flush(); err != nil {
				// log
			}
		case <-a.ctx.Done():
			return

		}
	}
}

func (a *AofPersistence) Close() {
	// TODO implement me
	panic("implement me")
}
