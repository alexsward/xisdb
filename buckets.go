package xisdb

import (
	"fmt"
	"sync"
)

// Bucket is the user-facing representation of a bucket that enables transctions
type Bucket struct {
	tx      *Tx
	managed *bucket
}

// Get retrieves a value by its key, or errors
func (b *Bucket) Get(key string) (string, error) {
	return b.tx.get(b.managed, key)
}

// Exists returns whether or not a key is present in the Bucket
func (b *Bucket) Exists(key string) bool {
	exists, _ := b.tx.exists(b.managed, key)
	return exists
}

// Set will add or update a value
func (b *Bucket) Set(key, value string) error {
	return b.tx.set(b.managed, key, value, nil)
}

// Delete will delete a key from the bucket. Returns whether or not it actually was
func (b *Bucket) Delete(key string) (bool, error) {
	return b.tx.delete(b.managed, key)
}

// Clear will empty a bucket of all of its keys
func (b *Bucket) Clear() error {
	return b.tx.clear(b.managed)
}

// AddIndex adds an index to the bucket's data
func (b *Bucket) AddIndex() error {
	return nil
}

// DeleteIndex removes an index from the bucket's data
func (b *Bucket) DeleteIndex(name string) error {
	return nil
}

// Size is how many items are in the bucket
func (b *Bucket) Size() int {
	return b.managed.size()
}

// bucket is a collection of key-value pairs, much like a traditional DB table
type bucket struct {
	name    string
	db      *DB
	mutex   sync.RWMutex      // lock on a per-bucket level -- TODO: maybe not
	data    map[string]Item   // the data itself
	indexes map[string]*index // indexes on the data
}

func newBucket(name string, db *DB) *bucket {
	return &bucket{
		name:    name,
		db:      db,
		data:    make(map[string]Item),
		indexes: make(map[string]*index),
	}
}

func (b *bucket) isRoot() bool {
	return b.db != nil && b == b.db.root()
}

func (b *bucket) get(key string) (*Item, bool) {
	value, exists := b.data[key]
	if !exists {
		return nil, false
	}
	return &value, exists
}

func (b *bucket) insert(item *Item) {
	b.data[item.Key] = *item
}

func (b *bucket) exists(key string) bool {
	_, exists := b.data[key]
	return exists
}

// Delete removes a key from a bucket and returns whether or not it was removed
func (b *bucket) delete(key string) bool {
	item, ok := b.data[key]
	if !ok {
		return ok
	}

	delete(b.data, key)
	for _, idx := range b.indexes {
		if idx.match(&item) {
			idx.remove(&item)
		}
	}
	return ok
}

func (b *bucket) clear() error {
	for _, idx := range b.indexes {
		fmt.Printf("Idx:%s\n", idx)
		// idx.clear()
	}

	for key := range b.data {
		delete(b.data, key)
	}

	return nil
}

func (b *bucket) size() int {
	return len(b.data)
}

func (b *bucket) rollback(info *rollbackInfo) error {
	for key, value := range info.items {
		if value == nil {
			b.delete(key)
			continue
		}
		b.insert(value)
	}
	return nil
}
