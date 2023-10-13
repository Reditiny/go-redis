package database

import (
	"fmt"
	"redis-based-on-go/redis/protocol"
)

// HSET KEY_NAME FIELD VALUE
// setHash 为哈希表中的字段赋值  若哈希表不存在则创建并进行 HSET 操作  若字段已经存在于哈希表中则覆盖
func setHash(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 3 || len(args)%2 != 1 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	var hash *Hash
	existed := db.isExisted(string(args[0]))
	if existed {
		data, err := db.getData(string(args[0]))
		if err != nil {
			return &protocol.Error{Data: err.Error()}
		}
		var ok bool
		hash, ok = data.(*Hash)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
	} else {
		hash = makeSimpleHash()
	}
	insertCount := 0
	for i := 0; i < len(args)/2; i++ {
		insertCount += hash.set(string(args[2*i+1]), string(args[2*i+2]))
	}
	db.setData(string(args[0]), hash)
	return &protocol.Integer{Data: insertCount}
}

// HGET KEY_NAME FIELD_NAME
// getHash 返回哈希表中指定字段的值
func getHash(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.Error{Data: err.Error()}
	} else {
		hash, ok := data.(*Hash)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
		value, ok := hash.table[string(args[1])]
		if !ok {
			return &protocol.Error{Data: fmt.Sprintf("hash field '%s' not found", args[1])}
		} else {
			return &protocol.BulkString{Data: value}
		}
	}
}

// HGETALL KEY_NAME
// getAllField 返回哈希表中所有的字段和值  若 key 不存在则返回空列表
func getAllField(db *DB, args [][]byte) protocol.Reply {
	if len(args) != 1 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.EMPTY_ARRAY
	} else {
		hash, ok := data.(*Hash)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
		all := hash.getAll()
		replys := make([]protocol.Reply, len(all))
		for i, s := range all {
			replys[i] = &protocol.BulkString{Data: s}
		}
		return &protocol.Array{Data: replys}
	}
}

func initHashCommand() {
	commandTable["hset"] = &command{name: "hset", executor: setHash}
	commandTable["hget"] = &command{name: "hget", executor: getHash}
	commandTable["hgetall"] = &command{name: "hgetall", executor: getAllField}
}
