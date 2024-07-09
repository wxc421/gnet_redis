package parser

import (
	"fmt"
	"gnet_redis/model"
	"gnet_redis/utils"
	"log/slog"

	"testing"
)

func TestParseOne(t *testing.T) {
	parser := NewParser()
	replies := []model.Reply{
		model.NewIntReply(1),
		model.NewSimpleStringReply("OK"),
	}
	for _, reply := range replies {
		b := reply.ToBytes()
		pReply, err := parser.ParseOne(b)
		if err != nil {
			t.Fatal(err)
		}
		pb := pReply.ToBytes()
		slog.Info("ParseOne",
			slog.Any("origin", string(b)),
			slog.Any("now", string(pb)),
		)
		if !utils.BytesEquals(b, pb) {
			t.Error("parse failed: " + string(pb))
		}
	}
}

func TestParse_Get(t *testing.T) {
	parser := NewParser()
	data := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	pReply, err := parser.ParseOne(data)
	slog.Info("ParseOne",
		slog.Any("pReply", pReply),
		slog.Any("err", err),
	)

	reply := model.NewMultiBulkReply(utils.ToCmdLine("get", "key", "value"))
	stringCmd := string(reply.ToBytes())
	stringCmd = fmt.Sprintf("%#v", stringCmd)
	// *3\r\n$3\r\nget\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	fmt.Println(stringCmd)

}

func TestParse(t *testing.T) {

	cmds := []*model.MultiBulkReply{
		model.NewMultiBulkReply(utils.ToCmdLine("get", "key")),
		model.NewMultiBulkReply(utils.ToCmdLine("set", "key", "value")),
	}
	for _, cmd := range cmds {
		stringCmd := string(cmd.ToBytes())
		stringCmd = fmt.Sprintf("%#v", stringCmd)
		fmt.Println(stringCmd)
	}
}
