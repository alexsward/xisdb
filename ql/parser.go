package ql

import "strconv"

// Parse is a wrapper for the entire API of Lexer->Parser->[]Statement
func Parse(statement string) ([]Statement, error) {
	lexer, err := NewLexer(statement)
	if err != nil {
		return nil, err
	}
	return NewParser(lexer).Parse()
}

// Parser will parse the tokens into executable statements
type Parser struct {
	lexer    *Lexer
	tokens   []*Token
	position int

	selects []*SelectStatement
}

// NewParser creates a new Parser for the given Lexer
func NewParser(l *Lexer) *Parser {
	return &Parser{
		lexer:   l,
		selects: make([]*SelectStatement, 0),
	}
}

// Parse produces executable statements from the raw text
func (p *Parser) Parse() ([]Statement, error) {
	var statements []Statement
	err := p.tokenize()
	if err != nil {
		return statements, err
	}

	for p.next() {
		token := p.current()
		if p.isStatement(token) {
			switch token.tokenType {
			case SELECT:
				s, err := p.parseSelectStatement()
				if err != nil {
					return statements, err
				}
				statements = append(statements, s)
				p.selects = append(p.selects, s)
			case GET:
				s, err := p.parseGetStatement()
				if err != nil {
					return statements, err
				}
				statements = append(statements, s)
			case SET:
				s, err := p.parseSetStatement()
				if err != nil {
					return statements, err
				}
				statements = append(statements, s)
			case DEL:
				s, err := p.parseDelStatement()
				if err != nil {
					return statements, err
				}
				statements = append(statements, s)
			default:
				return statements, ErrCannotParseStatement
			}
		}
	}

	return statements, nil
}

func (p *Parser) tokenize() error {
	tokens, err := p.lexer.Tokenize()
	if err != nil {
		return err
	}
	p.tokens = tokens
	p.position = -1
	return nil
}

func (p *Parser) current() *Token {
	return p.tokens[p.position]
}

func (p *Parser) next() bool {
	if _, exists := p.peek(); exists {
		p.position++
		return true
	}
	return false
}

// advance will move the position forward n steps, or until the end
func (p *Parser) advance(steps int) bool {
	for i := 0; i < steps; i++ {
		if !p.next() {
			return false
		}
	}
	return true
}

// indexOfNext will return the index of the next token matching the predicate
func (p *Parser) indexOfNext(predicate func(*Token) bool) (int, bool) {
	i := p.position
	for next, more := p.peekAt(i); more; next, more = p.peekAt(i) {
		if predicate(next) {
			return i, true
		}
		i++
	}
	return -1, false
}

// indexOfLast returns the last index of the token matching predicate
func (p *Parser) indexOfLast(predicate func(*Token) bool) (int, bool) {
	match := false
	i := p.position
	for next, more := p.peekAt(i); more; next, more = p.peekAt(i) {
		if !predicate(next) {
			break
		}
		match = true
		i++
	}
	return i, match
}

func (p *Parser) peekAt(idx int) (*Token, bool) {
	if len(p.tokens) <= idx {
		return nil, false
	}
	return p.tokens[idx], true
}

func (p *Parser) peek() (*Token, bool) {
	return p.peekAt(p.position + 1)
}

var statements = []TokenType{SELECT, GET, DEL, SET, EXISTS}

// isStatement tells you if the token is specific to a given statement
func (p *Parser) isStatement(tok *Token) bool {
	for _, s := range statements {
		if s == tok.tokenType {
			return true
		}
	}
	return false
}

func (p *Parser) isEnd() bool {
	return p.current().tokenType == SEMICOLON || p.position == len(p.tokens)-1
}

// extractIdentifiers finds all IDENTIFIER tokens, starting at the current position
// it is the responsibility of the caller method to advance the parser to the first IDENTIFIER
// the reason for this is if there are LPAREN/RPAREN surrounding them, or similar uses
func (p *Parser) extractIdentifiers() ([]*Token, error) {
	last, ok := p.indexOfLast(IDENTIFIER.equalsTokenType)
	if !ok {
		return nil, ErrCannotFindIdentifiers
	}
	ids := make([]*Token, last-p.position)
	copy(ids, p.tokens[p.position:last])
	return ids, nil
}

// parseGetStatement generates a GetStatement
func (p *Parser) parseGetStatement() (*GetStatement, error) {
	if p.isEnd() {
		return nil, ErrIncompleteStatement
	}

	g := NewGetStatement()
	for p.next() || !p.isEnd() {
		switch p.current().tokenType {
		case IDENTIFIER:
			ids, err := p.extractIdentifiers()
			if err != nil {
				return g, err
			}
			for _, tok := range ids {
				g.addKeys(string(tok.raw))
			}
			p.advance(len(ids))
		case SEMICOLON:
			// return g, nil
		default:
			return g, ErrUnknownToken
		}
	}

	return g, nil
}

func (p *Parser) parseSetStatement() (*SetStatement, error) {
	s := NewSetStatement()
	for p.next() || !p.isEnd() {
		switch p.current().tokenType {
		case IDENTIFIER:
			kvps, err := p.extractIdentifiers()
			if err != nil {
				return s, err
			}
			if len(kvps)%2 != 0 {
				return s, ErrBothKeyValueRequired
			}
			for i := 0; i < len(kvps); i += 2 {
				key := string(kvps[i].raw)
				s.pairs[key] = string(kvps[i+1].raw)
			}
			p.advance(len(kvps))
		case SEMICOLON:
		default:
			return s, ErrUnknownToken
		}
	}
	return s, nil
}

func (p *Parser) parseDelStatement() (*DelStatement, error) {
	s := NewDelStatement()
	for p.next() {
		switch p.current().tokenType {
		case IDENTIFIER:
			ids, err := p.extractIdentifiers()
			if err != nil {
				return s, err
			}
			for _, tok := range ids {
				s.addKeys(string(tok.raw))
			}
			p.advance(len(ids))
		case SEMICOLON:
		default:
			return s, ErrUnknownToken
		}
	}
	return s, nil
}

func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	s := NewSelectStatement()
	for p.next() || !p.isEnd() {
		token := p.current()
		switch token.tokenType {
		case FROM:
			err := p.parseFrom(s)
			if err != nil {
				return s, err
			}
		case IDENTIFIER:
			return s, ErrUnparsedIdentifier
		case USE:
			err := p.parseUse(s)
			if err != nil {
				return s, err
			}
		case LIMIT:
			err := p.parseLimit(s)
			if err != nil {
				return s, err
			}
		case SEMICOLON:
			return s, nil
		default:
			return s, ErrUnknownToken
		}
	}
	return s, nil
}

// parseFrom parses everything from a FROM clause
func (p *Parser) parseFrom(s *SelectStatement) error {
	_, ok := p.peek()
	if !ok {
		return ErrCannotParseStatement
	}

	p.advance(1)
	switch p.current().tokenType {
	case BUCKET:
		advanced := p.advance(1) // we advance from BUCKET to the first identifier
		if !advanced || p.current().tokenType != IDENTIFIER {
			return ErrCannotFindIdentifiers
		}
		buckets, err := p.extractIdentifiers()
		if err != nil {
			return err
		}
		if len(buckets) == 0 {
			return ErrNoBucketIdentifier
		}
		s.addBuckets(buckets...)
		p.advance(len(buckets) - 1)
	default:
		return ErrIllegalFromClause
	}
	return nil
}

// parseUse parses everything from a USE clause
func (p *Parser) parseUse(s *SelectStatement) error {
	_, ok := p.peek()
	if !ok {
		return ErrCannotParseStatement
	}
	p.advance(1)
	switch p.current().tokenType {
	case INDEX:
		p.advance(1)
		indexes, err := p.extractIdentifiers()
		if err != nil {
			return err
		}
		if len(indexes) == 0 {
			return ErrNoIndexIdentifiers
		}

		s.addIndexes(indexes...)
		p.advance(len(indexes) - 1)
	default:
		return ErrIllegalUseClause
	}
	return nil
}

// parseLimit parses data out of a LIMIT clause
func (p *Parser) parseLimit(s *SelectStatement) error {
	if !p.next() {
		return ErrNoLimitProvided
	}
	if p.current().tokenType != INTEGER {
		return ErrLimitMustBeInteger
	}
	limit, err := strconv.Atoi(string(p.current().raw))
	if err != nil {
		return ErrLimitMustBeInteger
	}
	p.advance(1)
	s.Limit = limit
	return nil
}
