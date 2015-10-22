package pseudodb

import (
	"time"
)

type DB struct {
	rows       int
	starttime  time.Time
	changetime time.Time
	messages   []Row
}

type Row struct {
	ctime time.Time
	mtime time.Time
	atime time.Time
	value interface{}
}

func New() DB {
	var db = &DB{
		rows:       0,
		starttime:  time.Now(),
		changetime: time.Now(),
		messages:   []Row{},
	}
	return *db
}

func (db *DB) Insert(data interface{}) {
	db.rows++
	db.changetime = time.Now()
	db.messages = append(db.messages, Row{
		ctime: time.Now(),
		mtime: time.Now(),
		atime: time.Now(),
		value: data,
	})
}

/*
func (db *DB) Select(key, relop, value string) interface{} {
	return
}*/

func (db *DB) Count() int {
	return db.rows
}

func (db *DB) Each(handler func(*int, *interface{})) {
	for i, x := range db.messages {
		handler(&i, &x.value)
	}
}
