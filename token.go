package vdf

type (
	// TokenType represents a lexical token.
	TokenType int
	// Token represents a token and value returned from the lexer.
	Token struct {
		Typ TokenType
		Val string
	}
)

// The list of tokens.
const (
	ILLEGAL TokenType = iota
	EOF
	IDENTIFIER
	LBRACE
	RBRACE
)
