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
	ws  = rune(0x0020) // Literal ' '.
	com = rune(0x002C) // Literal ','.
	dq  = rune(0x0022) // '\"' double quote.
	fsl = rune(0x002F) // '\/' solidus (forward slash).
	bck = rune(0x0008) // '\b' backspace.
	ff  = rune(0x000C) // '\f' form feed.
	lf  = rune(0x000A) // '\n' line feed (new line).
	cr  = rune(0x000D) // '\r' carriage return.
	tab = rune(0x0009) // '\t' horizontal tab.
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
		r1, _ := l.read()
		r2, _ := l.read()
		rs := []rune{r1, r2}
		if rs[0] == '"' && rs[1] == '"' {
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
func (l *Lexer) scanString(r rune) (Token, error) {
	// var w int

	// startPos := l.pos
	// startLPos := l.lpos
	startLine := l.line
	// startLRW := l.lrw

	runeStart := l.lpos

	// var bc int

	// var done bool
	// for !done {
	// 	r, w = l.read()
	// 	bc += w

	// 	switch {
	// 	case r == '"':
	// 		done = true

	// 	case r < ws && r != tab:
	// 		return Token{}, fmt.Errorf("invalid character within string: %q", r)

	// 	case r == bsl:
	// 		r, _ = l.read()

	// 		// No need to increment bc here, if we hit backslash, we should already have incremented
	// 		// the counter by 1. That one byte increment should satisfy the width of any escape
	// 		// sequence other than unicode escape sequences when decoded as a rune. We handle the
	// 		// unicode escape sequence case further down.
	// 		//bc += w

	// 		if r == 'u' {
	// 			_, _ = l.read()
	// 			_, _ = l.read()
	// 			_, _ = l.read()
	// 			_, _ = l.read()

	// 			// Increment bc by 3, because we've already incremented by 1 above at the start of
	// 			// this loop iteration. We increment by 3 here because we want to have incremented
	// 			// by 4 in total. 4 bytes being the maximum width of a valid unicode escape sequence
	// 			// supported by GraphQL.
	// 			bc += 3
	// 		}
	// 	}
	// }

	// l.pos = startPos
	// l.lpos = startLPos
	// l.line = startLine
	// l.lrw = startLRW

	// Sadly, allocations cannot be avoided here unless we modify the input byte slice to make
	// string scanning work. This is because we have to replace the escape sequences with their
	// actual rune counterparts and use that as the token's literal value. To store that data, we
	// need bytes to be allocated. In this case, in the form of a string - the compiler optimises it
	// for us to help keep our code simple.
	var str string
	for {
		r, _ = l.read()

		switch {
		case r == '"':
			return Token{
				Type:     token.StringValue,
				Literal:  str,
				Position: runeStart,
				Line:     startLine,
			}, nil

		case r < ws && r != tab:
			return Token{}, fmt.Errorf("invalid character within string: %q", r)

		case r == bsl:
			r, err := escapedChar(l)
			if err != nil {
				return Token{}, err
			}

			str += string(r)

		default:
			str += string(r)
		}
	}
}

// scanBlockString ...
func (l *Lexer) scanBlockString(r rune) (Token, error) {
	startLine := l.line
	runeStart := l.lpos
	var str string
	for {
		r, _ = l.read()

		switch {
		case r == '"':
			r1, _ := l.read()
			r2, _ := l.read()
			if r1 == '"' && r2 == '"' {
				return Token{
					Type:     token.StringValue,
					Literal:  str,
					Position: runeStart,
					Line:     startLine,
				}, nil
			}
			str += string(r)
			l.unread()
			l.unread()

		case r < ws && r != tab:
			return Token{}, fmt.Errorf("invalid character within string: %q", r)

		case r == bsl:
			r1, _ := l.read()
			r2, _ := l.read()
			r3, _ := l.read()
			if r1 == '"' && r2 == '"' && r3 == '"' {
				return Token{
					Type:     token.StringValue,
					Literal:  str,
					Position: runeStart,
					Line:     startLine,
				}, nil
			}
			l.unread()
			l.unread()

			r, err := escapedChar(l)
			if err != nil {
				return Token{}, err
			}

			str += string(r)

		default:
			str += string(r)
		}
	}
}

func escapedChar(l *Lexer) (rune, error) {
	r, _ := l.read()
	switch r {
	case '"':
		return dq, nil
	case '/':
		return fsl, nil
	case '\\': // escaped single backslash '\' == U+005C
		return bsl, nil
	case 'b':
		return bck, nil
	case 'f':
		return ff, nil
	case 'n':
		return lf, nil
	case 'r':
		return cr, nil
	case 't':
		return tab, nil

	case 'u':
		r1, _ := l.read()
		r2, _ := l.read()
		r3, _ := l.read()
		r4, _ := l.read()

		//rs := []rune{r1, r2, r3, r4}
		r := ucptor(r1, r2, r3, r4)
		//r := ucptor(rs[1], rs[2], rs[3], rs[4])
		if r < 0 {
			return 0, fmt.Errorf("invalid character escape sequence: %s", "\\u"+string([]rune{r1, r2, r3, r4}))
		}
		return r, nil
	}

	return 0, fmt.Errorf("invalid character escape sequence: %s", "\\"+string(r))
}

// TODO(seeruk): Here: https://github.com/graphql/graphql-js/blob/master/src/language/lexer.js#L689
// TODO(Luke-Vear): Discuss rename, I think ucptor isn't easy to grok.
func ucptor(ar, br, cr, dr rune) rune {
	ai, bi, ci, di := hexRuneToInt(ar), hexRuneToInt(br), hexRuneToInt(cr), hexRuneToInt(dr)
	return rune(ai<<12 | bi<<8 | ci<<4 | di<<0)
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
		r, _ = l.read()
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
		r, _ := l.read()

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
		r2, _ := l.read()
		r3, _ := l.read()

		rs := []rune{r, r2, r3}
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
		r, _ = l.read()
	}

	// Check if digits begins with zero
	if r == '0' {
		r, _ = l.read()

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

		r, _ = l.read()

		// Read the digits after the decimal place if the first character is not a digit, error.
		r, err = l.readDigits(r)
		if err != nil {
			return Token{}, err
		}
	}

	// Check for exponent sign, if there is an exponent sign this number is a float.
	if r == 'e' || r == 'E' {
		float = true

		r, _ = l.read()

		// Check for positive or negative symbol infront of the value.
		if r == '+' || r == '-' {
			r, _ = l.read()
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
		r, _ = l.read()

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
		r, _ = l.read()

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
func (l *Lexer) read() (rune, int) {
	if l.pos >= l.inputLen {
		return eof, 0
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

	return r, w
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
