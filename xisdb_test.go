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
