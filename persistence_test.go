package xisdb

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {
	fmt.Println("TestReadFile")

	db, _ := Open(&Options{InMemory: false})
	db.Read(func(tx *Tx) error {
		for item := range tx.Each() {
			fmt.Println(item)
		}
		return nil
	})
}
