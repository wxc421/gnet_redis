package database

import (
	"gnet_redis/model"
)

type DataBase interface {
	// Get string
	Get(*Command) model.Reply
}

type dataBase struct {
}

func NewDataBase() DataBase {
	return &dataBase{}
}

func (d *dataBase) Get(command *Command) model.Reply {
	return model.NewSuccessReply()
}
