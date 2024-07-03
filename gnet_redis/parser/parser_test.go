package parser

import (
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
