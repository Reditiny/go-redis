package database

import (
	"errors"
	"fmt"
	"redis-based-on-go/redis/protocol"
	"time"
)

type DB struct {
	kVTable  *concurrentMap // key -> value  value 为任意 reids 类型
	ttlTable *concurrentMap // key -> ttl
}

func NewDB() *DB {
	return &DB{kVTable: makeConcurrentMap(), ttlTable: makeConcurrentMap()}
}

// setData 修改 db 的 key -> value value 为任意 redis 类型  内部修改由各结构本身实现
func (db *DB) setData(key string, value interface{}) int {
	_, ok := db.kVTable.get(key)
	if ok {
		db.kVTable.delete(key)
		db.kVTable.set(key, value)
		return 0
	} else {
		db.kVTable.set(key, value)
		return 1
	}
}

// setTTL 修改 db 的 key -> ttl
func (db *DB) setTTL(key string, ttl time.Duration) {
	db.ttlTable.set(key, time.Now().Add(ttl))
}

// getData 获取 指定 key 的 value
func (db *DB) getData(key string) (interface{}, error) {
	expire, ok := db.ttlTable.get(key)
	if ok {
		t := expire.(time.Time)
		if t.Before(time.Now()) {
			db.kVTable.delete(key)
			db.ttlTable.delete(key)
			return nil, errors.New(fmt.Sprintf("key '%s' has expired", key))
		}
	}
	value, ok := db.kVTable.get(key)
	if !ok {
		return nil, errors.New(fmt.Sprintf("key '%s' not found", key))
	} else {
		return value, nil
	}
}

// isExisted 判断指定 key 是否存在
func (db *DB) isExisted(key string) bool {
	_, ok := db.kVTable.get(key)
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
