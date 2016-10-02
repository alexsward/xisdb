package xisdb

import "sync"

// DB -- the database itself!
type DB struct {
	lock sync.RWMutex
	data map[string]string
}

func (db *DB) execute(fn func(tx *Tx) error) error {
	txn := &Tx{}
	txn.db = db

	return fn(txn)
}

func (db *DB) Close() error {
	return nil
}

func (db *DB) Get(fn func(tx *Tx) error) error {
	return db.execute(fn)
}

func (db *DB) Set(fn func(tx *Tx) error) error {
	return db.execute(fn)
}

// Open creates a new database
func Open() (*DB, error) {
	db := &DB{}
	db.data = make(map[string]string)

	return db, nil
}
