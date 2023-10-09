package protocol

// Error 错误 -ERR unknown command 'red'\r\n
type Error struct {
	Data string
}

func (e *Error) ToBytes() []byte {
	return []byte("-" + e.Data + CRLF)
}
