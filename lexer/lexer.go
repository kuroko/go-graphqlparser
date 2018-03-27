package lexer

import (
	"fmt"
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

	lbs  []byte
	lbsl int
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input io.Reader) *Lexer {
	return &Lexer{
		input: input,
	}
}

// Scan attempts
func (l *Lexer) Scan() (Token, error) {
	r := l.read()

	switch {
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
		return l.scanName(r)
	case r == '{' || r == '}':
		return l.scanPunctuator(r)
	case (r >= '0' && r <= '9') || r == '-':
		return l.scanNumber(r)
	// Ignore spaces...
	case r == ' ':
		return l.Scan()
	}

	return Token{
		Type:     token.EOF,
		Position: l.pos,
		// TODO(seeruk): Line number.
	}, nil
}

func (l *Lexer) scanName(r rune) (Token, error) {
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
	}, nil
}

func (l *Lexer) scanPunctuator(r rune) (Token, error) {
	start := l.rpos - 1

	//r := l.read()

	return Token{
		Type:     token.Punctuator,
		Literal:  string(r),
		Position: start,
	}, nil
}

func (l *Lexer) scanNumber(r rune) (Token, error) {
	start := l.rpos - 1

	var rs []rune
	if r == '-' {
		rs = append(rs, r)
		r = l.read()
	}

	readDigits := func(r rune, l *Lexer, rs []rune) ([]rune, error) {
		if !(r >= '0' && r <= '9') {
			return nil, fmt.Errorf("Invalid number, expected digit but got: %v", r)
		}
		rs = append(rs, r)

		var done bool
		for !done {
			r := l.read()

			switch {
			case (r >= '0' && r <= '9'):
				rs = append(rs, r)
			default:
				done = true
			}
		}
		return rs, nil
	}

	if r == '0' {
		rs = append(rs, r)
		r := l.read()
		if r >= '0' && r <= '9' {
			return Token{}, fmt.Errorf("Invalid number, unexpected digit after 0: %v", r)
		}
	}

	rs, err := readDigits(r, l, rs)
	if err != nil {
		return Token{}, err
	}

	var float bool
	r = l.read()
	if r == '.' {
		float = true
		rs, err = readDigits(r, l, rs)
		if err != nil {
			return Token{}, err
		}
	}

	if float {
		return Token{
			Type:     token.FloatValue,
			Literal:  string(r),
			Position: start,
			// TODO(seeruk): Line number.
		}, nil
	}
	return Token{
		Type:     token.IntValue,
		Literal:  string(r),
		Position: start,
		// TODO(seeruk): Line number.
	}, nil
}

// Next will return true if there are more tokens yet to be scanned in this lexer's input.
func (l *Lexer) Next() bool {
	l.token, _ = l.Scan() // err checking?
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
	l.lbsl = maxBytes - w

	return r
}
