package protocol

// OkReply 回复 +OK
var OkReply = new(OK)

// SimpleString 简单字符串 +7redTiny\r\n
type SimpleString struct {
	Data string
}

func (ss *SimpleString) ToBytes() []byte {
	return []byte("+" + ss.Data + CRLF)
}

func (ss *SimpleString) ToArgs() [][]byte {
	args := make([][]byte, 1)
	args[0] = []byte(ss.Data)
	return args
}

type OK struct {
}

func (ss *OK) ToBytes() []byte {
	return []byte("+OK" + CRLF)
}
