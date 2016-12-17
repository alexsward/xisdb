package ql

import "errors"

var (
	// ErrIllegalToken when an unknown token is passed
	ErrIllegalToken = errors.New("Illegal token")
	// ErrEmptyQuery when you attempt to analyze a query of len=0
	ErrEmptyQuery = errors.New("Cannot parse an empty query")
	// ErrCannotParseStatement is a generic parse error
	ErrCannotParseStatement = errors.New("Cannot parse statement")
	// ErrUnsupportedStatement when the statement begins with something unknown
	ErrUnsupportedStatement = errors.New("Unknown satement")
	// ErrCannotFindIdentifiers when trying to pull IDENTIFIERs out of a statement
	ErrCannotFindIdentifiers = errors.New("Cannot extract desired identifiers from statement")
	// ErrInvalidBucketIdentifier when the bucket identifier is invalid
	ErrInvalidBucketIdentifier = errors.New("Invalid bucket identifier")
	// ErrNoBucketIdentifier when you attempt to select bucket without a bucket
	ErrNoBucketIdentifier = errors.New("Must provide a bucket identifier")
	// ErrCanOnlySelectSingleBucket when you attempt to select from multiple buckets
	ErrCanOnlySelectSingleBucket = errors.New("Can only select from a single bucket")
	// ErrNoLimitProvided when a LIMIT clause has no value
	ErrNoLimitProvided = errors.New("Limit clause requires a provided numeric value")
	// ErrLimitMustBeInteger when a LIMIT clause isn't followed by an integer
	ErrLimitMustBeInteger = errors.New("Limit value must be non-negative integer")
	// ErrUnparsedIdentifier when an identifier occurs that isn't expected, for example 'LIMIT 10 10'
	ErrUnparsedIdentifier = errors.New("An identifier wasn't accounted for")
	// ErrUnknownToken when something unknown and unaccounted for occurs
	ErrUnknownToken = errors.New("Unknown token in statement")
	// ErrNoIndexIdentifiers when a USE INDEX has no indexes
	ErrNoIndexIdentifiers = errors.New("No indexes provided")
	// ErrIllegalFromClause when the FROM clause is incorrect
	ErrIllegalFromClause = errors.New("Illegal FROM clause")
	// ErrIllegalUseClause when the USE clause is incorrect
	ErrIllegalUseClause = errors.New("Illegal USE clause")
	// ErrIncompleteStatement is a generic statement error
	ErrIncompleteStatement = errors.New("Incomplete statement")
	// ErrBothKeyValueRequired when a SET command doens't have a key and value
	ErrBothKeyValueRequired = errors.New("SET requires both key and value")
)
