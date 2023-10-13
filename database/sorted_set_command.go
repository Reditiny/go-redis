package database

import (
	"redis-based-on-go/redis/protocol"
	"strconv"
)

// ZADD KEY_NAME SCORE1 VALUE1..SCOREN VALUEN
// zAdd 将一个或多个成员元素及其分数值加入到有序集当中  如果成员已存在则更新分数  并通过重新插入这个成员元素
func zAdd(db *DB, args [][]byte) protocol.Reply {
	if len(args) < 3 || len(args)%2 != 1 {
		return &protocol.ARGUEMENTS_NUMBER_ERR_REPLY
	}
	var set *sortedSet
	data, err := db.getData(string(args[0]))
	if err != nil {
		set = makeSortedSet()
	} else {
		var ok bool
		set, ok = data.(*sortedSet)
		if !ok {
			return &protocol.TYPE_ERR_REPLY
		}
	}
	insertCount := 0
	for i := 0; i < len(args)/2; i++ {
		score, err := strconv.Atoi(string(args[2*i+1]))
		key := string(args[2*i+2])
		if err != nil {
			return &protocol.NOT_INTEGER_ERR_REPLY
		}
		ele, ok := set.dict[key]
		if !ok {
			set.skipList.insert(key, float64(score))
			insertCount++

		} else {
			if ele.score != float64(score) {
				set.skipList.remove(key, ele.score)
				set.skipList.insert(key, float64(score))
			}
		}
		set.dict[string(args[2*i+2])] = &element{key: key, score: float64(score)}

	}
	db.setData(string(args[0]), set)
	return &protocol.Integer{Data: insertCount}
}

// ZRANGE KEY_NAME start end
// getSortedSetRange  获取有序集合中指定区间内的成员  其中成员的位置按分数值递增(从小到大)来排序
func getSortedSetRange(db *DB, args [][]byte) protocol.Reply {
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
		return &protocol.EMPTY_ARRAY
	} else {
		set := data.(*sortedSet)
		all := set.skipList.rangeScope(start, end, false)
		reply := make([]protocol.Reply, 2*len(all))
		for i, key := range all {
			reply[2*i] = &protocol.BulkString{Data: key.key}
			reply[2*i+1] = &protocol.BulkString{Data: strconv.Itoa(int(key.score))}
		}
		return &protocol.Array{Data: reply}
	}
}

func initSortedSetCommand() {
	commandTable["zadd"] = &command{name: "zadd", executor: zAdd}
	commandTable["zrange"] = &command{name: "zrange", executor: getSortedSetRange}
}
