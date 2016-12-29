package xisdb

import (
	"github.com/alexsward/xisdb/indexes"
	"github.com/alexsward/xistree"
)

// This file contains xisdb functions that manage transactions for the user
// If you need complicated interactions, you need to use the db.Read() and
// db.ReadWrite() functions directly with a transaction function

// Bucket creates a bucket with the given name in the database, or no-op if it exists
func (db *DB) Bucket(name string) error {
	return db.ReadWrite(func(tx *Tx) error {
		_, err := tx.Bucket(name)
		return err
	})
}

// DeleteBucket will delete a bucket from the database, if it exists
func (db *DB) DeleteBucket(name string) error {
	return db.ReadWrite(func(tx *Tx) error {
		_, err := tx.DeleteBucket(name)
		return err
	})
}

// Get returns a value from the database
func (db *DB) Get(key string) (string, error) {
	var val string
	err := db.Read(func(tx *Tx) error {
		v, err := tx.Get(key)
		val = v
		return err
	})

	return val, err
}

// Exists will tell you if a key exists in the database
func (db *DB) Exists(key string) (bool, error) {
	does := false
	err := db.Read(func(tx *Tx) error {
		exists, err := tx.Exists(key)
		does = exists
		return err
	})
	return does, err
}

// Set adds/updates an object in the database
func (db *DB) Set(key, value string) error {
	return db.ReadWrite(func(tx *Tx) error {
		return tx.Set(key, value, nil)
	})
}

// Delete removes an object from the database
func (db *DB) Delete(key string) (bool, error) {
	var existed bool
	err := db.ReadWrite(func(tx *Tx) error {
		e, err := tx.Delete(key)
		existed = e
		return err
	})
	return existed, err
}

// AddIndex will add an index to the database
// Will match using the given Matcher and uses the xistree.Comparator function
func (db *DB) AddIndex(name string, it IndexType, m indexes.Matcher, c xistree.Comparator) error {
	if name == "" {
		return ErrInvalidIndexName
	}
	b := db.root()
	if _, exists := b.indexes[name]; exists {
		return ErrIndexAlreadyExists
	}

	return db.ReadWrite(func(tx *Tx) error {
		return tx.AddIndex(name, it, m, c)
	})
}

// DeleteIndex will remove an index by name, if it exists.
// Returns whether or not it was deleted along with any error that may have occured
func (db *DB) DeleteIndex(name string) (bool, error) {
	var deleted bool
	err := db.ReadWrite(func(tx *Tx) error {
		existed, err := tx.DeleteIndex(name)
		deleted = existed
		return err
	})
	return deleted, err
}
