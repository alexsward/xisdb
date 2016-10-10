package xisdb_test

import (
	"fmt"
	"testing"

	"github.com/alexsward/xisdb"
)

// TestCommitHooks -- verifies that functions execute upon commit
func TestCommitHooks(t *testing.T) {
	fmt.Println("TestCommitHooks")

	db := openTestDB()

	i := 0
	f := func() {
		i++
	}

	err := db.Read(func(tx *xisdb.Tx) error {
		tx.Hooks(f)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Expected increment of value to 1, instead value is [%d]", i)
	}
}
