package xisdb

import "errors"

var (
	// ErrorNoDatabase - When a transaction doesn't have access to the database
	ErrorNoDatabase = errors.New("There's no database")

	// ErrorKeyNotFound when the key isn't found
	ErrorKeyNotFound = errors.New("Key not found")
)
