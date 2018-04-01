package lexer

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/bucketd/go-graphqlparser/token"
)

const (
	// er represents an "empty" rune, but is also an invalid one.
	er = rune(-1)
	// eof represents the end of input.
	eof = rune(0)
	// maxBytes is the max bytes to read at a time.
	maxBytes = 4

	cr  = rune(0x000D) // Literal '\r'.
	lf  = rune(0x000A) // Literal '\n'.
	tab = rune(0x0009) // Literal '	'.
	ws  = rune(0x0020) // Literal ' '.
	com = rune(0x002C) // Literal ','.
	bom = rune(0xFEFF) // Unicode BOM.
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

	// Positional information.
	pos  int // The start position of the last rune read, in runes, on the current line.
	line int // The current line number.

	// Previously read information.
	lbs  []byte // Last bytes read.
	lbsl int    // Length of last bytes read.
	ur   rune   // Unread rune, will be read as next rune if not equal to `ef`.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input io.Reader) *Lexer {
	return &Lexer{
		input: input,
		line:  1,
		ur:    er,
	}
}

// Scan attempts to read the next significant token from the input. Tokens that are not understood
// will yield an "illegal" token.
func (l *Lexer) Scan() (Token, error) {
	r := l.readNextSignificant()

	switch {
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
		return l.scanName(r)
	case runeIn(r, '!', '$', '(', ')', '.', ':', '=', '@', '[', ']', '{', '|', '}'):
		return l.scanPunctuator(r)
	case (r >= '0' && r <= '9') || r == '-':
		return l.scanNumber(r)
	case r == '#':
		return l.scanComment(r)
	case r == eof:
		return Token{
			Type:     token.EOF,
			Position: l.pos,
			Line:     l.line,
		}, nil
	}

	// TODO(seeruk): Should this just be an error really?
	return Token{
		Type:     token.Illegal,
		Position: l.pos,
		Line:     l.line,
	}, nil
}

// scanComment ...
func (l *Lexer) scanComment(r rune) (Token, error) {
	var was000D bool

	for {
		r = l.read()

		// If on the last iteration we saw a CR, then we should check if we just read an LF on this
		// iteration. If we did, reset line position as the next character is still the start of the
		// next line, then scan.
		if was000D && r == lf {
			l.pos = 0

			return l.Scan()
		}

		// Otherwise, if we saw a CR, and this rune isn't an LF, then we have started reading the
		// next line's runes, so unread the rune we read, and scan the next token.
		if was000D && r != lf {
			l.unread(r)

			return l.Scan()
		}

		// If we encounter a CR at any point, this will be true.
		was000D = r == cr
		if was000D {
			// Carriage return, i.e. '\r'.
			l.line++
			l.pos = 0
			continue
		}

		// If we encounter a LF without a proceeding CR, this will be true.
		if r == lf {
			// Line feed, i.e. '\n'.
			l.line++
			l.pos = 0

			return l.Scan()
		}
	}

}

// scanName ...
func (l *Lexer) scanName(r rune) (Token, error) {
	start := l.pos

	rs := []rune{r}

	var done bool
	for !done {
		r := l.read()

		switch {
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
			rs = append(rs, r)
		default:
			l.unread(r)
			done = true
		}
	}

	return Token{
		Type:     token.Name,
		Literal:  string(rs),
		Position: start,
		Line:     l.line,
	}, nil
}

// scanPunctuator ...
func (l *Lexer) scanPunctuator(r rune) (Token, error) {
	start := l.pos

	if r == '.' {
		rs := []rune{r, l.read(), l.read()}
		if rs[1] != '.' || rs[2] != '.' {
			return Token{}, fmt.Errorf("invalid punctuator, expected \"...\" but got: %q", string(rs))
		}

		return Token{
			Type:     token.Punctuator,
			Literal:  "...",
			Position: start,
			Line:     l.line,
		}, nil
	}

	return Token{
		Type:     token.Punctuator,
		Literal:  string(r),
		Position: start,
		Line:     l.line,
	}, nil
}

// scanNumber ...
func (l *Lexer) scanNumber(r rune) (t Token, err error) {
	start := l.pos

	var rs []rune
	if r == '-' {
		rs = append(rs, r)
		r = l.read()
	}

	readDigits := func(r rune, rs []rune) (rune, []rune, error) {
		if !(r >= '0' && r <= '9') {
			return 0, nil, fmt.Errorf("invalid number, expected digit but got: %q", r)
		}

		rs = append(rs, r)

		var done bool
		for !done {
			r = l.read()

			switch {
			case r >= '0' && r <= '9':
				rs = append(rs, r)
			default:
				done = true
			}
		}
		return r, rs, nil
	}

	if r == '0' {
		rs = append(rs, r)
		r = l.read()

		if r >= '0' && r <= '9' {
			return t, fmt.Errorf("invalid number, unexpected digit after 0: %q", r)
		}
	} else {
		r, rs, err = readDigits(r, rs)
		if err != nil {
			return t, err
		}
	}

	var float bool
	if r == '.' {
		float = true

		rs = append(rs, r)
		r = l.read()

		r, rs, err = readDigits(r, rs)
		if err != nil {
			return t, err
		}
	}

	if runeIn(r, 'e', 'E') {
		float = true

		rs = append(rs, r)
		r = l.read()

		if runeIn(r, '+', '-') {
			rs = append(rs, r)
			r = l.read()
		}

		r, rs, err = readDigits(r, rs)
		if err != nil {
			return t, err
		}
	}

	l.unread(r)

	t.Literal = string(rs)
	t.Line = l.line
	t.Position = start

	t.Type = token.IntValue
	if float {
		t.Type = token.FloatValue
	}

	return t, nil
}

// readNextSignificant reads runes until a "significant" rune is read, i.e. a rune that could be a
// significant token (not whitespace, not tabs, not newlines, not commas, not encoding-specific
// characters, etc.). It also does part of the work for identifying when new lines are encountered
// to increment the line counter.
func (l *Lexer) readNextSignificant() rune {
	var done bool
	var was000D bool

	r := er

	for !done && r != eof {
		r = l.read()

		was000D = r == cr

		switch {
		case was000D:
			// Carriage return, i.e. '\r'.
			l.line++
			l.pos = 0
		case r == lf:
			// Line feed, i.e. '\n'.
			if !was000D {
				// \r\n is not 2 newlines, so we must check what the last rune was.
				l.line++
				l.pos = 0
			}
		case runeIn(r, tab, ws, com, bom):
			// Skip!
		default:
			// Done, this run was significant.
			done = true
		}
	}

	return r
}

// read attempts to read the next rune from the input. Returns the EOF rune if an error occurs. The
// return values are the rune that was read, and it's width in bytes.
func (l *Lexer) read() rune {
	// If we unread a rune, return the one that was unread.
	if l.ur != er {
		ur := l.ur
		l.ur = er
		l.pos++
		return ur
	}

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

	// Update position on current line.
	l.pos++

	// Update last bytes read.
	l.lbs = fbs[w:]
	l.lbsl = maxBytes - w

	return r
}

// unread doesn't really unread anything, it just stores a given rune to be read as the next rune.
// Actually doing an unread would be trickier given the use of a reader...
func (l *Lexer) unread(r rune) {
	l.ur = r
	l.pos--
}

// runeIn returns true if the rune `r` matches a code point in `rs`.
func runeIn(r rune, rs ...rune) bool {
	for _, cp := range rs {
		if cp == r {
			return true
		}
	}

	return false
}
