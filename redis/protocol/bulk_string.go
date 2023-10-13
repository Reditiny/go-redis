package protocol

import "strconv"

var NIL_REPLY = NIL{}

// BulkString 多行字符串 $7\r\nredTiny\r\n
type BulkString struct {
	Data string
}

var nullBulkReply = []byte("$-1\r\n")

func (bs *BulkString) ToBytes() []byte {
	if bs == nil {
		return nullBulkReply
	}
	return []byte("$" + strconv.Itoa(len(bs.Data)) + CRLF + bs.Data + CRLF)
}

func (bs *BulkString) ToArgs() [][]byte {
	args := make([][]byte, 1)
	args[0] = []byte(bs.Data)
	return args
}

type NIL struct{}

func (n *NIL) ToBytes() []byte {
	return []byte("$-1\r\n")
}

func (n *NIL) ToArgs() [][]byte {
	return nil
}
