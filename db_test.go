package xisdb

import (
	"fmt"
	"testing"
	"time"
)

func TestDBRollbackItemAdd(t *testing.T) {
	fmt.Println("-- TestDBRollbackItemAdd")
	db := openTestDB()
	tx := NewTransaction(true, db)
	assertDBKeyValue(t, db, "key", "value", false)
	tx.addRollback("", "key", &Item{"key", "value", nil})
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", false)
}

func TestDBRollbackItemUpdate(t *testing.T) {
	fmt.Println("-- TestDBRollbackItemUpdate")
	db := openTestDB()
	tx := NewTransaction(true, db)
	db.Set("key", "value")
	assertDBKeyValue(t, db, "key", "value", true)
	tx.addRollback("", "key", &Item{"key", "value", nil})
	db.Set("key", "value2")
	assertDBKeyValue(t, db, "key", "value2", true)
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", true)
}

func TestDBRollbackItemDelete(t *testing.T) {
	fmt.Println("-- TestDBRollbackItemDelete")
	db := openTestDB()
	tx := NewTransaction(true, db)
	db.Set("key", "value")
	assertDBKeyValue(t, db, "key", "value", true)
	tx.addRollback("", "key", nil)
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", false)
}

func TestDBRollbackBucketAdd(t *testing.T) {
	fmt.Println("-- TestDBRollbackBucketAdd")
	db := openTestDB()
	tx := NewTransaction(true, db)
	assertBucketExists(t, db, "b1", false)
	db.Bucket("b1")
	assertBucketExists(t, db, "b1", true)
	tx.addRollbackBucket("b1", nil) // didn't exist before
	db.lock(true)
	err := db.rollback(tx)
	if err != nil {
		t.Errorf("Got an error rolling back: %s", err)
	}
	assertBucketExists(t, db, "b1", false)
}

func TestDBRollbackBucketDelete(t *testing.T) {
	fmt.Println("-- TestDBRollbackBucketDelete")
	db := openTestDB()
	tx := NewTransaction(true, db)
	assertBucketExists(t, db, "b1", false)
	db.Bucket("b1")
	assertBucketExists(t, db, "b1", true)
	b := db.buckets["b1"]
	tx.addRollbackBucket("b1", b)
	db.lock(true)
	err := db.rollback(tx)
	if err != nil {
		t.Errorf("Got an error rolling back: %s", err)
	}
	assertBucketExists(t, db, "b1", true)
}

func TestDBRollbackBucketMultipleOperations(t *testing.T) {
	fmt.Println("-- TestDBRollbackBucketMultipleOperations")
	db := openTestDB()
	tx := NewTransaction(true, db)
	assertBucketExists(t, db, "b1", false)
	db.Bucket("b1")
	tx.addRollbackBucket("b1", nil) // didn't exist
	assertBucketExists(t, db, "b1", true)
	b := db.buckets["b1"]
	db.DeleteBucket("b1")
	tx.addRollbackBucket("b1", b) // did
	assertBucketExists(t, db, "b1", false)
	if len(tx.rollbackBuckets) != 1 {
		t.Errorf("Expected 1 item in bucket rollback log, got %d", len(tx.rollbackBuckets))
	}
	db.Bucket("b1")
	tx.addRollbackBucket("b1", nil) // didn't exist again
	assertBucketExists(t, db, "b1", true)
	if len(tx.rollbackBuckets) != 1 {
		t.Errorf("Expected 1 item in bucket rollback log, got %d", len(tx.rollbackBuckets))
	}
	db.lock(true)
	err := db.rollback(tx)
	if err != nil {
		t.Errorf("Got an error rolling back: %s", err)
	}
	assertBucketExists(t, db, "b1", false)
}

func TestDBAddRootBucket(t *testing.T) {
	fmt.Println("-- TestDBAddRootBucket")
	db := openTestDB()
	b, created := db.addBucket("")
	if created {
		t.Errorf("Expected bucket '' to not be created, it was")
	}
	if !b.isRoot() {
		t.Errorf("Expected bucket '' to be root, it wasn't")
	}
}

func TestDBDeleteRootBucket(t *testing.T) {
	fmt.Println("-- TestDBDeleteRootBucket")
	db := openTestDB()
	ok, err := db.deleteBucket("")
	if !ok || err != ErrCannotDeleteRootBucket {
		t.Errorf("Shouldn't be able to delete root bucket. got err:%s", err)
	}
}

// TestCommitHooks -- verifies that functions execute upon commit
func TestDBCommitHooks(t *testing.T) {
	fmt.Println("-- TestDBCommitHooks")

	db := openTestDB()
	i := 0
	f := func() {
		i++
	}
	tx := NewTransaction(false, db)
	tx.Hooks(f)
	db.lock(false)
	db.commit(tx)
	if i != 1 {
		t.Errorf("Expected function to run, it did not")
	}
}

// TestDBExpiration -- tests that expiring a key works
func TestDBExpiration(t *testing.T) {
	fmt.Println("-- TestDBExpiration")
	db, _ := Open(&Options{
		InMemory:           true,
		BackgroundInterval: 10,
	})

	db.ReadWrite(func(tx *Tx) error {
		err := tx.Set("expireme", "value", &SetMetadata{10})
		return err
	})

	v, err := db.Get("expireme")
	if err != nil {
		t.Error(err)
	}

	if v != "value" {
		t.Errorf("Expected expireme key to have value [value], got [%s]", v)
	}

	time.Sleep(30 * time.Millisecond)

	_, err = db.Get("expireme")
	if err != ErrKeyNotFound {
		t.Errorf("Expected [%s] as error, got [%s]", ErrKeyNotFound, err)
	}
}

func TestDBBackground(t *testing.T) {
	fmt.Println("-- TestDBBackground")
}

func TestDBClose(t *testing.T) {
	fmt.Println("-- TestDBClose")
	db := openTestDB()
	err := db.Close()
	if err != nil {
		t.Errorf("Got an error closing database: %s", err)
	}
}

func assertDBKeyValue(t *testing.T, db *DB, key, value string, exists bool) {
	v, err := db.Get(key)
	if err == ErrKeyNotFound && exists {
		t.Errorf("Expected key '%s' to exist, got %s", key, err)
		return
	}
	if exists && v != value {
		t.Errorf("Expected value '%s', got '%s'", value, v)
	}
}

func assertBucketExists(t *testing.T, db *DB, bucket string, exists bool) {
	_, ok := db.buckets[bucket]
	if ok != exists {
		t.Errorf("Expected bucket to exist: %t, got %t", exists, !exists)
	}
}
