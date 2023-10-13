package database

import "redis-based-on-go/redis/protocol"

// SADD KEY_NAME VALUE1..VALUEN
// addIntoSet 将若干个元素加入到集合中已存在的元素将被忽略  key 不存在则创建  当不是集合类型时返回一个错误
func addIntoSet(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 2 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	var set *Set
	data, err := db.getData(string(args[0]))
	if err != nil {
		set = makeSet()
	} else {
		var ok bool
		set, ok = data.(*Set)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
	}
	db.setData(string(args[0]), set)
	addCount := set.addKeys(args[1:])
	return &protocol.Integer{Data: addCount}
}

// SMEMBERS KEY_NAME
// getSetKeys 返回集合中的所有元素
func getSetKeys(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 1 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	data, err := db.getData(string(args[0]))
	if err != nil {
		return &protocol.EMPTY_ARRAY
	} else {
		set := data.(*Set)
		all := set.getAll()
		reply := make([]protocol.Reply, len(all))
		for i, key := range all {
			reply[i] = &protocol.BulkString{Data: key}
		}
		return &protocol.Array{Data: reply}
	}
}

func initSetCommand() {
	commandTable["sadd"] = &command{name: "sadd", executor: addIntoSet}
	commandTable["smembers"] = &command{name: "smembers", executor: getSetKeys}
}
