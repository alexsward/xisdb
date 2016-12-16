package indexes

import (
	"regexp"
	"strings"
)

// Matcher is a function that determines if a string matches
type Matcher func(string) bool

// PrefixMatcher creates a matcher that matches a prefix to a string
func PrefixMatcher(prefix string) Matcher {
	m := func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
	return m
}

// WildcardMatcher matches all strings
var WildcardMatcher = func(s string) bool {
	return true
}

// RegexMatcher returns a regular-expression matcher for the given regex
func RegexMatcher(regex string) (Matcher, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	m := func(s string) bool {
		return r.MatchString(s)
	}
	return m, err
}
