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

type OK struct {
}

func (ss *OK) ToBytes() []byte {
	return []byte("+OK" + CRLF)
}
