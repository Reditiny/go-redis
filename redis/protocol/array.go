package protocol

import "strconv"

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
