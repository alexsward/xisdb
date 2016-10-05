package xisdb_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/alexsward/xisdb"
)

// TestSetRollback -- tests an error creating a value
func TestSetRollback(t *testing.T) {
	fmt.Println("TestSetRollBack")

	db, _ := xisdb.Open(&xisdb.Options{})
	db.ReadWrite(func(tx *xisdb.Tx) error {
		return tx.Set("key", "value")
	})

	db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Set("key", "value2")
		return errors.New("Roll it back")
	})

	err := db.Read(func(tx *xisdb.Tx) error {
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

// TestSetUpdateRollback -- tests an error on update, gets rolled back
func TestSetUpdateRollback(t *testing.T) {
	fmt.Println("TestSetUpdateRollback")

	db, _ := xisdb.Open(&xisdb.Options{})
	db.ReadWrite(func(tx *xisdb.Tx) error {
		return tx.Set("key", "value")
	})

	err := db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Set("key", "updatedValue")
		return errors.New("Nope!")
	})

	if err != nil {
		t.Error("There should have been an error thrown")
	}

	err = db.Read(func(tx *xisdb.Tx) error {
		val, err := tx.Get("key")
		if err != nil {
			return err
		}

		if val != "value" {
			t.Error("Value was supposed to be value")
		}

		return nil
	})

	if err != nil {
		t.Error("There should have been an error thrown")
	}
}

// TestDeleteRollback -- tests rolling back a delete, transaction throws exception
func TestDeleteRollback(t *testing.T) {
	fmt.Println("TestDeleteRollback")

	db, _ := xisdb.Open(&xisdb.Options{})
	db.ReadWrite(func(tx *xisdb.Tx) error {
		return tx.Set("key", "value")
	})

	db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Delete("key")
		return errors.New("This is an error to cause a rollback")
	})

	db.Read(func(tx *xisdb.Tx) error {
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
