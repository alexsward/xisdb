package xisdb

import (
	"fmt"
	"math/rand"
	"time"
)

// Tx is a transaction against a database
type Tx struct {
	id    uint32
	db    *DB
	write bool

	rollbacks map[string]*Item // rollback values
	hooks     []func()         // functions to execute upon commit
}

func (tx *Tx) String() string {
	return fmt.Sprintf("id:[%d] write:[%t]", tx.id, tx.write)
}

var randomer = rand.New(rand.NewSource(time.Now().UnixNano()))

func (tx *Tx) initialize(db *DB) {
	tx.db = db
	tx.id = randomer.Uint32()
	tx.rollbacks = make(map[string]*Item)
	tx.hooks = make([]func(), 0)
}

func (tx *Tx) close() {
	tx.db = nil
	tx.rollbacks = make(map[string]*Item)
	tx.hooks = make([]func(), 0)
}

// Set changes/adds a value to the database
func (tx *Tx) Set(key, value string) error {
	if tx.db == nil {
		return ErrorNoDatabase
	}
	tx.write = true

	var oldValue *Item
	if old, exists := tx.db.data[key]; exists {
		oldValue = &old
	}

	tx.rollbacks[key] = oldValue
	tx.db.data[key] = Item{key, value}

	return nil
}

// Delete removes a key entirely from the database, if it exists
func (tx *Tx) Delete(key string) (bool, error) {
	if tx.db == nil {
		return false, ErrorNoDatabase
	}
	tx.write = true

	val, exists := tx.db.data[key]
	if !exists {
		tx.rollbacks[key] = nil
		return false, nil
	}

	tx.rollbacks[key] = &val
	delete(tx.db.data, key)

	return true, nil
}

// Get retrieves a value from the database, if it exists
func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrorNoDatabase
	}

	if item, exists := tx.db.data[key]; exists {
		return item.Value, nil
	}

	return "", ErrorKeyNotFound
}

// Hooks adds post-commit hooks to this transaction
func (tx *Tx) Hooks(hooks ...func()) {
	tx.hooks = append(tx.hooks, hooks...)
}
