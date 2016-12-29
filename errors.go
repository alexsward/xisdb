package xisdb

import "errors"

var (
	// ErrNoDatabase - When a transaction doesn't have access to the database
	ErrNoDatabase = errors.New("There's no database")

	// ErrKeyNotFound when the key isn't found
	ErrKeyNotFound = errors.New("Key not found")

	// ErrDatabaseReadOnly when the database is read-only and a write operation is attempted
	ErrDatabaseReadOnly = errors.New("Database is read only")

	// ErrIncorrectDatabaseFileFormat when the database file has errors
	ErrIncorrectDatabaseFileFormat = errors.New("Database file format is incorrect")

	// ErrNotWriteTransaction when an update/write operation is attempted with a Read transaction
	ErrNotWriteTransaction = errors.New("Cannot perform write operations in a read transaction")

	// ErrInvalidIndexName when an index name is invalid
	ErrInvalidIndexName = errors.New("Index name is invalid")

	// ErrIndexAlreadyExists when a duplicate index is attempted to be added to the database
	ErrIndexAlreadyExists = errors.New("Index already exists")

	// ErrIndexDoesNotExist when an index doesn't exist, usually when attempting to iterate/scan
	ErrIndexDoesNotExist = errors.New("Index doesn't exist")

	// ErrCannotDeleteRootBucket when an attempt to delete the root bucket is made
	ErrCannotDeleteRootBucket = errors.New("Cannot delete root bucket")

	// ErrCannotRollbackReadTransaction when you try and roll back a read-only transaction
	ErrCannotRollbackReadTransaction = errors.New("Read-only transactions cannot be rolled back")
)
