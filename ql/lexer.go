package ql

import "strings"

// Lexer is the lexer/tokenizer for
type Lexer struct {
	position int
	query    []byte
	char     byte
}

// NewLexer returns a new Lexer for the query after lowercasing it
func NewLexer(query string) (*Lexer, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if len(query) <= 0 {
		return nil, ErrEmptyQuery
	}
	l := &Lexer{
		query:    []byte(query),
		position: -1,
	}
	return l, nil
}

// Tokenize takes this lexer's query string and returns the tokens
func (l *Lexer) Tokenize() ([]*Token, error) {
	var tokens []*Token
	for l.next() {
		l.skip()
		switch l.char {
		case ',':
			tokens = appendToken(tokens, []byte{l.char}, COMMA)
		case ';':
			tokens = appendToken(tokens, []byte{l.char}, SEMICOLON)
		case ':':
			tokens = appendToken(tokens, []byte{l.char}, COLON)
		case '*':
			tokens = appendToken(tokens, []byte{l.char}, ASTERISK)
		case '(':
			tokens = appendToken(tokens, []byte{l.char}, LPAREN)
		case ')':
			tokens = appendToken(tokens, []byte{l.char}, RPAREN)
		default:
			start := l.position
			if isLetter(l.char) {
				ahead := l.match(isAlphanumeric)
				raw := l.query[start : start+ahead]
				tokens = append(tokens, &Token{raw, getToken(raw)})
				l.advance(ahead - 1)
			} else if isNumber(l.char) {
				ahead := l.match(isNumber)
				raw := l.query[start : start+ahead]
				tokens = append(tokens, &Token{raw, INTEGER})
				l.advance(ahead - 1)
			} else {
				tokens = append(tokens, &Token{[]byte{}, ILLEGAL})
				return tokens, ErrIllegalToken
			}
		}
	}
	return tokens, nil
}

// match tells you how many of the charcters, starting at l.position, match the predicate
func (l *Lexer) match(predicate func(byte) bool) int {
	i := 0
	for next, more := l.peekAt(l.position + i); more; next, more = l.peekAt(l.position + i) {
		if !predicate(next) {
			return i
		}
		i++
	}
	return i
}

// next moves the lexer ahead 1 character, if it's available
// sets position to current+1 and the current character to that next value
func (l *Lexer) next() bool {
	if b, more := l.peek(); more {
		l.position++
		l.char = b
		return true
	}
	return false
}

func (l *Lexer) advance(steps int) bool {
	for i := 0; i < steps; i++ {
		if !l.next() {
			return false
		}
	}
	return true
}

func (l *Lexer) peek() (byte, bool) {
	b, ok := l.peekAt(l.position + 1)
	return b, ok
}

func (l *Lexer) peekAt(i int) (byte, bool) {
	if i >= len(l.query) {
		return *new(byte), false
	}
	return l.query[i], true
}

// skip will move forward through whitspace as long as there are characters remaining
func (l *Lexer) skip() int {
	ahead := l.match(isWhitespace)
	l.advance(ahead)
	return ahead
}

func appendToken(tokens []*Token, raw []byte, t TokenType) []*Token {
	return append(tokens, &Token{raw, t})
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isAlphanumeric(ch byte) bool {
	return isLetter(ch) || isNumber(ch)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
