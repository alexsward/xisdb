package xisdb_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alexsward/xisdb"
)

// TestSetRollback -- tests an error creating a value
func TestSetRollback(t *testing.T) {
	fmt.Println("TestSetRollBack")

	db := openTestDB()
	db.ReadWrite(func(tx *xisdb.Tx) error {
		return tx.Set("key", "value", nil)
	})

	db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Set("key", "value2", nil)
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
		return tx.Set("key", "value", nil)
	})

	err := db.ReadWrite(func(tx *xisdb.Tx) error {
		tx.Set("key", "updatedValue", nil)
		return errors.New("Nope!")
	})

	if err == nil {
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
		t.Error("There should not have been an error thrown")
	}
}

// TestDeleteRollback -- tests rolling back a delete, transaction throws exception
func TestDeleteRollback(t *testing.T) {
	fmt.Println("TestDeleteRollback")

	db := openTestDB()
	db.ReadWrite(func(tx *xisdb.Tx) error {
		return tx.Set("key", "value", nil)
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

// TestExpiration -- tests that expiring a key works
func TestExpiration(t *testing.T) {
	db, _ := xisdb.Open(&xisdb.Options{
		InMemory:           true,
		BackgroundInterval: 10,
	})

	db.ReadWrite(func(tx *xisdb.Tx) error {
		err := tx.Set("expireme", "value", &xisdb.SetMetadata{10})
		return err
	})

	v, err := db.Get("expireme")
	if err != nil {
		t.Error(err)
	}

	if v != "value" {
		t.Errorf("Expected expireme key to have value [value], got [%s]", v)
	}

	time.Sleep(30 * time.Millisecond)

	_, err = db.Get("expireme")
	if err != xisdb.ErrorKeyNotFound {
		t.Errorf("Expected [%s] as error, got [%s]", xisdb.ErrorKeyNotFound, err)
	}
}
