package ql

import (
	"fmt"
	"testing"
)

func TestLexerFirstChar(t *testing.T) {
	fmt.Println("-- TestLexerFirstChar")
	tests := []struct {
		err      error
		str      string
		expected byte
	}{
		{ErrEmptyQuery, "", '0'},
		{ErrEmptyQuery, " ", '0'},
		{nil, "select", *new(byte)},
		{nil, " select", *new(byte)},
	}
	for i, test := range tests {
		lexer, err := NewLexer(test.str)
		if err != test.err {
			t.Errorf("Test %d failed: Expected error %s, got %s", i+1, test.err, err)
			continue
		}
		if lexer != nil && test.expected != lexer.char {
			t.Errorf("Test %d failed: Expected byte: '%s', got '%s'", i+1, string(test.expected), string(lexer.char))
		}
	}
}

func TestLexerMatch(t *testing.T) {
	fmt.Println("-- TestLexerRead")
	tests := []struct {
		str                string
		fn                 func(byte) bool
		position, expected int
	}{
		{"select", isWhitespace, 0, 0}, {"select", isWhitespace, 5, 0}, {"select;", isWhitespace, 5, 0},
		{"select s", isWhitespace, 6, 1}, {"select  s", isWhitespace, 6, 2}, {"select   s", isWhitespace, 6, 3},
		{"select", isAlphanumeric, 0, 6}, {"select", isAlphanumeric, 5, 1}, {"select;", isAlphanumeric, 0, 6},
		{"select;", isAlphanumeric, 6, 0}, {"select; select", isAlphanumeric, 6, 0},
	}
	for i, test := range tests {
		lexer, _ := NewLexer(test.str)
		lexer.advance(test.position + 1)
		read := lexer.match(test.fn)
		if read != test.expected {
			t.Errorf("Test %d failed: expected to read %d chars, got %d", i+1, test.expected, read)
		}
	}
}

func TestLexerSingleTokens(t *testing.T) {
	fmt.Println("-- TestLexerSingleTokens")
	tests := []struct {
		str      string
		expected TokenType
	}{
		{"select", SELECT},
		{"from", FROM},
		{"where", WHERE},
		{"limit", LIMIT},
		{"use", USE},
		{"index", INDEX},
		{"gt", GT},
		{"gte", GTE},
		{"lt", LT},
		{"lte", LTE},
		{"asc", ASC},
		{"desc", DESC},
		{",", COMMA},
		{";", SEMICOLON},
		{":", COLON},
		{"*", ASTERISK},
		{"(", LPAREN},
		{")", RPAREN},
	}
	for i, test := range tests {
		lexer, _ := NewLexer(test.str)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Test %d failed, got an error, %s", i+1, err)
			continue
		}
		if len(tokens) != 1 {
			for _, tok := range tokens {
				fmt.Printf("tok:'%s' ", tok)
			}
			t.Errorf("Test %d failed, expected a single token, got %d tokens", i+1, len(tokens))
			continue
		}
		if tokens[0].tokenType != test.expected {
			t.Errorf("Test %d failed. Expected '%s' (%d) token, got %d", i+1, test.str, test.expected, tokens[0].tokenType)
		}
	}
}

func TestLexerMultipleTokens(t *testing.T) {
	fmt.Println("-- TestLexerMultipleTokens")
	tests := []struct {
		str      string
		expected []TokenType
	}{
		{
			str:      "select select select",
			expected: []TokenType{SELECT, SELECT, SELECT},
		},
		{
			str:      "select from where limit use index gt gte lt lte asc desc",
			expected: []TokenType{SELECT, FROM, WHERE, LIMIT, USE, INDEX, GT, GTE, LT, LTE, ASC, DESC},
		},
		{
			str:      "select;",
			expected: []TokenType{SELECT, SEMICOLON},
		},
	}
	for i, test := range tests {
		lexer, _ := NewLexer(test.str)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Test %d failed: got an error, %s", i+1, err)
			continue
		}
		if len(tokens) != len(test.expected) {
			t.Errorf("Test %d failed: expected %d tokens, got %d tokens", i+1, len(test.expected), len(tokens))
			continue
		}
		for idx, tok := range tokens {
			if tok.tokenType != test.expected[idx] {
				t.Errorf("Test %d failed: expected %d token to be '%s', got '%s'", i+1, idx, test.expected[idx], tok)
			}
		}
	}
}

func TestLexerNonKeywords(t *testing.T) {
	fmt.Println("-- TestLexerIdentifiers")
	tests := []struct {
		str      string
		err      error
		expected Token
	}{
		{"index1", nil, Token{[]byte("index1"), IDENTIFIER}},
		{"abc", nil, Token{[]byte("abc"), IDENTIFIER}},
		{"123", nil, Token{[]byte("123"), INTEGER}},
	}
	for i, test := range tests {
		lexer, _ := NewLexer(test.str)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("Test %d failed: got an error, %s", i+1, err)
			continue
		}
		if len(tokens) != 1 {
			t.Errorf("Test %d failed: expected 1 token, got %d tokens", i+1, len(tokens))
			continue
		}
	}
}
