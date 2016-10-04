package xisdb

import (
	"errors"
	"fmt"
	"testing"
)

func TestSetRollback(t *testing.T) {
	fmt.Println("TestSetRollBack")

	db, _ := Open()
	db.ReadWrite(func(tx *Tx) error {
		return tx.Set("key", "value")
	})

	db.ReadWrite(func(tx *Tx) error {
		tx.Set("key", "value2")
		return errors.New("Roll it back")
	})

	err := db.Read(func(tx *Tx) error {
		val, err := tx.Get("key")
		if err != nil {
			return err
		}

		if val != "value" {
			return fmt.Errorf("Incorrect value, expected [value], got %s ", val)
		}

		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

// testDeleteRollback -- tests rolling back a delete, transaction throws exception
func TestDeleteRollback(t *testing.T) {
	fmt.Println("TestDeleteRollback")

	db, _ := Open()
	db.ReadWrite(func(tx *Tx) error {
		return tx.Set("key", "value")
	})

	db.ReadWrite(func(tx *Tx) error {
		tx.Delete("key")
		return errors.New("This is an error to cause a rollback")
	})

	db.Read(func(tx *Tx) error {
		val, err := tx.Get("key")
		if err != nil {
			return err
		}

		if val == "" {
			t.Errorf("key [key] was not found: %s", err)
		}
		return err
	})
}
