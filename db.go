package xisdb

import (
	"fmt"
	"sync"
	"time"
)

// DB -- the database itself!
type DB struct {
	mutex sync.RWMutex
	data  map[string]string
}

func (db *DB) run() error {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		// perform anything that needs to be performed periodically
	}

	return nil
}

func (db *DB) execute(fn func(tx *Tx) error, write bool) error {
	txn := &Tx{}
	txn.initialize(db)

	db.lock(write)

	err := fn(txn)
	if err != nil {
		fmt.Printf("There was an error executing the transaction: %s\n", err)

		err = db.rollback(txn)
		if err != nil {
			fmt.Println("Error rolling back transaction")
		}
		return err
	}

	err = db.commit(txn)
	if err != nil {
		fmt.Println("Error committing transaction")
	}

	return err
}

func (db *DB) commit(tx *Tx) error {
	db.unlock(tx.write)
	return nil
}

func (db *DB) rollback(tx *Tx) error {
	for key, value := range tx.rollbacks {
		if value == nil {
			delete(db.data, key)
			continue
		}

		db.data[key] = *value
	}

	db.unlock(tx.write)
	return nil
}

func (db *DB) lock(write bool) {
	if write {
		db.mutex.Lock()
		return
	}

	db.mutex.RLock()
}

func (db *DB) unlock(write bool) {
	if write {
		db.mutex.Unlock()
		return
	}

	db.mutex.RUnlock()
}

// Close shuts down the database instnace
func (db *DB) Close() error {
	return nil
}

// Read performs a read-only transaction against the database
func (db *DB) Read(fn func(tx *Tx) error) error {
	return db.execute(fn, false)
}

// ReadWrite performs a write-allowed transaction against the database
func (db *DB) ReadWrite(fn func(tx *Tx) error) error {
	return db.execute(fn, true)
}

// Open creates a new database
func Open() (*DB, error) {
	db := &DB{
		data: make(map[string]string),
	}

	go db.run()

	return db, nil
}
