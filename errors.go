package xisdb

import "errors"

var (
	// ErrorNoDatabase - When a transaction doesn't have access to the database
	ErrorNoDatabase = errors.New("There's no database")

	// ErrorKeyNotFound when the key isn't found
	ErrorKeyNotFound = errors.New("Key not found")

	// ErrorDatabaseReadOnly when the database is read-only and a write operation is attempted
	ErrorDatabaseReadOnly = errors.New("Database is read only")

	// ErrorIncorrectDatabaseFileFormat when the database file has errors
	ErrorIncorrectDatabaseFileFormat = errors.New("Database file format is incorrect")

	// ErrorNotWriteTransaction when an update/write operation is attempted with a Read transaction
	ErrorNotWriteTransaction = errors.New("Cannot perform write operations in a read transaction")
)
