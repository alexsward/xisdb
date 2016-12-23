package xisdb

import (
	"fmt"
	"testing"
)

// TestReadOnly -- ensures that things cannot be written to the DB
func TestReadOnly(t *testing.T) {
	fmt.Println("TestReadOnly")

	db, _ := Open(&Options{ReadOnly: true, InMemory: true})
	err := db.Set("key", "value")
	if err == nil {
		t.Error("Expected an error")
	}

	if err != ErrorDatabaseReadOnly {
		t.Errorf("Expected %s as an error, got %s", ErrorDatabaseReadOnly, err)
	}
}
