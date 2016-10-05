package xisdb

import (
	"fmt"
	"testing"
)

// TestXisDBGet -- tests the higher-level Get API
func TestXisDBGet(t *testing.T) {
	fmt.Println("TestXisDBGet")

	db, _ := Open()
	db.ReadWrite(func(tx *Tx) error {
		return tx.Set("key", "value")
	})

	v, err := db.Get("key")
	if err != nil {
		t.Error(err)
	}

	if v != "value" {
		t.Errorf("Expected [value], got [%s]", v)
	}
}

// TestXisDBSet -- tests higher-level Set API
func TestXisDBSet(t *testing.T) {
	fmt.Println("TestXisDBGet")

	db, _ := Open()
	err := db.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	v, err := db.Get("key")
	if err != nil {
		t.Error(err)
	}

	if v != "value" {
		t.Errorf("Expected [value], got [%s]", v)
	}
}

// TestXisDBDelete -- tests high level delete API
func TestXisDBDelete(t *testing.T) {
	fmt.Println("TestXisDBDelete")

	db, _ := Open()
	err := db.Set("key", "value")
	if err != nil {
		t.Error(err)
	}

	removed, err := db.Delete("key")
	if err != nil {
		t.Error(err)
	}

	if !removed {
		t.Error("Expected to remove key, did not")
	}
}