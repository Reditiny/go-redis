package database

type Set struct {
	table map[string]interface{}
}

// SetConsumer 遍历函数
type SetConsumer func(key string) bool

// MakeSet makes a new set
func makeSet() *Set {
	return &Set{
		table: make(map[string]interface{}),
	}
}

// exist check if key exists in set
func (sh *Set) exist(key string) bool {
	_, ok := sh.table[key]
	return ok
}

// put key into set
func (sh *Set) put(key string) int {
	if sh.exist(key) {
		return 0
	}
	sh.table[key] = nil
	return 1
}

// addKeys add keys to set
func (sh *Set) addKeys(keys [][]byte) int {
	count := 0
	for _, k := range keys {
		count += sh.put(string(k))
	}
	return count
}

// traverse set
func (sh *Set) forEach(consumer SetConsumer) {
	for key := range sh.table {
		consumer(key)
	}
}

// get set's size
func (sh *Set) size() int {
	size := len(sh.table)
	return size
}

// get all keys in set
func (sh *Set) getAll() []string {
	all := make([]string, sh.size())
	i := 0
	sh.forEach(func(key string) bool {
		all[i] = key
		i++
		return true
	})
	return all
}
