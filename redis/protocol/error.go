package protocol

var ARGUEMENTS_NUMBER_ERR_REPLY = Error{Data: "wrong number of arguments"}
var NOT_INTEGER_ERR_REPLY = Error{Data: "value is not an integer"}
var TYPE_ERR_REPLY = Error{Data: "WRONGTYPE Operation against a key holding the wrong kind of value"}

// Error 错误 -ERR unknown command 'red'\r\n
type Error struct {
	Data string
}

func (e *Error) ToBytes() []byte {
	return []byte("-" + e.Data + CRLF)
}

func (e *Error) ToArgs() [][]byte {
	args := make([][]byte, 1)
	args[0] = []byte(e.Data)
	return args
}
