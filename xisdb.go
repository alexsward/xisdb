package xisdb

// This file contains xisdb functions that manage transactions for the user
// If you need complicated interactions, you need to use the db.Read() and
// db.ReadWrite() functions directly with a transaction function

// Get returns a value from the database
func (db *DB) Get(key string) (string, error) {
	var val string
	err := db.Read(func(tx *Tx) error {
		v, err := tx.Get(key)
		val = v
		return err
	})

	return val, err
}

// Set adds/updates an object in the database
func (db *DB) Set(key, value string) error {
	return db.ReadWrite(func(tx *Tx) error {
		return tx.Set(key, value)
	})
}

// Delete removes an object from the database
func (db *DB) Delete(key string) (bool, error) {
	var existed bool
	err := db.ReadWrite(func(tx *Tx) error {
		e, err := tx.Delete(key)
		existed = e
		return err
	})
	return existed, err
}
