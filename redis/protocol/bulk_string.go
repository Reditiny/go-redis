package protocol

import "strconv"

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
