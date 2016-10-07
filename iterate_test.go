package xisdb_test

import (
	"fmt"
	"testing"

	"github.com/alexsward/xisdb"
)

func TestIterateEach(t *testing.T) {
	fmt.Println("TestIterateEach")

	db, _ := xisdb.Open(&xisdb.Options{})
	db.ReadWrite(func(tx *xisdb.Tx) error {
		for i := 1; i <= 10; i++ {
			k := fmt.Sprintf("key%d", i)
			v := fmt.Sprintf("value%d", i)
			tx.Set(k, v)
		}
		return nil
	})

	db.Read(func(tx *xisdb.Tx) error {
		for item := range tx.Each() {
			fmt.Println(item)
		}
		return nil
	})
}
