package xisdb

import (
	"fmt"
	"sync"
	"time"
)

// DB is the data base object itself. It encapsulates all functionality for xisdb.
// Do not create an instance of this struct directly as you may introduce undesired
// side-effects through improper initialization.
type DB struct {
	mutex      sync.RWMutex // sync.RWMutex enables multiple read clients but only a single writer
	readOnly   bool
	persistent bool            // whether to persist to disk or not (not enabled currently)
	data       map[string]Item // the data itself
}

// Item is an item in the database, includes both the key and value of the object
type Item struct {
	Key, Value string
}

func (i Item) String() string {
	return fmt.Sprintf("key:[%s] value:[%s]", i.Key, i.Value)
}

// Open creates a new database
func Open(opts *Options) (*DB, error) {
	db := &DB{
		data:     make(map[string]Item),
		readOnly: opts.ReadOnly,
	}

	go db.run()

	return db, nil
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
	defer txn.close()

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
	db.hooks(tx)

	db.unlock(tx.write)
	return nil
}

func (db *DB) hooks(tx *Tx) {
	for _, fn := range tx.hooks {
		fn()
	}
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

// lock makes the database locked (uses RWMutex, so multiple readers available)
func (db *DB) lock(write bool) {
	if write {
		db.mutex.Lock()
		return
	}

	db.mutex.RLock()
}

// unlock makes the database accessible again
func (db *DB) unlock(write bool) {
	if write {
		db.mutex.Unlock()
		return
	}

	db.mutex.RUnlock()
}

// Close shuts down the database instance
func (db *DB) Close() error {
	return nil
}

// Read performs a read-only transaction against the database
func (db *DB) Read(fn func(tx *Tx) error) error {
	return db.execute(fn, false)
}

// ReadWrite performs a write-allowed transaction against the database
func (db *DB) ReadWrite(fn func(tx *Tx) error) error {
	if db.readOnly {
		return ErrorDatabaseReadOnly
	}

	return db.execute(fn, true)
}
