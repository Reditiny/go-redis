package database

import "redis-based-on-go/redis/protocol"

// commandTable 记录所有客户端命令
var commandTable = make(map[string]*command)

// ExecFunc 命令执行函数
type ExecFunc func(db *DB, args [][]byte) protocol.Reply

type command struct {
	name     string
	executor ExecFunc
}
