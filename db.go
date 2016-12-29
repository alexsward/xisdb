package xisdb

import (
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
	mutex      sync.RWMutex       // sync.RWMutex enables multiple read clients but only a single writer
	persistent bool               // if false, do not persist to disk
	file       *os.File           // where to save the data
	fileErrors bool               // if loading a file should return an error
	readOnly   bool               // if this database is read-only
	bginterval int                // how often to perform background cleanup
	expires    bool               // if expiring keys are enabled
	buckets    map[string]*bucket // buckets
}

// Item is an item in the database, includes both the key and value of the object
type Item struct {
	Key, Value string
	metadata   *itemMetadata
}

type itemMetadata struct {
	expiration *time.Time
}

// Open creates a new database
func Open(opts *Options) (*DB, error) {
	db := &DB{
		readOnly:   opts.ReadOnly,
		persistent: !opts.InMemory,
		expires:    !opts.DisableExpiration,
		bginterval: opts.BackgroundInterval,
		buckets:    make(map[string]*bucket),
	}
	db.buckets[""] = newBucket("", db) // adding the rootBucket
	db.start()
	return db, nil
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
		return ErrDatabaseReadOnly
	}

	return db.execute(fn, true)
}

// addBucket will create a new bucket with the name, otherwise returns the existing
// returns the bucket and whether or not it was created
func (db *DB) addBucket(name string) (*bucket, bool) {
	if bucket, exists := db.buckets[name]; exists {
		return bucket, false
	}

	bucket := newBucket(name, db)
	db.buckets[name] = bucket
	return bucket, true
}

// deleteBucket removes a bucket from the database, returns if the bucket exists
// retruns an error if an attempt to remove the root bucket is made
func (db *DB) deleteBucket(name string) (bool, error) {
	bucket, exists := db.buckets[name]
	if !exists {
		return exists, nil
	}

	if bucket.isRoot() {
		return true, ErrCannotDeleteRootBucket
	}

	delete(db.buckets, name)
	return exists, nil
}

func (db *DB) root() *bucket {
	return db.buckets[""]
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
			buckets, err := tx.Buckets()
			if err != nil {
				return err
			}
			// TODO: potentially make expiration non-transactional
			now := time.Now().UnixNano()
			for _, bucket := range buckets {
				for _, item := range bucket.managed.data {
					if item.metadata != nil && item.metadata.expiration != nil {
						if now > item.metadata.expiration.Unix() {
							_, err := tx.delete(bucket.managed, item.Key)
							if err != nil {
								return err
							}
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
	db.lock(write) // TODO: defer db.unlock(write) ?
	defer txn.close()
	defer db.unlock(txn.write)

	err := fn(txn)
	if !write {
		// TODO: make this a slice?
		return firstNonNil(db.commit(txn), err)
	}

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
	db.hooks(tx)
	// persist to disk
	// pub-sub
	return nil
}

func (db *DB) hooks(tx *Tx) {
	for _, fn := range tx.hooks {
		fn()
	}
}

func (db *DB) rollback(tx *Tx) error {
	defer db.unlock(tx.write)
	if !tx.write {
		return ErrCannotRollbackReadTransaction
	}

	for name, bucket := range tx.rollbackBuckets {
		if bucket == nil {
			db.deleteBucket(name)
			continue
		}

		db.buckets[name] = bucket
	}

	for bucket, rollback := range tx.rollbacks {
		b, exists := db.buckets[bucket]
		if exists {
			err := b.rollback(rollback)
			if err != nil {
				return err
			}
		}
	}

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

func firstNonNil(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
