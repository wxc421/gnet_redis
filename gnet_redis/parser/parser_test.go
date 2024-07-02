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

func TestName(t *testing.T) {
	// data := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	// data = []byte("+OK\r\n")
	// {
	// 	data = []byte("+OK\r\n")
	// 	data = []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	// 	// fmt.Println(len(data))
	// 	// reader := bufio.NewReader(bytes.NewReader(data))
	// 	// slice, err := reader.ReadSlice('\n')
	// 	// fmt.Println(slice, err)
	// 	replies, err := ParseBytes(data)
	// 	fmt.Println(err)
	// 	for _, reply := range replies {
	// 		fmt.Println(string(reply.ToBytes()))
	// 	}
	// }
	// {
	// 	data = []byte(`+OK\r\n`)
	// 	fmt.Println(len(data))
	// 	reader := bufio.NewReader(bytes.NewReader(data))
	// 	slice, err := reader.ReadSlice('\n')
	// 	fmt.Println(slice, err)
	// }
	// reader.
	// 	replies, err := ParseBytes(bytes.NewReader(data))
	// fmt.Println(err)
	// for _, reply := range replies {
	// 	fmt.Println(reply)
	// }
}
