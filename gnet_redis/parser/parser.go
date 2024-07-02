package parser

import (
	"bufio"
	"bytes"
	"errors"
	"gnet_redis/model"
	"log"

	"io"
	"strconv"
)

type lineParser func(header []byte, reader *bufio.Reader) *model.Droplet

type Parser struct {
	lineParsers map[byte]lineParser
	logger      log.Logger
}

func NewParser() *Parser {
	p := Parser{}
	p.lineParsers = map[byte]lineParser{
		'+': p.parseSimpleString,
		'-': p.parseError,
		':': p.parseInt,
		'$': p.parseBulk,
		'*': p.parseMultiBulk,
	}
	return &p
}

// 解析
func (p *Parser) parseMultiBulk(header []byte, reader *bufio.Reader) (droplet *model.Droplet) {
	var _err error
	defer func() {
		if _err != nil {
			droplet = &model.Droplet{
				Reply: model.NewErrReply(_err.Error()),
				Err:   _err,
			}
		}
	}()

	// 获取数组长度
	length, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil {
		_err = err
		return
	}

	if length <= 0 {
		return &model.Droplet{
			Reply: model.NewEmptyMultiBulkReply(),
		}
	}

	lines := make([][]byte, 0, length)
	for i := int64(0); i < length; i++ {
		// 获取每个 bulk 首行
		firstLine, err := reader.ReadBytes('\n')
		if err != nil {
			_err = err
			return
		}

		// bulk 首行格式校验
		length := len(firstLine)
		if length < 4 || firstLine[length-2] != '\r' || firstLine[length-1] != '\n' || firstLine[0] != '$' {
			continue
		}

		// bulk 解析
		bulkBody, err := p.parseBulkBody(firstLine[:length-2], reader)
		if err != nil {
			_err = err
			return
		}

		lines = append(lines, bulkBody)
	}

	return &model.Droplet{
		Reply: model.NewMultiBulkReply(lines),
	}
}

// parseSimpleString 解析简单 string 类型
func (p *Parser) parseSimpleString(header []byte, reader *bufio.Reader) *model.Droplet {
	content := header[1:]
	return &model.Droplet{
		Reply: model.NewSimpleStringReply(string(content)),
	}
}

// parseInt 解析简单 int 类型
func (p *Parser) parseInt(header []byte, reader *bufio.Reader) *model.Droplet {

	i, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil {
		return &model.Droplet{
			Err:   err,
			Reply: model.NewErrReply(err.Error()),
		}
	}

	return &model.Droplet{
		Reply: model.NewIntReply(i),
	}
}

// parseError 解析错误类型
func (p *Parser) parseError(header []byte, reader *bufio.Reader) *model.Droplet {
	return &model.Droplet{
		Reply: model.NewErrReply(string(header[1:])),
	}
}

// 解析定长 string 类型
func (p *Parser) parseBulk(header []byte, reader *bufio.Reader) *model.Droplet {
	// 解析定长 string
	body, err := p.parseBulkBody(header, reader)
	if err != nil {
		return &model.Droplet{
			Reply: model.NewErrReply(err.Error()),
			Err:   err,
		}
	}
	return &model.Droplet{
		Reply: model.NewBulkReply(body),
	}
}

// parseBulkBody 解析定长 string
func (p *Parser) parseBulkBody(header []byte, reader *bufio.Reader) ([]byte, error) {
	// 获取 string 长度
	strLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	if err != nil {
		return nil, err
	}

	// 长度 + 2，把 CRLF 也考虑在内
	body := make([]byte, strLen+2)
	// 从 reader 中读取对应长度
	if _, err = io.ReadFull(reader, body); err != nil {
		return nil, err
	}
	return body[:len(body)-2], nil
}

// ParseStream reads data from io.Reader and send payloads through channel
func (p *Parser) ParseStream(reader io.Reader) <-chan *model.Droplet {
	ch := make(chan *model.Droplet)
	go p.parse0(reader, ch)
	return ch
}

// ParseBytes reads data from []byte and return all replies
func (p *Parser) ParseBytes(data []byte) ([]model.Reply, error) {
	ch := make(chan *model.Droplet)
	reader := bytes.NewReader(data)
	go p.parse0(reader, ch)
	var results []model.Reply
	for droplet := range ch {
		if droplet == nil {
			return nil, errors.New("no protocol")
		}
		if droplet.Err != nil {
			if droplet.Err == io.EOF {
				break
			}
			return nil, droplet.Err
		}
		results = append(results, droplet.Reply)
	}
	return results, nil
}

// ParseReader reads data from reader and return all replies
func (p *Parser) ParseReader(reader io.Reader) ([]model.Reply, error) {
	ch := make(chan *model.Droplet)
	go p.parse0(reader, ch)
	var results []model.Reply
	for droplet := range ch {
		if droplet == nil {
			return nil, errors.New("no protocol")
		}
		if droplet.Err != nil {
			if droplet.Err == io.EOF {
				break
			}
			return nil, droplet.Err
		}
		results = append(results, droplet.Reply)
	}
	return results, nil
}

// ParseOne reads data from []byte and return the first payload
func (p *Parser) ParseOne(data []byte) (model.Reply, error) {
	ch := make(chan *model.Droplet)
	reader := bytes.NewReader(data)
	go p.parse0(reader, ch)
	droplet := <-ch
	if droplet == nil {
		return nil, errors.New("no protocol")
	}
	return droplet.Reply, droplet.Err
}

func (p *Parser) parse0(rawReader io.Reader, ch chan<- *model.Droplet) {
	reader := bufio.NewReader(rawReader)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			ch <- &model.Droplet{
				Reply: model.NewErrReply(err.Error()),
				Err:   err,
			}
			close(ch)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// there are some empty lines within replication traffic, ignore this error
			// protocolError(ch, "empty line")
			continue
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		lineParseFunc, ok := p.lineParsers[line[0]]
		if !ok {
			continue
		}
		ch <- lineParseFunc(line, reader)
	}
}

func (p *Parser) Parse(rawReader io.Reader) []*model.Droplet {
	reader := bufio.NewReader(rawReader)
	droplets := make([]*model.Droplet, 0, 1)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			droplets = append(droplets, &model.Droplet{
				Reply: model.NewErrReply(err.Error()),
				Err:   err,
			})
			return droplets
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// there are some empty lines within replication traffic, ignore this error
			// protocolError(ch, "empty line")
			continue
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		lineParseFunc, ok := p.lineParsers[line[0]]
		if !ok {
			continue
		}
		droplets = append(droplets, lineParseFunc(line, reader))
	}
}
