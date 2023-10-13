package database

import (
	"redis-based-on-go/redis/protocol"
	"strconv"
)

func getString(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 1 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.Error{Data: err.Error()}
	} else {
		return &protocol.BulkString{Data: data.(string)}
	}
}

func setString(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	return &protocol.Integer{Data: db.setData(string(args[0]), string(args[1]))}
}

// SETNX KEY_NAME VALUE
// setIfNotExists  指定的 key 不存在时为 key 设置指定的值
func setIfNotExists(db *DB, args [][]byte) protocol.Reply {
	existed := db.isExisted(string(args[0]))
	if existed {
		return &protocol.Integer{Data: 0}
	} else {
		return setString(db, args)
	}
}

// GETRANGE KEY_NAME start end
// getStringRange  获取存储在指定 key 中字符串的子字符串 字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)
func getStringRange(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 3 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return &protocol.NOT_INTEGER_ERR_REPLY
	}
	end, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return &protocol.NOT_INTEGER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.Error{Data: err.Error()}
	} else {
		str := data.(string)
		if end < start || end >= len(str) || start < 0 {
			start = 0
			end = len(str) - 1
		}
		return &protocol.BulkString{Data: str[start : end+1]}
	}
}

// MSET key1 value1 key2 value2 .. keyN valueN
// setStrings  同时设置一个或多个 key-value 对
func setStrings(db *DB, args [][]byte) protocol.Reply {
	if len(args)%2 != 0 || len(args) == 0 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	count := 0
	for i := 0; i < (len(args) / 2); i++ {
		count += db.setData(string(args[2*i]), string(args[2*i+1]))
	}
	return &protocol.Integer{Data: count}
}

func initStringCommand() {
	commandTable["get"] = &command{name: "get", executor: getString}
	commandTable["set"] = &command{name: "set", executor: setString}
	commandTable["setnx"] = &command{name: "setnx", executor: setIfNotExists}
	commandTable["getrange"] = &command{name: "getrange", executor: getStringRange}
	commandTable["mset"] = &command{name: "mset", executor: setStrings}
}
