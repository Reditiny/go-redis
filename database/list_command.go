package database

import (
	"redis-based-on-go/redis/protocol"
	"strconv"
)

// RPUSH KEY_NAME VALUE1..VALUEN
// rPush 将一个或多个值插入到列表的尾部(最右边)  若 key 不存在则创建  若 key 对应的值不是列表类型则返回错误
func rPush(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	var list *List
	data, err := db.getData(string(args[0]))
	if err != nil {
		list = makeList()
	} else {
		var ok bool
		list, ok = data.(*List)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
	}
	db.setData(string(args[0]), list)
	pushCount := list.rPush(args[1:])
	return &protocol.Integer{Data: pushCount}
}

// LINDEX KEY_NAME INDEX
// lIndex 通过索引获取列表中的元素  如果指定索引值不在列表的区间范围内，返回 nil  下标从 0 开始
func lIndex(db *DB, args [][]byte) protocol.Reply {
	if len(args) != 2 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.Error{Data: err.Error()}
	} else {
		list, ok := data.(*List)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
		index, err := strconv.Atoi(string(args[1]))
		if err != nil {
			return &protocol.Error{Data: err.Error()}
		}
		node := list.get(index)
		if node == nil {
			return &protocol.NIL_REPLY
		}
		return &protocol.BulkString{Data: node.value}
	}
}

func initListCommand() {
	commandTable["rpush"] = &command{name: "rpush", executor: rPush}
	commandTable["lindex"] = &command{name: "lindex", executor: lIndex}
}
