package ql

import (
	"errors"
	"time"
)

// Statement represents something that will be executed against the DB
type Statement interface {
	// Validate ensures that this statement is executable
	Validate() error
	// Equals determines if two statements are equivalent
	Equals(Statement) bool
}

// GetStatement represents getting a specific number of items from the database
// get (key1, key2, key3) from bucket
type GetStatement struct {
	keys []string

	buckets []string
}

// NewGetStatement returns a new GetStatement object
func NewGetStatement() *GetStatement {
	return &GetStatement{}
}

// Keys gives you all keys this statement is asking for
func (s *GetStatement) Keys() []string {
	return s.keys
}

func (s *GetStatement) addKeys(keys ...string) {
	if s.keys == nil {
		s.keys = make([]string, 0)
	}

	for _, key := range keys {
		s.keys = append(s.keys, key)
	}
}

// Equals tells you if g is the same statement as other, semantically
func (s *GetStatement) Equals(other Statement) bool {
	return false
}

// Validate tells you if this is a valid GetStatement
func (s *GetStatement) Validate() error {
	if len(s.Keys()) == 0 {
		return errors.New("Cannot get nothing")
	}
	return nil
}

// KeyValuePair for SET commands. Potentially change to a struct{key,value}
type KeyValuePair map[string]string

// SetStatement is the SET statement/command
type SetStatement struct {
	pairs KeyValuePair
}

// NewSetStatement creates a new default SetStatement
func NewSetStatement() *SetStatement {
	return &SetStatement{
		pairs: make(KeyValuePair, 0),
	}
}

// Pairs returns all the key-value pairs being set by this statement
func (s *SetStatement) Pairs() KeyValuePair {
	return s.pairs
}

// Validate ensures at least one key-value pair is included for the statement
func (s *SetStatement) Validate() error {
	if len(s.Pairs()) == 0 {
		return errors.New("Cannot use set command with no key-value pairs")
	}
	return nil
}

// Equals determines if two statements are equivalent
func (s *SetStatement) Equals(other Statement) bool {
	return false
}

// DelStatement for deleting things from xisdb
type DelStatement struct {
	keys, buckets []string
}

// NewDelStatement creates a new defauly DEL statement
func NewDelStatement() *DelStatement {
	return &DelStatement{
		keys:    make([]string, 0),
		buckets: make([]string, 0),
	}
}

// Keys returns the keys to delete from the database
func (s *DelStatement) Keys() []string {
	return s.keys
}

// Buckets returns the buckets to delete from the database
func (s *DelStatement) Buckets() []string {
	return s.buckets
}

// Validate ensures we have at least one key or bucket
func (s *DelStatement) Validate() error {
	if len(s.Buckets()) == 0 && len(s.Keys()) == 0 {
		return errors.New("DEL statement cannot be executed on nothing")
	}
	return nil
}

// Equals determines if two statements are equivalent
func (s *DelStatement) Equals(other Statement) bool {
	return false
}

func (s *DelStatement) addKeys(keys ...string) {
	if s.keys == nil {
		s.keys = make([]string, 0)
	}

	for _, key := range keys {
		s.keys = append(s.keys, key)
	}
}

// SelectStatement represents asking for something from the database
// select from bucket bucket1 use index index1 limit 10;
type SelectStatement struct {
	Interrupt chan<- bool
	Max       time.Time

	buckets []string
	indexes []string
	Limit   int
}

// NewSelectStatement creates a new SelectStatement object
func NewSelectStatement() *SelectStatement {
	return &SelectStatement{
		Interrupt: make(chan<- bool),
	}
}

// Buckets returns this buckets used for this query
func (s *SelectStatement) Buckets() []string {
	return s.buckets
}

// Indexes returns the indexes for this query
func (s *SelectStatement) Indexes() []string {
	return s.indexes
}

// addBuckets adds all of the buckets to the SelectStatement
func (s *SelectStatement) addBuckets(buckets ...*Token) {
	if s.buckets == nil {
		s.buckets = make([]string, 0)
	}

	for _, bucket := range buckets {
		s.buckets = append(s.buckets, string(bucket.raw))
	}
}

// addIndexes adds all of the indexes to the select statement
func (s *SelectStatement) addIndexes(indexes ...*Token) {
	if s.indexes == nil {
		s.indexes = make([]string, 0)
	}

	for _, index := range indexes {
		s.indexes = append(s.indexes, string(index.raw))
	}
}

// Validate ensures that this statement is executable
func (s *SelectStatement) Validate() error {
	if len(s.Buckets()) > 1 {
		return ErrCanOnlySelectSingleBucket
	}
	return nil
}

// Equals tells you if the two statements are semantically the same
func (s *SelectStatement) Equals(other Statement) bool {
	return false
}
