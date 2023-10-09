# redis-based-on-go
go 语言实现 redis 服务器
# Redis 简介
Redis（Remote Dictionary Server）是一种开源的内存数据库，是一个高性能、非关系型的数据存储系统。
# Redis 通信协议
Redis 服务器使用的协议是 Redis 协议，也称为 RESP（REdis Serialization Protocol）。 RESP 是二进制安全的。
是一种文本协议，用于客户端与 Redis 服务器之间的通信。RESP 协议设计得非常简单，使得它易于实现和解析，同时也很高效，适用于处理大量请求。
RESP 能序列化不同的数据类型，例如整型(integers)、字符串(strings)、数组(arrays)。额外还有特殊的错误类型。  
请求从客户端以字符串数组的形式发送到 redis 服务器，这些字符串表示要执行的命令的参数。Redis 用特定于命令的数据类型回复。
### 协议数据类型
+ 单行字符串 *Simple Strings*
  > 非二进制安全  
  > 首字节是 "+"  
  > 例如 "+OK/r/n" 表示字符串 "OK"
+ 错误 *Errors*   
  > 错误应答只在发生异常时发送  
  > 首字节是 "-"  
  > 例如 "-ERR unknown command 'redTiny'"
+ 整型 *Integers*
  > 整数返回也被用来返回布尔值
  > 首字节是 ":"  
  > 例如 ":42\r\n" 表示整数 42
+ 多行字符串 *Bulk Strings*
  > 二进制安全  
  > 首字节是 "\$"  
  > 例如 "$7\r\nredTiny\r\n" 表示字符串 "redTiny"   "$0\r\n\r\n"与"$-1\r\n" 表示空字符串
+ 数组 *Arrays*
  > 响应的首字节是 "*" 接着是表示数组中元素个数的十进制数 元素可以是任意RESP元素类型  
  > 例如 "*2\r\n:1\r\n$7\r\nredTiny\r\n" 表示数组 [1, "redTiny"]
### 示例
+ 设置键值对：
  > 请求：SET mykey Hello   *3\r\n$3\r\nSET\r\n$5\r\nmykey$5\r\nHello\r\n  
  > 响应：+OK\r\n
+ 获取键值：
  > 请求：GET mykey\r\n  
  > 响应：$5\r\nHello\r\n
+ 哈希操作：
  > 请求：HSET myhash field1 Hello\r\n  
  > 响应：+:1\r\n（表示设置成功） 
  > 请求：HGET myhash field1\r\n  
  > 响应：$5\r\nHello\r\n
+ 列表操作：
  > 请求：LPUSH mylist world\r\n  
  > 响应：+:1\r\n（表示插入成功）  
  > 请求：LRANGE mylist 0 -1\r\n  
  > 响应：*2\r\n$5\r\nworld\r\n$5\r\nHello\r\n




