package xisdb_test

import (
	"fmt"
	"testing"

	"github.com/alexsward/xisdb"
)

func TestReadFile(t *testing.T) {
	fmt.Println("TestReadFile")

	db, _ := xisdb.Open(&xisdb.Options{InMemory: false})
	db.Read(func(tx *xisdb.Tx) error {
		for item := range tx.Each() {
			fmt.Println(item)
		}
		return nil
	})
}

/*
func TestWrite(t *testing.T) {
	fmt.Println("TestWrite")
	if 1 < 2 {
		return
	}

	fmt.Println("TestWrite")
	db := openTestDB()
	db.Set("key1", "value1")
	db.Delete("key1")
	db.Set("key2", "value2")

	db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Set("key3", "value3")
		tx.Set("key4", "value4")
		tx.Set("key5", "value5")
		tx.Set("key6", "value6")
		return nil
	})
}
*/
