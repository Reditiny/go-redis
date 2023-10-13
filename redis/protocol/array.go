package protocol

import "strconv"

var EMPTY_ARRAY = Array{}

// Array 数组 *3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n:25\r\n
type Array struct {
	Data []Reply
}

func (a *Array) ToBytes() []byte {
	str := "*" + strconv.Itoa(len(a.Data)) + CRLF
	for _, d := range a.Data {
		if d != nil {
			str += string(d.ToBytes())
		}
	}
	return []byte(str)
}

func (a *Array) ToArgs() [][]byte {
	args := make([][]byte, 0)
	for _, reply := range a.Data {
		for _, arg := range reply.ToArgs() {
			args = append(args, arg)
		}
	}
	return args
}
