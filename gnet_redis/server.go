//go:build linux

package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gnet_redis/handler"
	"gnet_redis/parser"
	"gnet_redis/utils/iox"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type RedisServer struct {
	gnet.BuiltinEventEngine
	pool    *ants.Pool
	handler *handler.Handler
	parser  *parser.Parser
	eng     *gnet.Engine
}

func (server *RedisServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	server.eng = &eng
	slog.Info("OnBoot...")
	return gnet.None
}

func (server *RedisServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	// buf := make([]byte, 512)
	// n, err := io.ReadFull(c, buf)
	// slog.Info("OnTraffic io.ReadFull(c, buf)...", slog.Any("n", n), slog.Any("err", err))
	// n, err = c.Read(buf)
	// slog.Info("OnTraffic c.Read(buf)...", slog.Any("n", n), slog.Any("err", err))
	// bytes, err := c.Next(-1)
	// slog.Info("OnTraffic c.Next(-1)...", slog.Any("bytes", bytes), slog.Any("err", err))
	writer := iox.CountWriter{}
	buf, err := c.Peek(-1)
	io.TeeReader(bytes.NewReader(buf), &writer)
	slog.Info("c.Peek(-1)...", slog.Any("buf", buf), slog.Any("err", err))
	if errors.Is(err, io.ErrShortBuffer) {
		return gnet.None
	}
	droplets := server.parser.Parse(io.TeeReader(bytes.NewReader(buf), &writer))
	slog.Info("parser.Parse(buf)",
		slog.Any("droplets", droplets),
		slog.Any("err", err),
	)
	_, _ = c.Discard(writer.Count())
	// 具体业务在 worker pool中处理
	server.pool.Submit(func() {
		server.handler.HandleDroplet(c, droplets)
	})
	return gnet.None
}

func (server *RedisServer) Run() error {

	defer func() {
		if server.eng != nil {
			slog.Info("start close eng")
			_ = server.eng.Stop(context.Background())
			server.pool.Release()
		}
	}()

	_ = server.pool.Submit(func() {
		server.handler.Run()
	})
	_ = server.pool.Submit(func() {
		// 监听系统信号量
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case s := <-osSignal:
				slog.Info("received signal", slog.Any("signal", s))
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT,
					syscall.SIGHUP:
					if server.eng != nil {
						err := server.eng.Stop(context.Background())
						if err != nil {
							slog.Warn("stop server err", slog.Any("err", err))
						}
					}
					return
				default:
				}
			}
		}
	})
	return gnet.Run(
		server, "tcp://:8081",
		gnet.WithLogLevel(logging.DebugLevel),
		gnet.WithMulticore(true),
	)
}

func main() {
	pool, err := ants.NewPool(ants.DefaultAntsPoolSize)
	if err != nil {
		log.Fatal(err)
	}
	server := &RedisServer{
		parser:  parser.NewParser(),
		handler: handler.NewHandler(),
		pool:    pool,
	}
	log.Fatal(server.Run())
}
