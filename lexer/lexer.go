package lexer

import (
	"unicode/utf8"

	"github.com/bucketd/go-graphqlparser/token"
)

// eof represents the end of input.
const eof = rune(0)

// Token represents a small, easily categorisable data structure that is fed to the parser to
// produce the abstract syntax tree (AST).
type Token struct {
	Type     token.Type // The token type.
	Literal  string     // The literal value consumed.
	Position int        // The starting position, in runes, of this token in the input.
	Line     int        // The line number at the start of this item.
}

// Lexer holds the state of a state machine for lexically analysing GraphQL queries.
type Lexer struct {
	input    string // The input query string.
	inputLen int    // The input length, in bytes.
	pos      int    // The start position of the last rune read, in bytes.
	rpos     int    // The start position of the last rune read, in runes.
	token    Token  // The last token that was lexed.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input string) *Lexer {
	return &Lexer{
		input:    input,
		inputLen: len(input),
	}
}

// Scan attempts
func (l *Lexer) Scan() Token {
	r := l.read()

	switch {
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
		return l.scanName()
	case r == '{' || r == '}':
		return l.scanPunctuator(r)
	// Ignore spaces...
	case r == ' ':
		l.read()
		return l.Scan()
	}

	return Token{
		Type:     token.EOF,
		Position: l.pos,
		// TODO(seeruk): Line number.
	}
}

func (l *Lexer) scanName() Token {
	start := l.pos
	end := l.pos + 1

	// We already know the first rune is valid part of a name.
	//l.read()

	var done bool
	for !done {
		r := l.read()

		switch {
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
			end++
		default:
			done = true
		}
	}

	return Token{
		Type:     token.Name,
		Literal:  string(l.input[start:end]),
		Position: start,
		// TODO(seeruk): Line number.
	}
}

func (l *Lexer) scanPunctuator(r rune) Token {
	start := l.pos

	//r := l.read()

	return Token{
		Type:     token.Punctuator,
		Literal:  string(r),
		Position: start,
	}
}

// Next will return true if there are more tokens yet to be scanned in this lexer's input.
func (l *Lexer) Next() bool {
	l.token = l.Scan()
	if l.token.Type == token.EOF {
		return false
	}
	return true
}

// Token returns the last token that was lexed by this lexer.
func (l *Lexer) Token() Token {
	return l.token
}

func (l *Lexer) Read() rune {
	return l.read()
}

func (l *Lexer) Unread() {
	l.unread()
}

// read attempts to read the next rune from the input. Returns the EOF rune if an error occurs. The
// return values are the rune that was read, and it's width in bytes.
func (l *Lexer) read() rune {
	if l.pos+1 > l.inputLen {
		return eof
	}

	// This looks pretty weird, but it's the fastest way I've found to get a rune out of a string,
	// one at a time, reliably.
	var r rune
	for _, r = range l.input[l.pos:] {
		break
	}

	w := utf8.RuneLen(r)

	l.pos += w
	l.rpos++

	return r
}

// unread attempts to rewind the underlying buffered reader, allowing a previously read rune to be
// read again.
func (l *Lexer) unread() {
	// If we've not read anything, we can't unread.
	if l.pos <= 0 {
		return
	}

	l.pos--
}
