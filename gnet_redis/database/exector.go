package database

import (
	"context"
	"fmt"
	"log/slog"
)

type Executor struct {
	ctx         context.Context
	cmdHandlers map[CmdType]CmdHandler
	receiver    CmdReceiver
	database    DataBase
}

func NewExecutor() *Executor {
	e := Executor{
		database: NewKVStore(),
	}
	e.cmdHandlers = map[CmdType]CmdHandler{
		// string
		CmdTypeGet: e.database.Get,
		CmdTypeSet: e.database.Set,
	}
	return &e
}

func (executor *Executor) Do(command *Command) {
	slog.Info(fmt.Sprintf("%v", command.args))
	cmdHandler := executor.getCmdHandler(command)
	// do
	reply := cmdHandler(command)
	command.receiver <- reply
}

func (executor *Executor) getCmdHandler(command *Command) CmdHandler {
	handler := executor.cmdHandlers[command.cmdType.Lower()]
	return handler
}
