package protocol

import "strconv"

// Integer 整型 :25\r\n
type Integer struct {
	Data int
}

func (i *Integer) ToBytes() []byte {
	return []byte(":" + strconv.Itoa(i.Data) + CRLF)
}

func (i *Integer) ToArgs() [][]byte {
	args := make([][]byte, 1)
	args[0] = []byte(strconv.Itoa(i.Data))
	return args
}
