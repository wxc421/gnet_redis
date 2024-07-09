package database

import (
	"gnet_redis/model"
)

type DataBase interface {
	// Get string
	Get(*Command) model.Reply
	Set(*Command) model.Reply
}
