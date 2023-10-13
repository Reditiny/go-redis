package database

import (
	"sync"
)

const (
	TABLE_COUNT = 16
	PRIME_32    = uint32(16777619)
)

// DB 中 key-value 的并发安全结构
type concurrentMap struct {
	tables      []*table
	memberCount int
	tableCount  int
}

type table struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

// makeConcurrentMap makes a new concurrent map
func makeConcurrentMap() *concurrentMap {
	cm := &concurrentMap{
		tables:     make([]*table, TABLE_COUNT),
		tableCount: TABLE_COUNT,
	}
	for i := 0; i < TABLE_COUNT; i++ {
		cm.tables[i] = &table{
			m: make(map[string]interface{}),
		}
	}
	return cm
}

// fnv32 hash function
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= PRIME_32
		hash ^= uint32(key[i])
	}
	return hash
}

// getTable returns the table of the given key
func (cm *concurrentMap) getTable(key string) *table {
	return cm.tables[fnv32(key)%uint32(cm.tableCount)]
}

// set 并发安全地修改 db 的 key -> value
func (cm *concurrentMap) set(key string, value interface{}) int {
	table := cm.getTable(key)
	table.mutex.Lock()
	defer table.mutex.Unlock()
	_, ok := table.m[key]
	if ok {
		table.m[key] = value
		return 0
	} else {
		table.m[key] = value
		cm.memberCount++
		return 1
	}
}

// get 并发安全地获取指定 key 的 value
func (cm *concurrentMap) get(key string) (interface{}, bool) {
	table := cm.getTable(key)
	table.mutex.RLock()
	defer table.mutex.RUnlock()
	value, ok := table.m[key]
	if !ok {
		return nil, false
	} else {
		return value, true
	}
}

// isExisted 并发安全地判断指定 key 是否存在
func (cm *concurrentMap) isExisted(key string) bool {
	table := cm.getTable(key)
	table.mutex.RLock()
	defer table.mutex.RUnlock()
	_, ok := table.m[key]
	return ok
}

// delete 并发安全地删除指定 key
func (cm *concurrentMap) delete(key string) int {
	table := cm.getTable(key)
	table.mutex.Lock()
	defer table.mutex.Unlock()
	_, ok := table.m[key]
	if ok {
		delete(table.m, key)
		cm.memberCount--
		return 1
	} else {
		return 0
	}
}
