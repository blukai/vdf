package vdf

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type (
	// stateFn represents the state of the lexer as a function that
	// returns the next state.
	stateFn func(*lexer) stateFn

	// lexer holds the lexing state.
	lexer struct {
		src    string     // the data being scanned
		rPos   int        // current rune position in the input
		tPos   int        // start position of current token
		rWidth int        // width of last rune read from input
		tokens chan token // channel of scanned tokens
		state  stateFn    // the next lexing function to enter
	}
)

const eof = -1

// -------------------------------------------------------

// token returns the next token emitted by the lexer.
// If called for the first time, makes token channel.
func (l *lexer) token() token {
	if l.tokens == nil {
		l.tokens = make(chan token)
		go l.run()
	}

	return <-l.tokens
}

// -------------------------------------------------------

func lex(l *lexer) stateFn {
	switch r := l.next(); r {
	case eof: // end of file
		l.emit(EOF)
		return nil
	case '{':
		l.emit(LBRACE)
	case '}':
		l.emit(RBRACE)
	case '"':
		l.ignore() // ignore opening quote mark
		return lexIdentifier
	case '/':
		return lexComment
	default:
		switch {
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return l.errorf("unrecognized character: %#U", r)
		}
	}

	return lex
}

// lexIdentifier lex an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
loop:
	for {
		switch r := l.next(); {
		case isAplphameric(r):
			// absorb
		case r == '"': // ignore closing quote mark
			l.backup()
			l.emit(IDENTIFIER)
			l.next()
			break loop
		}
	}

	return lex
}

// lexComment lex(ignore) comments
func lexComment(l *lexer) stateFn {
	for r := l.next(); r != '\r' && r != '\n' && r != eof; {
		r = l.next()
	}
	l.backup()
	l.ignore()

	return lex
}

// -------------------------------------------------------

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.rPos -= l.rWidth
}

// emit receives a token type and pushes a new token into the tokens channel.
func (l *lexer) emit(typ tokenType) {
	l.tokens <- token{typ, l.src[l.tPos:l.rPos]}
	l.tPos = l.rPos
}

// errorf returns an error token and terminates the lexing by passing
// back a nil pointer that will be the next state, terminating l.token.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{ILLEGAL, fmt.Sprintf(format, args...)}

	return nil
}

// ignore skips over the pending input before this point
func (l *lexer) ignore() {
	l.tPos = l.rPos
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.rPos >= len(l.src) {
		l.rWidth = 0

		return eof
	}
	// rune, width
	r, w := utf8.DecodeRuneInString(l.src[l.rPos:])
	l.rWidth = w
	l.rPos += l.rWidth

	return r
}

// run lexes the input by executing state functions until the state is nil.
func (l *lexer) run() {
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
