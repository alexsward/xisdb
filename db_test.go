package xisdb

import (
	"fmt"
	"testing"
)

// TestCommitHooks -- verifies that functions execute upon commit
func TestCommitHooks(t *testing.T) {
	fmt.Println("TestCommitHooks")

	db, _ := Open()
	txn := &Tx{}
	txn.initialize(db)

	i := 0
	f := func() {
		i++
	}

	err := db.execute(func(tx *Tx) error {
		tx.Hooks(f)
		return nil
	}, false)

	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Expected increment of value to 1, instead value is [%d]", i)
	}
}
