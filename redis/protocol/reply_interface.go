package protocol

var (
	// CRLF redis 协议中的分隔符
	CRLF = "\r\n"
)

type Reply interface {
	ToBytes() []byte
}
