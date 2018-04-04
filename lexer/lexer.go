package lexer

import (
	"fmt"
	"unicode/utf8"
	"unsafe"

	"github.com/bucketd/go-graphqlparser/token"
)

const (
	// er represents an "empty" rune, but is also an invalid one.
	er = rune(-1)
	// eof represents the end of input.
	eof = rune(0)

	bom = rune(0xFEFF) // Unicode BOM.
	cr  = rune(0x000D) // Literal '\r'.
	lf  = rune(0x000A) // Literal '\n'.
	tab = rune(0x0009) // Literal '	'.
	ws  = rune(0x0020) // Literal ' '.
	com = rune(0x002C) // Literal ','.
	bsl = rune(0x005C) // Literal reverse solidus (backslash).
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
	input    []byte // Raw input is just a byte slice. It is expected to be UTF-8 encoded characters.
	inputLen int    // Length of the input, in bytes.

	// Positional information.
	pos  int // The start position of the last rune read, in bytes.
	lpos int // The start position of the last rune read, in runes, on the current line.
	line int // The current line number.

	// Previous read information.
	lrw int // The width of the last rune read.
}

// New returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func New(input []byte) *Lexer {
	return &Lexer{
		input:    input,
		inputLen: len(input),
		line:     1,
	}
}

// Scan attempts to read the next significant token from the input. Tokens that are not understood
// will yield an "illegal" token.
func (l *Lexer) Scan() (Token, error) {
	r := l.readNextSignificant()

	switch {
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
		return l.scanName(r)

	case r == '!' || r == '$' || r == '(' || r == ')' || r == '.' || r == ':' || r == '=' || r == '@' || r == '[' || r == ']' || r == '{' || r == '|' || r == '}':
		return l.scanPunctuator(r)

	case (r >= '0' && r <= '9') || r == '-':
		return l.scanNumber(r)

	case r == '#':
		return l.scanComment(r)

	case r == '"':
		rs := []rune{r, l.read(), l.read()}
		if rs[1] == '"' && rs[2] == '"' {
			return l.scanBlockString(r)
		}
		l.unread()
		l.unread()
		return l.scanString(r)

	case r == eof:
		return Token{
			Type:     token.EOF,
			Position: l.lpos,
			Line:     l.line,
		}, nil
	}

	// TODO(seeruk): Should this just be an error really?
	// https://github.com/graphql/graphql-js/blob/master/src/language/lexer.js#L347
	return Token{
		Type:     token.Illegal,
		Position: l.lpos,
		Line:     l.line,
	}, nil
}

// scanString ...
// TODO(Luke-Vear): finish logic and helper functions.
// TODO(Luke-Vear): are appends needed?
func (l *Lexer) scanString(r rune) (Token, error) {
	//byteStart := l.pos - 1
	runeStart := l.lpos
	rs := []rune{r}

	var done bool
	for !done {
		r = l.read()

		switch {
		case r == '"':
			rs = append(rs, r)
			return Token{
				Type:     token.StringValue,
				Literal:  string(rs),
				Position: runeStart,
				Line:     l.line,
			}, nil

		case r < ws && r != tab:
			return Token{}, fmt.Errorf("invalid character within string: %q", r)

		case r == bsl:
			r, err := escapedChar(l)
			if err != nil {
				return Token{}, err
			}
			rs = append(rs, r)

		default:
			rs = append(rs, r)
		}
	}

	return Token{
		Type:     token.StringValue,
		Literal:  string(rs),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

// scanBlockString ...
func (l *Lexer) scanBlockString(r rune) (Token, error) {
	return Token{}, nil
}

func escapedChar(l *Lexer) (rune, error) {
	r := l.read()
	switch r {
	case '"':
		return '"', nil
	case '/':
		return '/', nil
	case bsl:
		return bsl, nil
	case 'b':
		return '\b', nil
	case 'f':
		return '\f', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil

	case 'u':
		rs := []rune{l.read(), l.read(), l.read(), l.read()}
		r := ucptor(rs)
		if r < 0 {
			return 0, fmt.Errorf("invalid character escape sequence: %s,", "\\u"+string(r))
		}
		return r, nil
	}

	return 0, fmt.Errorf("invalid character escape sequence: %s", "\\"+string(r))
}

// TODO(seeruk): Here: https://github.com/graphql/graphql-js/blob/master/src/language/lexer.js#L689
// TODO(Luke-Vear): Discuss rename, I think ucptor isn't easy to grok.
func ucptor(rs []rune) rune {
	a, b, c, d := hexRuneToInt(rs[0]), hexRuneToInt(rs[1]), hexRuneToInt(rs[2]), hexRuneToInt(rs[3])
	return rune(a<<12 | b<<8 | c<<4 | d<<0)
}

// hexRuneToInt changes a character into its integer value in hexadecimal. For example:
// the character 'A' is 65 in decimal representation but its value is 10 in hexadecimal.
func hexRuneToInt(r rune) int {
	switch {
	case r >= '0' && r <= '9':
		return int(r - 48)
	case r >= 'A' && r <= 'F':
		return int(r - 55)
	case r >= 'a' && r <= 'f':
		return int(r - 87)
	}
	return -1
}

// scanComment scans valid GraphQL comments.
func (l *Lexer) scanComment(r rune) (Token, error) {
	var was000D bool

	for {
		r = l.read()
		if r == eof {
			return l.Scan()
		}

		// If on the last iteration we saw a CR, then we should check if we just read an LF on this
		// iteration. If we did, reset line position as the next character is still the start of the
		// next line, then scan.
		if was000D && r == lf {
			l.lpos = 0

			return l.Scan()
		}

		// Otherwise, if we saw a CR, and this rune isn't an LF, then we have started reading the
		// next line's runes, so unread the rune we read, and scan the next token.
		// Q: not hit by tests? can this code be reached?
		if was000D && r != lf {
			l.unread()

			return l.Scan()
		}

		// If we encounter a CR at any point, this will be true.
		was000D = r == cr
		if was000D {
			// Carriage return, i.e. '\r'.
			l.line++
			l.lpos = 0
			continue
		}

		// If we encounter a LF without a proceeding CR, this will be true.
		if r == lf {
			// Line feed, i.e. '\n'.
			l.line++
			l.lpos = 0

			return l.Scan()
		}
	}
}

// scanName scans valid GraphQL name tokens.
func (l *Lexer) scanName(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	var done bool
	for !done {
		r := l.read()

		switch {
		case r == eof:
			done = true
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
			continue
		default:
			l.unread()
			done = true
		}
	}

	return Token{
		Type:     token.Name,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

// scanPunctuator scans valid GraphQL punctuation tokens.
func (l *Lexer) scanPunctuator(r rune) (Token, error) {
	byteStart := l.pos
	runeStart := l.lpos

	if r == '.' {
		rs := []rune{r, l.read(), l.read()}
		if rs[1] != '.' || rs[2] != '.' {
			return Token{}, fmt.Errorf("invalid punctuator, expected \"...\" but got: %q", string(rs))
		}

		return Token{
			Type:     token.Punctuator,
			Literal:  "...",
			Position: runeStart,
			Line:     l.line,
		}, nil
	}

	// TODO(seeruk): Using other token types for each type of punctuation may actually be faster.
	return Token{
		Type:     token.Punctuator,
		Literal:  btos(l.input[byteStart-l.lrw : byteStart]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

// scanNumber scans valid GraphQL integer and float value tokens.
func (l *Lexer) scanNumber(r rune) (Token, error) {
	byteStart := l.pos - 1

	var float bool // If true, number is float.
	var err error  // So no shadowing of r.

	// Check for preceeding minus sign
	if r == '-' {
		r = l.read()
	}

	// Check if digits begins with zero
	if r == '0' {
		r = l.read()

		// If there is another digit after zero, error.
		if r >= '0' && r <= '9' {
			// TODO(seeruk): Unread here?
			// LV: I think on error we hard stop processing.
			return Token{}, fmt.Errorf("invalid number, unexpected digit after 0: %q", r)
		}

		// If number does not begin with zero, read the digits.
		// If the first character is not a digit, error.
	} else {
		r, err = l.readDigits(r)
		if err != nil {
			return Token{}, err
		}
	}

	// Check for a decimal place, if there is a decimal place this number is a float.
	if r == '.' {
		float = true

		r = l.read()

		// Read the digits after the decimal place if the first character is not a digit, error.
		r, err = l.readDigits(r)
		if err != nil {
			return Token{}, err
		}
	}

	// Check for exponent sign, if there is an exponent sign this number is a float.
	if r == 'e' || r == 'E' {
		float = true

		r = l.read()

		// Check for positive or negative symbol infront of the value.
		if r == '+' || r == '-' {
			r = l.read()
		}

		// Read the exponent digitas, if the first character is not a digit, error.
		r, err = l.readDigits(r)
		if err != nil {
			return Token{}, err
		}
	}

	if r != eof {
		l.unread()
	}

	t := Token{
		Literal:  btos(l.input[byteStart:l.pos]),
		Line:     l.line,
		Position: byteStart,
	}

	t.Type = token.IntValue
	if float {
		t.Type = token.FloatValue
	}

	return t, nil
}

// readDigits reads up until the next non-numeric character in the input.
func (l *Lexer) readDigits(r rune) (rune, error) {
	if !(r >= '0' && r <= '9') {
		return eof, fmt.Errorf("invalid number, expected digit but got: %q", r)
	}

	var done bool
	for !done {
		r = l.read()

		switch {
		case r >= '0' && r <= '9':
			continue
		default:
			// No need to unread here. We actually want to read the character after the numbers.
			done = true
		}
	}

	return r, nil
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
			l.lpos = 0
		case r == lf:
			// Q: not hit by tests? can this code be reached?
			// Line feed, i.e. '\n'.
			if !was000D {
				// \r\n is not 2 newlines, so we must check what the last rune was.
				l.line++
				l.lpos = 0
			}
		case r == tab || r == ws || r == com || r == bom:
			// Skip!
		default:
			// Done, this run was significant.
			done = true
		}
	}

	return r
}

// read moves forward in the input, and returns the next rune available. This function also updates
// the position(s) that the lexer keeps track of in the input so the next read continues from where
// the last left off. Returns the EOF rune if we hit the end of the input.
func (l *Lexer) read() rune {
	if l.pos >= l.inputLen {
		return eof
	}

	var r rune
	var w int
	if sbr := l.input[l.pos]; sbr < utf8.RuneSelf {
		r = rune(sbr)
		w = 1
	} else {
		r, w = utf8.DecodeRune(l.input[l.pos:])
	}

	l.pos += w
	l.lpos++

	l.lrw = w

	return r
}

// unread goes back one rune's worth of bytes in the input, changing the
// positions we keep track of.
// Does not currently go back a line.
func (l *Lexer) unread() {
	l.pos -= l.lrw

	if l.pos > 0 {
		// update rune width for further rewind
		_, l.lrw = utf8.DecodeLastRune(l.input[:l.pos])
	} else {
		// If we're already at the start, set this to so we don't end up with a negative position.
		l.lrw = 0
	}

	if l.lpos > 0 {
		l.lpos--
	}
}

// btos takes the given bytes, and turns them into a string.
// Q: naming btos or bbtos? :D
// TODO(seeruk): Is this actually portable then?
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
