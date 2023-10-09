package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"redis-based-on-go/log"
	"redis-based-on-go/redis/protocol"
	"runtime/debug"
	"strconv"
	"strings"
)

/*
	客户端请求解析器
*/

// PayLoad 保存 redis 协议数据或错误
type PayLoad struct {
	Data protocol.Reply
	Err  error
}

const (
	SIMPLE_STRING = '+'
	ERROR         = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

// Parse 从 conn 中读取数据并放入 channel
func Parse(conn net.Conn) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	go parse(conn, ch)
	return ch
}

func parse(conn net.Conn, ch chan<- *PayLoad) {
	defer func() {
		if err := recover(); err != nil {
			mylog.Logger.Error(err, string(debug.Stack()))
		}
	}()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &PayLoad{Err: err}
			close(ch)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// 该行不以 \r\n 结尾 忽略该行
			continue
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		payLoad := getNextPayLoad(line, reader)
		ch <- payLoad
		if payLoad.Err != nil && !strings.Contains(payLoad.Err.Error(), "数字解析失败:") {
			close(ch)
		}
		//switch line[0] {
		//case SIMPLE_STRING:
		//	ch <- &PayLoad{Data: &protocol.SimpleString{Data: string(line[1:])}}
		//case ERROR:
		//	ch <- &PayLoad{Data: &protocol.Error{Data: string(line[1:])}}
		//case INTEGER:
		//	value, err := strconv.Atoi(string(line[1:]))
		//	if err != nil {
		//		ch <- &PayLoad{Err: errors.New("数字解析失败:" + string(line[1:]))}
		//	} else {
		//		ch <- &PayLoad{Data: &protocol.Integer{Data: value}}
		//	}
		//case BULK_STRING:
		//	err = parseBulkString(line, reader, ch)
		//	if err != nil {
		//		ch <- &PayLoad{Err: err}
		//		close(ch)
		//		return
		//	}
		//case ARRAY:
		//	err = parseArray(line, reader, ch)
		//	if err != nil {
		//		ch <- &PayLoad{Err: err}
		//		close(ch)
		//		return
		//	}
		//default:
		//	ch <- &PayLoad{Err: errors.New("未知首字节")}
		//	close(ch)
		//	return
		//}
	}
}

// getNextPayLoad 解析下一个数据(客户端命令)
func getNextPayLoad(line []byte, reader *bufio.Reader) *PayLoad {
	switch line[0] {
	case SIMPLE_STRING:
		return &PayLoad{Data: &protocol.SimpleString{Data: string(line[1:])}}
	case ERROR:
		return &PayLoad{Data: &protocol.Error{Data: string(line[1:])}}
	case INTEGER:
		value, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return &PayLoad{Err: errors.New("数字解析失败:" + string(line[1:]))}
		} else {
			return &PayLoad{Data: &protocol.Integer{Data: value}}
		}
	case BULK_STRING:
		err, payLoad := parseBulkString(line, reader)
		if err != nil {
			return &PayLoad{Err: err}
		} else {
			return &PayLoad{Data: payLoad.Data}
		}
	case ARRAY:
		err, payLoad := parseArray(line, reader)
		if err != nil {
			return &PayLoad{Err: err}
		} else {
			return &PayLoad{Data: payLoad.Data}
		}
	default:
		return &PayLoad{Err: errors.New("未知首字节")}
	}
}

// parseArray 解析数组 line[1:] 为数组元素个数
func parseArray(line []byte, reader *bufio.Reader) (error, *PayLoad) {
	elementCount, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return err, nil
	}
	array := &protocol.Array{Data: make([]protocol.Reply, elementCount)}
	// 往后再读相应个数
	for i := 0; i < elementCount; i++ {
		nextLine, err := reader.ReadBytes('\n')
		if err != nil {
			return err, nil
		}
		length := len(nextLine)
		if length <= 2 || nextLine[length-2] != '\r' {
			continue
		}
		nextLine = bytes.TrimSuffix(nextLine, []byte{'\r', '\n'})
		payLoad := getNextPayLoad(nextLine, reader)
		if payLoad.Err != nil {
			return payLoad.Err, nil
		}
		array.Data[i] = payLoad.Data
	}
	return nil, &PayLoad{Data: array}
}

// parseBulkString 解析多行字符串 line[1:] 为数组元素个数
func parseBulkString(line []byte, reader *bufio.Reader) (error, *PayLoad) {
	strLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return err, nil
	}
	buf := make([]byte, strLen)
	_, err = io.ReadFull(reader, buf)
	//mylog.Logger.Info("读取字符数 ", n, buf)
	_, _ = reader.ReadBytes('\n')
	if err != nil {
		return err, nil
	}
	return nil, &PayLoad{Data: &protocol.BulkString{Data: string(buf)}}
}
