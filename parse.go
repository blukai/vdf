package vdf

import (
	"log"
)

// I'm too lazy to comment thi..

type Parser struct {
	lex       *Lexer
	token     [3]Token
	peekCount int
}

// -------------------------------------------------------

func Parse(data string) map[string]interface{} {
	p := &Parser{lex: &Lexer{src: data}}
	m := make(map[string]interface{})

	if t := p.next(); t.Typ == IDENTIFIER && p.peek().Typ == LBRACE {
		m[t.Val] = p.parseInside()
	}

	return m
}

// -------------------------------------------------------

func (p *Parser) parseInside() map[string]interface{} {
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
func (p *Parser) next() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.Token()
	}

	return p.token[p.peekCount]
}

// peek returns but does not consume the next token.
func (p *Parser) peek() Token {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.Token()

	return p.token[0]
}
