package ql

// TokenType -- see constants below
type TokenType int

func (tt TokenType) String() string {
	for str, t := range tokenTypes {
		if tt == t {
			return str
		}
	}
	return "ILLEGAL"
}

func (tt TokenType) equalsTokenType(token *Token) bool {
	return tt.equals(token.tokenType)
}

func (tt TokenType) notEqualsTokenType(token *Token) bool {
	return tt.notEquals(token.tokenType)
}

func (tt TokenType) equals(other TokenType) bool {
	return tt == other
}

func (tt TokenType) notEquals(other TokenType) bool {
	return tt != other
}

// Token represents something in a query
type Token struct {
	raw       []byte
	tokenType TokenType
}

// NewToken creates a new token using the default raw value for the TokenType
func NewToken(t TokenType) Token {
	return Token{tokenIndex[t], t}
}

func (t Token) String() string {
	if t.raw == nil {
		return ""
	}
	return string(t.raw)
}

const (
	IDENTIFIER TokenType = iota

	SELECT
	FROM
	WHERE
	LIMIT
	USE
	INDEX
	BUCKET

	GET
	DEL
	SET
	EXISTS

	GT
	GTE
	LT
	LTE
	ASC
	DESC

	COMMA
	SEMICOLON
	COLON
	ASTERISK
	RPAREN
	LPAREN

	INTEGER

	// EOQ - End of Query
	EOQ
	ILLEGAL
)

var tokenTypes map[string]TokenType
var tokenIndex = make(map[TokenType][]byte)

func init() {
	tokenTypes = map[string]TokenType{
		"identifier": IDENTIFIER,

		"select": SELECT,
		"from":   FROM,
		"where":  WHERE,
		"limit":  LIMIT,

		"get":    GET,
		"del":    DEL,
		"set":    SET,
		"exsits": EXISTS,

		"use":    USE,
		"index":  INDEX,
		"bucket": BUCKET,

		"gt":   GT,
		"gte":  GTE,
		"lt":   LT,
		"lte":  LTE,
		"asc":  ASC,
		"desc": DESC,

		",": COMMA,
		";": SEMICOLON,
		":": COLON,
		"*": ASTERISK,
		"(": LPAREN,
		")": RPAREN,

		"INTEGER": INTEGER,

		"":        EOQ,
		"ILLEGAL": ILLEGAL,
	}

	for raw, token := range tokenTypes {
		tokenIndex[token] = []byte(raw)
	}
}

func getToken(token []byte) TokenType {
	t, ok := tokenTypes[string(token)]
	if !ok {
		return IDENTIFIER
	}
	return t
}
