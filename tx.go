package xisdb

import (
	"fmt"
	"time"

	"github.com/alexsward/xisdb/indexes"
	"github.com/alexsward/xistree"
)

// Tx is a transaction against a database
type Tx struct {
	id              int64              // timestamp, in ns, of the transaction
	db              *DB                // the database
	write           bool               // if this is a write transaction
	rollbacks       map[string]*Item   // rollback values
	rollbackIndexes map[string][]*Item // index updates to roll back
	commits         map[string]*Item   // commit values
	hooks           []func()           // functions to execute upon commit
	closed          bool
}

// NewTransaction creates a new transaction against the DB
func NewTransaction(writeable bool, db *DB) *Tx {
	return &Tx{
		id:              time.Now().UnixNano(),
		db:              db,
		write:           writeable,
		rollbacks:       make(map[string]*Item),
		rollbackIndexes: make(map[string][]*Item),
		commits:         make(map[string]*Item),
		hooks:           make([]func(), 0),
	}
}

func (tx *Tx) String() string {
	return fmt.Sprintf("id:[%d] write:[%t]", tx.id, tx.write)
}

func (tx *Tx) addRollback(key string, item *Item) {
	tx.rollbacks[key] = item
}

func (tx *Tx) addCommit(key string, item *Item) {
	tx.commits[key] = item
}

func (tx *Tx) close() {
	tx.db = nil
	tx.rollbacks = make(map[string]*Item)
	tx.rollbackIndexes = make(map[string][]*Item)
	tx.commits = make(map[string]*Item)
	tx.hooks = make([]func(), 0)
	tx.closed = true
}

// Hooks adds post-commit hooks to this transaction
func (tx *Tx) Hooks(hooks ...func()) {
	tx.hooks = append(tx.hooks, hooks...)
}

// SetMetadata includes any additional non key-value parameters for setting a key
type SetMetadata struct {
	TTL int64
}

// Get retrieves a value from the database, if it exists
func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrNoDatabase
	}

	item := tx.db.get(key)
	if item != nil {
		return item.Value, nil
	}

	return "", ErrKeyNotFound
}

// Exists tells you if a key exists
func (tx *Tx) Exists(key string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}

	exists := tx.db.exists(key)
	return exists, nil
}

// Set changes/adds a value to the database
func (tx *Tx) Set(key, value string, md *SetMetadata) error {
	if tx.db == nil {
		return ErrNoDatabase
	}

	if !tx.write {
		return ErrNotWriteTransaction
	}

	var oldValue *Item
	if tx.db.exists(key) {
		oldValue = tx.db.get(key)
	}

	tx.addRollback(key, oldValue)

	imd := &itemMetadata{}
	if md != nil && md.TTL > 0 {
		t := time.Now().Add(time.Millisecond * time.Duration(md.TTL))
		imd.expiration = &t
	}

	item := Item{key, value, imd}
	tx.db.insert(&item)
	tx.addCommit(key, &item)

	return nil
}

// Delete removes a key entirely from the database, if it exists
func (tx *Tx) Delete(key string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}

	if !tx.write {
		return false, ErrNotWriteTransaction
	}

	if !tx.db.exists(key) {
		tx.rollbacks[key] = nil
		return false, ErrKeyNotFound
	}

	item := tx.db.get(key)
	tx.rollbacks[key] = item
	tx.addCommit(key, nil)
	tx.db.remove(&Item{key, "", item.metadata})

	return true, nil
}

// AddIndex creates a new index in the database using a read-write transaction
func (tx *Tx) AddIndex(name string, it IndexType, m indexes.Matcher, c xistree.Comparator) error {
	if tx.db == nil {
		return ErrNoDatabase
	}

	if !tx.write {
		return ErrNotWriteTransaction
	}

	_, exists := tx.db.indexes[name]
	if exists {
		return ErrIndexAlreadyExists
	}

	idx, err := newIndex(name, it, m, c)
	if err != nil {
		return err
	}

	for _, value := range tx.db.data {
		if !idx.match(&value) {
			continue
		}
		idx.add(&value)
	}
	tx.db.indexes[name] = idx
	return nil
}

// DeleteIndex will delete an index from the database by name, if it exists
func (tx *Tx) DeleteIndex(name string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}

	_, exists := tx.db.indexes[name]
	if exists {
		delete(tx.db.indexes, name)
	}
	return exists, nil
}

func (tx *Tx) iterate(indexName string, limit int) (<-chan Item, error) {
	if tx.db == nil {
		return nil, ErrNoDatabase
	}
	idx, exists := tx.db.indexes[indexName]
	if !exists {
		return nil, ErrIndexDoesNotExist
	}

	return idx.iterate(), nil
}
