package vdf

type (
	// tokenType represents a lexical token.
	tokenType int
	// token represents a token and value returned from the lexer.
	token struct {
		Typ tokenType
		Val string
	}
)

// The list of tokens.
const (
	ILLEGAL tokenType = iota
	EOF
	IDENTIFIER
	LBRACE
	RBRACE
)
