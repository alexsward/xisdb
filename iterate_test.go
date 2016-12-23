package xisdb

import (
	"fmt"
	"testing"
)

func TestIterateEach(t *testing.T) {
	fmt.Println("TestIterateEach")

	db := openTestDB()
	db.ReadWrite(func(tx *Tx) error {
		for i := 1; i <= 10; i++ {
			k := fmt.Sprintf("key%d", i)
			v := fmt.Sprintf("value%d", i)
			tx.Set(k, v, nil)
		}
		return nil
	})

	db.Read(func(tx *Tx) error {
		for item := range tx.Each() {
			fmt.Println(item)
		}
		return nil
	})
}
