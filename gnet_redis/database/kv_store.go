package database

import (
	"errors"
	"gnet_redis/model"
	"strconv"
	"strings"
)

const emptyStr = ""

type KVStore struct {
	data map[string]string
}

func NewKVStore() DataBase {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (k *KVStore) get(key string) (string, error) {
	str, ok := k.data[key]
	if !ok {
		return emptyStr, errors.New("not found")
	}
	return str, nil
}

func (k *KVStore) put(key, value string, insertStrategy bool) int64 {
	if _, ok := k.data[key]; ok && insertStrategy {
		return 0
	}
	k.data[key] = value
	return 1
}

func (k *KVStore) Get(cmd *Command) model.Reply {
	// args: get key value
	args := cmd.Args()
	// todo no memory copy
	key := string(args[0])
	str, err := k.get(key)
	if err != nil {
		return model.NewErrReply(err.Error())
	}
	return model.NewBulkReply([]byte(str))
}

func (k *KVStore) Set(cmd *Command) model.Reply {
	args := cmd.Args()
	key := string(args[0])
	value := string(args[1])

	// 支持 NX EX
	var (
		insertStrategy bool
		ttlStrategy    bool
		// ttlSeconds     int64
		ttlIndex = -1
	)

	for i := 2; i < len(args); i++ {
		flag := strings.ToLower(string(args[i]))
		switch flag {
		case "nx":
			insertStrategy = true
		case "ex":
			// 重复的 ex 指令
			if ttlStrategy {
				return model.NewSyntaxErrReply()
			}
			if i == len(args)-1 {
				return model.NewSyntaxErrReply()
			}
			ttl, err := strconv.ParseInt(string(args[i+1]), 10, 64)
			if err != nil {
				return model.NewSyntaxErrReply()
			}
			if ttl <= 0 {
				return model.NewErrReply("ERR invalid expire time")
			}

			ttlStrategy = true
			// ttlSeconds = ttl
			ttlIndex = i
			i++
		default:
			return model.NewSyntaxErrReply()
		}
	}

	// 将 args 剔除 ex 部分，进行持久化
	if ttlIndex != -1 {
		args = append(args[:ttlIndex], args[ttlIndex+2:]...)
	}

	// 设置
	affected := k.put(key, value, insertStrategy)
	if affected > 0 && ttlStrategy {
		// todo 过期时间处理
	}

	// 持久化
	if affected > 0 {
		// todo
	}

	return model.NewNilReply()
}

var _ DataBase = (*KVStore)(nil)
