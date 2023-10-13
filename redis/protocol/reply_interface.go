package protocol

var (
	// CRLF redis 协议中的分隔符
	CRLF = "\r\n"
)

// Reply 客户端服务端通信载体接口
type Reply interface {
	ToBytes() []byte  // 转换为网络字节流
	ToArgs() [][]byte // 转换为命令字符串数组
}
