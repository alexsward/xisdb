package xisdb

import (
	"time"

	"github.com/alexsward/xisdb/indexes"
	"github.com/alexsward/xisdb/tree"
)

// Tx is a transaction against a database
type Tx struct {
	id              int64                    // timestamp, in ns, of the transaction
	db              *DB                      // the database
	write           bool                     // if this is a write transaction
	rollbackBuckets map[string]*bucket       // buckets to rollback
	rollbacks       map[string]*rollbackInfo // how to roll back the entire transaction
	commits         map[string]*Item         // commit values
	hooks           []func()                 // functions to execute upon commit
	closed          bool
}

type rollbackInfo struct {
	items   map[string]*Item
	indexes map[string][]*Item
}

func newRollbackInfo() *rollbackInfo {
	return &rollbackInfo{
		items:   make(map[string]*Item),
		indexes: make(map[string][]*Item),
	}
}

// NewTransaction creates a new transaction against the DB
func NewTransaction(writeable bool, db *DB) *Tx {
	return &Tx{
		id:              time.Now().UnixNano(),
		db:              db,
		write:           writeable,
		rollbacks:       make(map[string]*rollbackInfo),
		rollbackBuckets: make(map[string]*bucket),
		commits:         make(map[string]*Item),
		hooks:           make([]func(), 0),
	}
}

func (tx *Tx) addRollback(bucket, key string, item *Item) {
	if !tx.write {
		return
	}

	if _, exists := tx.rollbacks[bucket]; !exists {
		tx.rollbacks[bucket] = newRollbackInfo()
	}

	if _, exists := tx.rollbacks[bucket].items[key]; exists {
		// this item has been added to be rolled back once, don't do it again
		// for example:
		//		db.Get("key") = "value"
		//    db.Set("key", "value1"), db.Set("key", "value2")
		//    --> don't rollback to value1
		return
	}
	tx.rollbacks[bucket].items[key] = item
}

func (tx *Tx) addRollbackBucket(bucket string, b *bucket) {
	if _, exists := tx.rollbackBuckets[bucket]; exists {
		// don't perform additional rollbacks for a bucket
		// delta from first change is what to roll back to
		return
	}
	tx.rollbackBuckets[bucket] = b
}

func (tx *Tx) addCommit(key string, item *Item) {
	tx.commits[key] = item
}

func (tx *Tx) close() {
	tx.db = nil
	tx.rollbacks = make(map[string]*rollbackInfo)
	tx.rollbackBuckets = make(map[string]*bucket)
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

// Bucket adds a bucket to the database by name
func (tx *Tx) Bucket(name string) (*Bucket, error) {
	if tx.db == nil {
		return nil, ErrNoDatabase
	}
	if !tx.write {
		return nil, ErrNotWriteTransaction
	}

	bucket, _ := tx.db.addBucket(name)
	b := &Bucket{
		tx:      tx,
		managed: bucket,
	}
	return b, nil
}

// DeleteBucket deletes a bucket from the database, if it exists. Returns whether or not it was deleted
func (tx *Tx) DeleteBucket(name string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}
	if !tx.write {
		return false, ErrNotWriteTransaction
	}

	return tx.db.deleteBucket(name)
}

// Buckets returns all buckets in the database. The root bucket will be first no matter what
func (tx *Tx) Buckets() ([]*Bucket, error) {
	var buckets []*Bucket
	if tx.db == nil {
		return buckets, ErrNoDatabase
	}
	if !tx.write {
		return buckets, ErrNotWriteTransaction
	}

	buckets = append(buckets, &Bucket{tx, tx.db.root()})
	for _, bucket := range tx.db.buckets {
		if bucket.isRoot() {
			continue
		}
		buckets = append(buckets, &Bucket{tx, bucket})
	}
	return buckets, nil
}

// Get retrieves a value from the database, if it exists
func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrNoDatabase
	}

	return tx.get(tx.db.root(), key)
}

func (tx *Tx) get(b *bucket, key string) (string, error) {
	item, exists := b.get(key)
	if !exists {
		return "", ErrKeyNotFound
	}
	return item.Value, nil
}

// Exists tells you if a key exists
func (tx *Tx) Exists(key string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}
	return tx.exists(tx.db.root(), key)
}

func (tx *Tx) exists(b *bucket, key string) (bool, error) {
	return b.exists(key), nil
}

// Set will add or update a key in the database
func (tx *Tx) Set(key, value string, md *SetMetadata) error {
	if tx.db == nil {
		return ErrNoDatabase
	}
	if !tx.write {
		return ErrNotWriteTransaction
	}

	return tx.set(tx.db.root(), key, value, md)
}

func (tx *Tx) set(b *bucket, key, value string, md *SetMetadata) error {
	var oldValue *Item
	if actual, exists := b.get(key); exists {
		oldValue = actual
	}
	tx.addRollback(b.name, key, oldValue)

	imd := &itemMetadata{}
	if md != nil && md.TTL > 0 {
		t := time.Now().Add(time.Millisecond * time.Duration(md.TTL))
		imd.expiration = &t
	}

	item := Item{key, value, imd}
	b.insert(&item)
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
	return tx.delete(tx.db.root(), key)
}

func (tx *Tx) delete(b *bucket, key string) (bool, error) {
	if !b.exists(key) {
		tx.rollbacks[key] = nil
		return false, ErrKeyNotFound
	}

	item, _ := b.get(key)
	tx.addRollback(b.name, key, item)
	tx.addCommit(key, nil)
	return b.delete(key), nil
}

func (tx *Tx) clear(b *bucket) error {
	if tx.db == nil {
		return ErrNoDatabase
	}
	if !tx.write {
		return ErrNotWriteTransaction
	}
	return b.clear()
}

// AddIndex creates a new index in the database using a read-write transaction
func (tx *Tx) AddIndex(name string, it IndexType, m indexes.Matcher, c tree.Comparator) error {
	if tx.db == nil {
		return ErrNoDatabase
	}

	if !tx.write {
		return ErrNotWriteTransaction
	}
	b := tx.db.root()
	_, exists := b.indexes[name]
	if exists {
		return ErrIndexAlreadyExists
	}

	idx, err := newIndex(name, it, m, c)
	if err != nil {
		return err
	}

	for _, value := range b.data {
		if !idx.match(&value) {
			continue
		}
		idx.add(&value)
	}
	b.indexes[name] = idx
	return nil
}

// DeleteIndex will delete an index from the database by name, if it exists
func (tx *Tx) DeleteIndex(name string) (bool, error) {
	if tx.db == nil {
		return false, ErrNoDatabase
	}
	b := tx.db.root()
	_, exists := b.indexes[name]
	if exists {
		delete(b.indexes, name)
	}
	return exists, nil
}

func (tx *Tx) iterate(indexName string, limit int) (<-chan Item, error) {
	if tx.db == nil {
		return nil, ErrNoDatabase
	}

	b := tx.db.root()
	idx, exists := b.indexes[indexName]
	if !exists {
		return nil, ErrIndexDoesNotExist
	}

	return idx.iterate(), nil
}
