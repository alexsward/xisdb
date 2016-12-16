package xisdb

import (
	"fmt"
	"testing"
	"time"
)

func TestDBRollbackItemAdd(t *testing.T) {
	fmt.Println("--- TestDBRollbackItemAdd")
	db := openTestDB()
	tx := NewTransaction(true, db)
	assertDBKeyValue(t, db, "key", "value", false)
	tx.addRollback("key", &Item{"key", "value", nil})
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", false)
}

func TestDBRollbackItemUpdate(t *testing.T) {
	fmt.Println("--- TestDBRollbackItemUpdate")
	db := openTestDB()
	tx := NewTransaction(true, db)
	db.Set("key", "value")
	assertDBKeyValue(t, db, "key", "value", true)
	tx.addRollback("key", &Item{"key", "value", nil})
	db.Set("key", "value2")
	assertDBKeyValue(t, db, "key", "value2", true)
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", true)
}

func TestDBRollbackItemDelete(t *testing.T) {
	fmt.Println("--- TestDBRollbackItemDelete")
	db := openTestDB()
	tx := NewTransaction(true, db)
	db.Set("key", "value")
	assertDBKeyValue(t, db, "key", "value", true)
	tx.addRollback("key", nil)
	db.lock(true)
	db.rollback(tx)
	assertDBKeyValue(t, db, "key", "value", false)
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

// TestCommitHooks -- verifies that functions execute upon commit
func TestDBCommitHooks(t *testing.T) {
	fmt.Println("--- TestDBCommitHooks")

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
	fmt.Println("--- TestDBExpiration")
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
