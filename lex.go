package vdf

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type (
	// stateFn represents the state of the lexer as a function that
	// returns the next state.
	stateFn func(*Lexer) stateFn

	// Lexer holds the lexing state.
	Lexer struct {
		Data   string     // the data being scanned
		rPos   int        // current rune position in the input
		tPos   int        // start position of current token
		rWidth int        // width of last rune read from input
		tokens chan Token // channel of scanned tokens
		state  stateFn    // the next lexing function to enter
	}
)

const eof = -1

// -------------------------------------------------------

// Token returns the next token emitted by the Lexer.
// If called for the first time, makes token channel.
func (l *Lexer) Token() Token {
	if l.tokens == nil {
		l.tokens = make(chan Token)
		go l.run()
	}

	return <-l.tokens
}

// -------------------------------------------------------

func lex(l *Lexer) stateFn {
	switch r := l.next(); r {
	case eof: // end of file
		l.emit(EOF)
		return nil
	case '{':
		l.emit(LBRACE)
	case '}':
		l.emit(RBRACE)
	case '"': // lex identifier
	loop:
		for {
			switch r := l.next(); {
			case isAplphameric(r):
				// absorb
			default:
				l.emit(IDENTIFIER)
				break loop
			}
		}
	case '/': // skip comments
		for r := l.next(); r != '\r' && r != '\n' && r != eof; {
			r = l.next()
		}
		l.backup()
		l.ignore()
	default:
		switch {
		case unicode.IsSpace(r): // skip spaces
			l.ignore()
		default:
			return l.errorf("unrecognized character: %#U", r)
		}
	}

	return lex
}

// -------------------------------------------------------

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.rPos -= l.rWidth
}

// emit receives a token type and pushes a new token into the tokens channel.
func (l *Lexer) emit(typ TokenType) {
	l.tokens <- Token{typ, l.Data[l.tPos:l.rPos]}
	l.tPos = l.rPos
}

// errorf returns an error token and terminates the lexing by passing
// back a nil pointer that will be the next state, terminating l.Token.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{ILLEGAL, fmt.Sprintf(format, args...)}

	return nil
}

// ignore skips over the pending input before this point
func (l *Lexer) ignore() {
	l.tPos = l.rPos
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.rPos >= len(l.Data) {
		l.rWidth = 0

		return eof
	}
	// rune, width
	r, w := utf8.DecodeRuneInString(l.Data[l.rPos:])
	l.rWidth = w
	l.rPos += l.rWidth

	return r
}

// run lexes the input by executing state functions until the state is nil.
func (l *Lexer) run() {
	for l.state = lex; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

// -------------------------------------------------------

// isAlphameric reports whether r is an alphabetic, digit, underscore, dot or dash.
func isAplphameric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-'
}
