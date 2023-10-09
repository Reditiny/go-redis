package database

import (
	"errors"
	"sync"
)

// TODO 完善内存数据库与数据结构

type DB struct {
	table sync.Map
}

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Set(key, value string) int {
	_, ok := db.table.Load(key)
	if ok {
		db.table.Delete(key)
		db.table.Store(key, value)
		return 0
	} else {
		db.table.Store(key, value)
		return 1
	}
}

func (db *DB) Get(key string) (string, error) {
	value, ok := db.table.Load(key)
	if !ok {
		return "", errors.New("not found")
	} else {
		return value.(string), nil
	}
}
