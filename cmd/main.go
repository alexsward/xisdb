package main

import (
	"fmt"

	"github.com/alexsward/xisdb"
)

func main() {
	db, err := xisdb.Open()
	if err != nil {
		fmt.Println("error opening database")
		return
	}

	err = db.Set(func(tx *xisdb.Tx) error {
		return tx.Put("derp", "is_derp")
	})

	fmt.Println(err)

	err = db.Get(func(tx *xisdb.Tx) error {
		v, err := tx.Get("derp")
		if err != nil {
			return err
		}

		fmt.Println(v)
		return nil
	})

	fmt.Println(err)
}
