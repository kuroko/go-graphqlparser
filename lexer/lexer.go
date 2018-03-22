package lexer

import (
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

// Lexer holds the state of a state machine for lexically analysing GraphQL queries.
type Lexer struct {
	input io.ReadSeeker // The input query string.
	pos   int           // The start position of the last rune read, in bytes.
	rpos  int           // The start position of the last rune read, in runes.
	wid   int           // The width of the last rune read, in bytes.
	token Token         // The last token that was lexed.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input io.ReadSeeker) *Lexer {
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
	// 4 bytes is the maximum size of a valid UTF-8 encoded character point.
	bs := make([]byte, 4)

	// Read as much as possible from the reader. Or hit the end of the input.
	_, err := l.input.Read(bs)
	if err != nil {
		return eof
	}

	// Hack to get runes from strings faster.
	var r rune
	for _, r = range string(bs) {
		break
	}

	// Find start position of next character.
	w := utf8.RuneLen(r)

	l.pos += w
	l.rpos++

	// Seek to the start position of the next character, otherwise we'll get garbage results.
	l.input.Seek(int64(l.pos), io.SeekStart)

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
