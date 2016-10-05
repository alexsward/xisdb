package xisdb

// Options represents configurable properties that initialize the database
type Options struct {
	// ReadOnly is to indicate this database is read-only
	ReadOnly bool

	// DisableHooks disables custom post-commit hooks entirely on the database
	DisableHooks bool
}
