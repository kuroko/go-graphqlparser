package lexer

import (
	"bufio"
	"fmt"
	"io"
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

// String returns a string representation of a token.
func (t Token) String() string {
	res, ok := token.TypeNames[t.Type]
	if !ok {
		panic(fmt.Sprintf("invalid type given: %d", t.Type))
	}

	if t.Literal != "" {
		res += "(" + t.Literal + ")"
	}

	return res
}

// Lexer holds the state of a state machine for lexically analysing GraphQL queries.
type Lexer struct {
	reader        *bufio.Reader // A buffered reader, for reading runes safely.
	bytePos       int           // The start position of the last rune read, in bytes.
	runePos       int           // The start position of the last rune read, in runes.
	prevRune      rune          // The last rune that was read.
	prevRuneWidth int           // The width of the last rune read, in bytes.
	token         Token         // The last token that was lexed.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(input),
	}
}

// Scan attempts
func (l *Lexer) Scan() Token {
	return Token{}
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

func (l *Lexer) Read() (rune, int) {
    return l.read()
}

// read attempts to read the next rune from the buffered reader. Returns the EOF rune if an error
// occurs. The return values are the rune that was read, and it's width in bytes.
func (l *Lexer) read() (rune, int) {
	r, w, err := l.reader.ReadRune()
	if err != nil {
		// TODO(seeruk): Should all errors yield EOF?
		return eof, utf8.RuneLen(eof)
	}

	// Update previous rune.
	l.prevRune = r
	l.prevRuneWidth = w

	// Update positions.
	l.bytePos += w
	l.runePos++

	return r, w
}

// unread attempts to rewind the underlying buffered reader, allowing a previously read rune to be
// read again.
func (l *Lexer) unread() {
	// If we've just unread, we can't do it again.
	if l.prevRune < 0 || l.prevRuneWidth < 0 {
		return
	}

	// If we've not read anything, we can't unread.
	if l.runePos <= 0 || l.bytePos <= 0 {
		return
	}

	err := l.reader.UnreadRune()
	if err != nil {
		// TODO(seeruk): Should maybe handle this better?
		return
	}

	// Disallow another unread, we don't have the information to allow this.
	l.prevRune = rune(-1)
	l.prevRuneWidth = -1

	// Update position.
	l.bytePos -= l.prevRuneWidth
	l.runePos--
}
