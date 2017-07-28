package vdf

import (
	"log"
)

type parser struct {
	lex       *lexer
	token     [3]token
	peekCount int
}

// -------------------------------------------------------

// Parse returns a map of parsed data
func Parse(data string) map[string]interface{} {
	p := &parser{
		lex: &lexer{
			src: data,
		},
	}
	m := make(map[string]interface{})

	if t := p.next(); t.Typ == IDENTIFIER && p.peek().Typ == LBRACE {
		m[t.Val] = p.parseInside()
	}

	return m
}

// -------------------------------------------------------

// parseInside parses inside braces
func (p *parser) parseInside() map[string]interface{} {
	m := make(map[string]interface{})
	p.next() // skip left brace

loop:
	for {
		switch t := p.next(); t.Typ {
		case EOF:
			break loop
		case ILLEGAL:
			log.Fatalf("%+v", t)
		case RBRACE:
			return m
		case IDENTIFIER:
			if p.peek().Typ == LBRACE {
				m[t.Val] = p.parseInside()
			} else if p.peek().Typ == IDENTIFIER {
				m[t.Val] = p.next().Val
			}
		default:
			log.Fatalf("unexpected token: %+v", t)
		}
	}

	return m
}

// -------------------------------------------------------

// next returns the next token.
func (p *parser) next() token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.token()
	}

	return p.token[p.peekCount]
}

// peek returns but does not consume the next token.
func (p *parser) peek() token {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.token()

	return p.token[0]
}
