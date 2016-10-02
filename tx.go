package xisdb

// Tx is a transaction against a database
type Tx struct {
	db *DB
}

func (tx *Tx) Put(key, value string) error {
	if tx.db == nil {
		return ErrorNoDatabase
	}

	tx.db.data[key] = value

	return nil
}

func (tx *Tx) Get(key string) (string, error) {
	if tx.db == nil {
		return "", ErrorNoDatabase
	}

	val, exists := tx.db.data[key]
	if !exists {
		return val, ErrorKeyNotFound
	}

	return val, nil
}
