package database

type Hash struct {
	table map[string]string
}

// HashConsumer 遍历函数
type HashConsumer func(key, value string) bool

// MakeSimple makes a new hash table
func makeSimpleHash() *Hash {
	return &Hash{
		table: make(map[string]string),
	}
}

func (sh *Hash) get(field string) (string, bool) {
	value, ok := sh.table[field]
	return value, ok
}

// 0 update  1 insert
func (sh *Hash) set(field, value string) int {
	_, ok := sh.table[field]
	sh.table[field] = value
	if ok {
		return 0
	} else {
		return 1
	}
}

func (sh *Hash) size() int {
	size := len(sh.table)
	return size
}

func (sh *Hash) getAll() []string {
	all := make([]string, sh.size()*2)
	i := 0
	sh.forEach(func(key, value string) bool {
		all[i] = key
		all[i+1] = value
		i += 2
		return true
	})
	return all
}

func (sh *Hash) forEach(consumer HashConsumer) {
	for key, value := range sh.table {
		consumer(key, value)
	}
}
