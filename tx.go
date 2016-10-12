package xisdb

import (
	"fmt"
	"time"
)

// Tx is a transaction against a database
type Tx struct {
	id        int64            //timestamp, in ns, of the transaction
	db        *DB              // the database
	write     bool             // if this is a write transaction
	rollbacks map[string]*Item // rollback values
	commits   map[string]*Item // commit values
	hooks     []func()         // functions to execute upon commit
}

func (tx *Tx) String() string {
	return fmt.Sprintf("id:[%d] write:[%t]", tx.id, tx.write)
}

func (tx *Tx) initialize(db *DB) {
	tx.db = db
	tx.id = time.Now().UnixNano()
	tx.rollbacks = make(map[string]*Item)
	tx.commits = make(map[string]*Item)
	tx.hooks = make([]func(), 0)
}

func (tx *Tx) close() {
	tx.db = nil
	tx.rollbacks = make(map[string]*Item)
	tx.commits = make(map[string]*Item)
	tx.hooks = make([]func(), 0)
}

// Get retrieves a value from the database, if it exists
func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrorNoDatabase
	}

	item := tx.db.get(key)
	if item != nil {
		return item.Value, nil
	}

	return "", ErrorKeyNotFound
}

// Set changes/adds a value to the database
func (tx *Tx) Set(key, value string, md *SetMetadata) error {
	if tx.db == nil {
		return ErrorNoDatabase
	}

	if !tx.write {
		return ErrorNotWriteTransaction
	}

	var oldValue *Item
	if tx.db.exists(key) {
		oldValue = tx.db.get(key)
	}

	tx.rollbacks[key] = oldValue

	imd := &itemMetadata{}
	if md != nil && md.TTL > 0 {
		t := time.Now().Add(time.Millisecond * time.Duration(md.TTL))
		imd.expiration = &t
	}

	item := Item{key, value, imd}
	tx.db.insert(&item)
	tx.commits[key] = &item

	return nil
}

// Delete removes a key entirely from the database, if it exists
func (tx *Tx) Delete(key string) (bool, error) {
	if tx.db == nil {
		return false, ErrorNoDatabase
	}

	if !tx.write {
		return false, ErrorNotWriteTransaction
	}

	if !tx.db.exists(key) {
		tx.rollbacks[key] = nil
		return false, ErrorKeyNotFound
	}

	item := tx.db.get(key)
	tx.rollbacks[key] = item
	tx.commits[key] = nil
	tx.db.remove(&Item{key, "", item.metadata})

	return true, nil
}

// Hooks adds post-commit hooks to this transaction
func (tx *Tx) Hooks(hooks ...func()) {
	tx.hooks = append(tx.hooks, hooks...)
}

// SetMetadata includes any additional non key-value parameters for setting a key
type SetMetadata struct {
	TTL int64
}
