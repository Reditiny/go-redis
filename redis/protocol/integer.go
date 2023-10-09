package protocol

import "strconv"

// Integer 整型 :25\r\n
type Integer struct {
	Data int
}

func (i *Integer) ToBytes() []byte {
	return []byte(":" + strconv.Itoa(i.Data) + CRLF)
}
