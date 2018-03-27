package lexer

import (
	"io"
	"unicode/utf8"

	"github.com/bucketd/go-graphqlparser/token"
)

const (
	// maxBytes is the max bytes to read at a time.
	maxBytes = 4
	// eof represents the end of input.
	eof = rune(0)
)

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
	input io.Reader // The input query string.
	pos   int       // The start position of the last rune read, in bytes.
	rpos  int       // The start position of the last rune read, in runes.
	wid   int       // The width of the last rune read, in bytes.
	token Token     // The last token that was lexed.

	lbs  []byte // Last bytes read.
	lbsl int    // Length of last bytes read.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input io.Reader) *Lexer {
	return &Lexer{
		input: input,
	}
}

// Scan attempts
func (l *Lexer) Scan() Token {
	r := l.read()

	switch {
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
		return l.scanName(r)
	case r == '{' || r == '}':
		return l.scanPunctuator(r)
	// Ignore spaces...
	case r == ' ':
		return l.Scan()
	}

	return Token{
		Type:     token.EOF,
		Position: l.pos,
		// TODO(seeruk): Line number.
	}
}

func (l *Lexer) scanName(r rune) Token {
	start := l.rpos - 1

	rs := []rune{r}

	var done bool
	for !done {
		r := l.read()

		switch {
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
			rs = append(rs, r)
		default:
			done = true
		}
	}

	return Token{
		Type:     token.Name,
		Literal:  string(rs),
		Position: start,
		// TODO(seeruk): Line number.
	}
}

func (l *Lexer) scanPunctuator(r rune) Token {
	start := l.rpos - 1

	//r := l.read()

	return Token{
		Type:     token.Punctuator,
		Literal:  string(r),
		Position: start,
	}
}

// read attempts to read the next rune from the input. Returns the EOF rune if an error occurs. The
// return values are the rune that was read, and it's width in bytes.
func (l *Lexer) read() rune {
	// Start off with the maximum amount of bytes we should be reading. If we have some previous
	// bytes leftover from another read, then we should cut down the length of the next set of bytes
	// that are going to be read (`bl`).
	bl := maxBytes - l.lbsl

	// The length of leftover bytes (l.lbsl) and the length of this byte slice together should be as
	// long as `maxBytes`.
	bs := make([]byte, bl)

	// Read as much as possible from the reader. Or hit the end of the input. If we shouldn't read
	// anything because we've already got `maxBytes` leftover from a previous read, then we don't
	// bother attempting a read.
	if l.lbsl < maxBytes {
		_, err := l.input.Read(bs)
		if err != nil && l.lbsl == 0 {
			return eof
		}
	}

	// Combine the leftover bytes from the last read with the bytes we've just read.
	fbs := append(l.lbs, bs...)

	// Hack to get runes from strings faster.
	var r rune
	for _, r = range string(fbs) {
		break
	}

	// Find start position of next character.
	w := utf8.RuneLen(r)

	l.pos += w
	l.rpos++

	l.lbs = fbs[w:]
	l.lbsl = 4 - w

	return r
}
