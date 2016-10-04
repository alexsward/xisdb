package xisdb

import (
	"math/rand"
	"time"
)

// Tx is a transaction against a database
type Tx struct {
	id    uint32
	db    *DB
	write bool

	rollbacks map[string]*string
	commits   map[string]*string
	hooks     []func()
}

var randomer = rand.New(rand.NewSource(time.Now().UnixNano()))

func (tx *Tx) initialize(db *DB) {
	tx.db = db
	tx.id = randomer.Uint32()
	tx.rollbacks = make(map[string]*string)
}

// Set changes/adds a value to the database
func (tx *Tx) Set(key, value string) error {
	if tx.db == nil {
		return ErrorNoDatabase
	}
	tx.write = true

	var oldValue *string
	if old, exists := tx.db.data[key]; exists {
		oldValue = &old
	}

	tx.rollbacks[key] = oldValue
	tx.db.data[key] = value

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

	return false, nil
}

// Get retrieves a value from the database, if it exists
func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrorNoDatabase
	}

	return tx.db.data[key], nil
}
