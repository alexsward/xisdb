package ql

import (
	"fmt"
	"testing"
)

func TestParserBasicFunctions(t *testing.T) {
	fmt.Println("-- TestParserBasicFunctions")
	p := NewParser(nil)
	p.tokens = []*Token{}
	p.position = -1
	if p.next() {
		t.Errorf("Did not expect a next() token")
	}
	p.tokens = append(p.tokens, &Token{[]byte{}, SELECT})
	peek, ok := p.peek()
	if peek == nil || !ok {
		t.Errorf("Expected non nil peeked Token: peek:%p, ok:%t\n", peek, ok)
		return
	}
	if peek.tokenType != SELECT {
		t.Errorf("Expected peeked token: %s, got %s", SELECT, peek.tokenType)
	}
	if !p.next() {
		t.Error("Expected a next token")
	}
	if p.position != 0 {
		t.Errorf("Expected position 0, got %d", p.position)
	}
}

func TestParserIsStatement(t *testing.T) {
	fmt.Println("-- TestParserIsStatement")
	tests := []struct {
		token    Token
		expected bool
	}{
		{Token{[]byte("select"), SELECT}, true},
		{Token{[]byte("get"), GET}, true},
		{Token{[]byte("set"), SET}, true},
		{Token{[]byte("del"), DEL}, true},

		{Token{[]byte("index"), INDEX}, false},
		{Token{[]byte("use"), USE}, false},
		{Token{[]byte("from"), USE}, false},
		{Token{[]byte("where"), USE}, false},
		{Token{[]byte("limit"), USE}, false},
		{Token{[]byte("GT"), GT}, false},
		{Token{[]byte("GTE"), GTE}, false},
		{Token{[]byte("LT"), LT}, false},
		{Token{[]byte("LTE"), LTE}, false},
		{Token{[]byte("ASC"), ASC}, false},
		{Token{[]byte("DESC"), DESC}, false},
	}

	for i, test := range tests {
		p := NewParser(nil)
		is := p.isStatement(&test.token)
		if is != test.expected {
			t.Errorf("Test %d failed: expected %t, got %t", i+1, test.expected, is)
		}
	}
}

func TestParserExtractIdentifiers(t *testing.T) {
	fmt.Println("-- TestParserExtractIdentifiers")
	tests := []struct {
		statement string
		err       error
		start     int
		ids       []string
	}{
		{"select from bucket bucket0 bucket1 limit 10", nil, 4, []string{"bucket0", "bucket1"}},
		{"select from bucket bucket0 bucket1 limit 10;", nil, 4, []string{"bucket0", "bucket1"}},
		{"get key", nil, 2, []string{"key"}},
		{"get key;", nil, 2, []string{"key"}},
		{"set key value", nil, 2, []string{"key", "value"}},
		{"set key value;", nil, 2, []string{"key", "value"}},
	}
	for i, test := range tests {
		l, _ := NewLexer(test.statement)
		p := NewParser(l)
		p.tokenize()
		p.advance(test.start)
		ids, err := p.extractIdentifiers()
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
			continue
		}
		if err != nil {
			continue
		}
		if len(ids) != len(test.ids) {
			t.Errorf("Test %d failed: Expected %d identifiers, got %d", i+1, len(test.ids), len(ids))
			continue
		}
		for j, id := range test.ids {
			raw := string(ids[j].raw)
			if id != raw {
				t.Errorf("Test %d failed: expected IDENTIFIER at position %d to be '%s', got '%s'", i+1, j, test.ids[j], raw)
			}
		}
	}
}

func TestParserGetStatement(t *testing.T) {
	fmt.Println("-- TestParserGetStatement")
	tests := []struct {
		statement string
		err       error
		expected  []string
	}{
		{"get key", nil, []string{"key"}},
		{"get key;", nil, []string{"key"}},
		{"get", ErrIncompleteStatement, []string{}},
		// {"get;", ErrIncompleteStatement, []string{}},
	}
	for i, test := range tests {
		s, err := parseSingleStatement(test.statement)
		if err != test.err {
			t.Errorf("Test %d failed: expected error:'%s', got:'%s'", i+1, test.err, err)
			continue
		}
		if test.err != nil {
			continue
		}
		statement, ok := s.(*GetStatement)
		if !ok {
			t.Errorf("Test %d failed: Expected a GetStatement, got a %T", i+1, s)
			continue
		}
		if len(statement.keys) != len(test.expected) {
			t.Errorf("Test %d failed: Expected %d keys, got %d", i+1, len(test.expected), len(statement.keys))
		}
		for j, key := range statement.keys {
			if key != test.expected[j] {
				t.Errorf("Test %d failed: Expected key at position %d: '%s' got '%s'", i+1, j, test.expected[j], key)
			}
		}
	}
}

func TestParserSetStatement(t *testing.T) {
	fmt.Println("-- TestParserSetStatement")
	tests := []struct {
		statement string
		err       error
		expected  map[string]string
	}{
		{"set key value;", nil, map[string]string{"key": "value"}},
		{"set key value", nil, map[string]string{"key": "value"}},
		{"set key value key2 value2", nil, map[string]string{"key": "value", "key2": "value2"}},
		{"set key value key2 value2;", nil, map[string]string{"key": "value", "key2": "value2"}},
		// {"set", ErrCannotFindIdentifiers, nil},
		// {"set;", ErrCannotFindIdentifiers, nil},
		{"set key", ErrBothKeyValueRequired, nil},
		{"set key;", ErrBothKeyValueRequired, nil},
	}
	for i, test := range tests {
		s, err := parseSingleStatement(test.statement)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
		}
		if test.err != nil {
			continue
		}
		statement, ok := s.(*SetStatement)
		if !ok {
			t.Errorf("Test %d failed: Expected a SetStatement, got a %T", i+1, s)
			continue
		}
		pairs := statement.Pairs()
		if len(pairs) != len(test.expected) {
			t.Errorf("Test %d failed: Expected %d items to be set, got %d", i+1, len(test.expected), len(pairs))
			continue
		}
		for key, value := range test.expected {
			v, ok := pairs[key]
			if !ok {
				t.Errorf("Test %d failed: Expected key %s to have value %s, it didn't exist", i+1, key, value)
				continue
			}
			if v != value {
				t.Errorf("Test %d failed: Expected key %s to have value %s, got %s", i+1, key, value, v)
			}
		}
	}
}

func TestParserDelStatement(t *testing.T) {
	fmt.Println("-- TestParserDelStatement")
	tests := []struct {
		statement string
		err       error
		ids       []string
	}{
		{"del key", nil, []string{"key"}},
		{"del key;", nil, []string{"key"}},
		{"del key1 key2", nil, []string{"key1", "key2"}},
		{"del key1 key2;", nil, []string{"key1", "key2"}},
	}
	for i, test := range tests {
		s, err := parseSingleStatement(test.statement)
		if err != test.err {
			t.Errorf("Test %d failed: expected error '%s', got '%s'", i+1, test.err, err)
		}
		if test.err != nil {
			continue
		}
		statement, ok := s.(*DelStatement)
		if !ok {
			t.Errorf("Test %d failed: Expected a DelStatement, got a %T", i+1, s)
			continue
		}
		for j, key := range statement.keys {
			if key != test.ids[j] {
				t.Errorf("Test %d failed: Expected key at position %d: '%s' got '%s'", i+1, j, test.ids[j], key)
			}
		}
	}
}

func TestParserSingleSelect(t *testing.T) {
	fmt.Println("-- TestParserSelect")
	tests := []struct {
		statement  string
		parseError error
		buckets    []string
		indexes    []string
		limit      int
	}{
		{"select from bucket bucket1 limit 10;", nil, []string{"bucket1"}, nil, 10},
		{"select from bucket bucket1 bucket2 limit 10;", nil, []string{"bucket1", "bucket2"}, nil, 10},
		{"select from bucket bucket1 limit;", ErrLimitMustBeInteger, nil, nil, 0},
		{"select from bucket bucket1 limit a;", ErrLimitMustBeInteger, nil, nil, 0},
		{"select from bucket bucket1 limit 17.3;", ErrIllegalToken, []string{"bucket1"}, nil, 10},
		{"select nothing", ErrUnparsedIdentifier, nil, nil, 0},
		{"select use index index1 limit 1;", nil, nil, []string{"index1"}, 1},
		{"select use index index1 index2 limit 1;", nil, nil, []string{"index1", "index2"}, 1},
		// {"select limit 10", nil, nil, nil, 10},
		{"select use", ErrCannotParseStatement, nil, nil, 0},
		{"select use", ErrCannotParseStatement, nil, nil, 0},
		{"select from", ErrCannotParseStatement, nil, nil, 0},
		{"select use limit 10;", ErrIllegalUseClause, nil, nil, 0},
		{"select from limit 10;", ErrIllegalFromClause, nil, nil, 0},
		{"select limit", ErrNoLimitProvided, nil, nil, 0},
		{"select from bucket;", ErrCannotFindIdentifiers, nil, nil, 0},
	}
	for i, test := range tests {
		l, _ := NewLexer(test.statement)
		p := NewParser(l)
		_, err := p.Parse()
		if err != test.parseError {
			t.Errorf("Test %d failed: Expected parse error:%s, got: %s", i+1, test.parseError, err)
			continue
		}
		if test.parseError != nil {
			continue // this makes the checks below too elaborate, just fail if it doesn't error properly
		}
		if len(p.selects) != 1 {
			t.Errorf("Test %d failed: Expected number of selects: 1, got %d", i+1, len(p.selects))
		}
		statement := p.selects[0]
		if len(statement.Buckets()) != len(test.buckets) {
			t.Errorf("Test %d failed: Expected %d buckets, got %d", i+1, len(test.buckets), len(statement.Buckets()))
			continue
		}
		for j, bucket := range test.buckets {
			if statement.Buckets()[j] != bucket {
				t.Errorf("Test %d failed: Expected bucket %d to be '%s', got '%s'", i+1, j+1, bucket, statement.Buckets()[j])
			}
		}
		if statement.Limit != test.limit {
			t.Errorf("Test %d failed: Expected limit %d, got %d", i+1, test.limit, statement.Limit)
		}
	}
}

func parseSingleStatement(statement string) (Statement, error) {
	l, _ := NewLexer(statement)
	s, err := NewParser(l).Parse()
	if len(s) > 0 {
		return s[0], err
	}
	return nil, err
}
