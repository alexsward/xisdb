package xisdb

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	defaultFilename = "xisdb.data"
)

// DB is the data base object itself. It encapsulates all functionality for xisdb.
// Do not create an instance of this struct directly as you may introduce undesired
// side-effects through improper initialization.
type DB struct {
	mutex         sync.RWMutex            // sync.RWMutex enables multiple read clients but only a single writer
	persistent    bool                    // if false, do not persist to disk
	file          *os.File                // where to save the data
	fileErrors    bool                    // if loading a file should return an error
	readOnly      bool                    // if this database is read-only
	bginterval    int                     // how often to perform background cleanup
	expires       bool                    // if expiring keys are enabled
	data          map[string]Item         // the data itself
	indexes       map[string]*index       // indexes on the data
	subscriptions map[string]subscription // subscriptions
}

// Item is an item in the database, includes both the key and value of the object
type Item struct {
	Key, Value string
	metadata   *itemMetadata
}

type itemMetadata struct {
	expiration *time.Time
}

func (i Item) String() string {
	return fmt.Sprintf("key:[%s] value:[%s]", i.Key, i.Value)
}

// Open creates a new database
func Open(opts *Options) (*DB, error) {
	db := &DB{
		data:          make(map[string]Item),
		indexes:       make(map[string]*index),
		subscriptions: make(map[string]subscription),
		readOnly:      opts.ReadOnly,
		persistent:    !opts.InMemory,
		expires:       !opts.DisableExpiration,
		bginterval:    opts.BackgroundInterval,
	}

	if db.persistent {
		filename := opts.Filename
		if filename == "" {
			filename = defaultFilename
		}

		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, err
		}
		db.file = f

		err = db.load()
		if err != nil {
			return nil, err
		}
	}

	db.start()

	return db, nil
}

// Close shuts down the database instance
func (db *DB) Close() error {
	for _, s := range db.subscriptions {
		for _, ch := range s.channels {
			close(ch)
		}
	}

	return nil
}

// Read performs a read-only transaction against the database
func (db *DB) Read(fn func(tx *Tx) error) error {
	return db.execute(fn, false)
}

// ReadWrite performs a write-allowed transaction against the database
func (db *DB) ReadWrite(fn func(tx *Tx) error) error {
	if db.readOnly {
		return ErrDatabaseReadOnly
	}

	return db.execute(fn, true)
}

func (db *DB) start() {
	if db.bginterval < 0 {
		go func() {
			select {}
		}()
	} else {
		if db.bginterval == 0 {
			db.bginterval = 1000
		}

		go db.background()
	}
}

// background performs background tasks, like cleanp of TTL keys
// TTL cleanup happens in a transaction, so pubsub and persistence and everything else
// takes place with the expirations as well
func (db *DB) background() error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(db.bginterval))
	defer ticker.Stop()
	for range ticker.C {
		if !db.expires {
			continue
		}

		err := db.ReadWrite(func(tx *Tx) error {
			now := time.Now().UnixNano()
			for _, item := range db.data {
				if item.metadata != nil && item.metadata.expiration != nil {
					if now > item.metadata.expiration.Unix() {
						_, err := tx.Delete(item.Key)
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) execute(fn func(tx *Tx) error, write bool) error {
	txn := NewTransaction(write, db)
	defer txn.close()

	db.lock(write) // TODO: defer db.unlock(write) ?

	err := fn(txn)
	if err != nil {
		rollbackErr := db.rollback(txn)
		if rollbackErr != nil {
			err = rollbackErr
		}
		return err
	}

	err = db.commit(txn)
	if err != nil {
		err = db.rollback(txn) // no idea how to handle an error here...
	}

	return err
}

func (db *DB) commit(tx *Tx) error {
	if db.persistent {
		err := tx.persist()
		if err != nil {
			return err
		}
	}

	db.hooks(tx)

	if len(db.subscriptions) > 0 {
		var items []Item
		for _, item := range tx.commits {
			items = append(items, *item)
		}
		db.publish(items...)
	}

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

func (db *DB) get(key string) *Item {
	if item, exists := db.data[key]; exists {
		return &item
	}

	return nil
}

func (db *DB) exists(key string) bool {
	_, exists := db.data[key]
	return exists
}

func (db *DB) insert(item *Item) {
	db.data[item.Key] = *item
}

func (db *DB) remove(item *Item) {
	delete(db.data, item.Key)
}

func (db *DB) hasIndex(name string) bool {
	_, exists := db.indexes[name]
	return exists
}

func (db *DB) index(item *Item) []string {
	var indexes []string
	for name, index := range db.indexes {
		if index.match(item) {
			index.add(item)
			indexes = append(indexes, name)
		}
	}
	return indexes
}
