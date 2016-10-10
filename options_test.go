package xisdb_test

import (
	"fmt"
	"testing"

	"github.com/alexsward/xisdb"
)

// TestReadOnly -- ensures that things cannot be written to the DB
func TestReadOnly(t *testing.T) {
	fmt.Println("TestReadOnly")

	db, _ := xisdb.Open(&xisdb.Options{ReadOnly: true, InMemory: true})
	err := db.Set("key", "value")
	if err == nil {
		t.Error("Expected an error")
	}

	if err != xisdb.ErrorDatabaseReadOnly {
		t.Errorf("Expected %s as an error, got %s", xisdb.ErrorDatabaseReadOnly, err)
	}
}
