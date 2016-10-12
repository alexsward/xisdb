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

	// DisableExpiration will turn off TTL keys, even if TTL is provided with a key
	DisableExpiration bool

	// BackgroundInterval (in ms) determines how frequently to perform background cleanup, < 0 means never, 0 defaults to 1000
	BackgroundInterval int
}
