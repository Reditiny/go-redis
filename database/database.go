package database

import (
	"errors"
	"fmt"
	"redis-based-on-go/redis/protocol"
)

type DB struct {
	table *concurrentMap // key -> value  value 为任意 reids 类型
}

func NewDB() *DB {
	return &DB{table: makeConcurrentMap()}
}

// setData 修改 db 的 key -> value value 为任意 redis 类型  内部修改由各结构本身实现
func (db *DB) setData(key string, value interface{}) int {
	_, ok := db.table.get(key)
	if ok {
		db.table.delete(key)
		db.table.set(key, value)
		return 0
	} else {
		db.table.set(key, value)
		return 1
	}
}

// getData 获取 指定 key 的 value
func (db *DB) getData(key string) (interface{}, error) {
	value, ok := db.table.get(key)
	if !ok {
		return nil, errors.New(fmt.Sprintf("key '%s' not found", key))
	} else {
		return value, nil
	}
}

// isExisted 判断指定 key 是否存在
func (db *DB) isExisted(key string) bool {
	_, ok := db.table.get(key)
	return ok
}

func InitDBCommand() {
	initStringCommand()
	initHashCommand()
	initSetCommand()
	initListCommand()
	initSortedSetCommand()
}

// Execute 根据命令字段执行相应命令
func (db *DB) Execute(args [][]byte) protocol.Reply {
	command, ok := commandTable[string(args[0])]
	if !ok {
		return &protocol.Error{Data: fmt.Sprintf("command '%s' not found", args[0])}
	} else {
		return command.executor(db, args[1:])
	}
}
