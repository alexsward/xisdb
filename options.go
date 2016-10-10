package xisdb

// Options represents configurable properties that initialize the database
type Options struct {
	// Filename is the location of the file to use, or to create
	Filename string

	// InMemory means whether to only save the data in memory
	InMemory bool

	// ReadOnly is to indicate this database is read-only
	ReadOnly bool

	// SkipDatabaseFileErrors will just pass over database file errors on load
	SkipDatabaseFileErrors bool
}
