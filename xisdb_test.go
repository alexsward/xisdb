package xisdb

import (
	"fmt"
	"testing"
)

func openTestDB() *DB {
	db, _ := Open(&Options{
		InMemory:           true,
		BackgroundInterval: -1,
	})
	return db
}

// NOTE: many of these tests will fail unexpectedly as they just test high-level, not underlying
// Underlying APIs will be used that these tests will assume pass appropriately

func TestXisGet(t *testing.T) {
	fmt.Println("--- TestXisGet")
	tests := []struct {
		key, value string
		add        bool
		err        error
	}{
		{"key", "value", true, nil},
		{"key", "value", false, ErrKeyNotFound},
	}
	for i, test := range tests {
		db := openTestDB()
		if test.add {
			db.ReadWrite(func(tx *Tx) error {
				return tx.Set(test.key, test.value, nil)
			})
		}
		val, err := db.Get(test.key)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if val != test.value {
			t.Errorf("Test %d failed: expected value '%s', got '%s'", i+1, test.value, val)
		}
	}
}

func TestXisExists(t *testing.T) {
	fmt.Println("--- TestXisExists")
	tests := []struct {
		key         string
		add, exists bool
		err         error
	}{
		{"key", true, true, nil},
		{"key", false, false, nil},
	}
	for i, test := range tests {
		db := openTestDB()
		if test.add {
			db.ReadWrite(func(tx *Tx) error {
				return tx.Set(test.key, "test", nil)
			})
		}
		exists, err := db.Exists(test.key)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		if exists != test.exists {
			t.Errorf("Test %d failed: expected exists %t, got %t", i+1, test.exists, exists)
		}
	}
}
