package language

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
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

	// nothing to see here
	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	// nothing to see here
	maskx = 0x3F // 0011 1111

	// nothing to see here
	runeError    = '\uFFFD'
	maxRune      = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	// these are not the bytes you're looking for
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000
)

// Lexer holds the state of a state machine for lexically analysing GraphQL queries.
type Lexer struct {
	input    []byte // Raw input is just a byte slice. It is expected to be UTF-8 encoded characters.
	inputLen int    // Length of the input, in bytes.

	// Positional information.
	pos  int // The start position of the last rune read, in bytes.
	lpos int // The start position of the last rune read, in runes, on the current line.
	line int // The current line number.
}

// NewLexer returns a new lexer, for lexically analysing GraphQL queries from a given reader.
func NewLexer(input []byte) *Lexer {
	return &Lexer{
		input:    input,
		inputLen: len(input),
		line:     1,
	}
}

// Scan attempts to read the next significant token from the input. Tokens that are not understood
// will yield an "illegal" token.
func (l *Lexer) Scan() Token {
	r, w := l.readNextSignificant()

	switch {
	case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_':
		return l.scanName(r)

	// TODO: Analyse frequency of occurrences of each punctuator in common queries to figure
	// out the best order for these to be in.
	case r == '{' || r == '}' || r == '[' || r == ']' || r == '!' || r == '$' || r == '(' || r == ')' || r == ':' || r == '@' || r == '&' || r == '.' || r == '=' || r == '|':
		return l.scanPunctuator(r, w)

	case (r >= '0' && r <= '9') || r == '-':
		return l.scanNumber(r)

	case r == '#':
		return l.scanComment(r)

	case r == '"':
		r1, w1 := l.read()
		r2, w2 := l.read()
		if r1 == '"' && r2 == '"' {
			return l.scanBlockString(r)
		}
		l.unread(w2)
		l.unread(w1)

		return l.scanString(r)

	case r == eof:
		return Token{
			Kind:   TokenKindEOF,
			Column: l.lpos + 1,
			Line:   l.line,
		}

	default:
		return Token{
			Kind:    TokenKindIllegal,
			Literal: string(r),
			Column:  l.lpos,
			Line:    l.line,
		}
	}
}

// scanString scans a valid GraphQL string.
func (l *Lexer) scanString(r rune) Token {
	var w int

	startPos := l.pos
	startLPos := l.lpos
	startLine := l.line

	var bc int
	var hasEscape bool

Loop:
	for {
		r, w = l.read()
		bc += w

		switch {
		case r == '"':
			break Loop

		case r < ws && r != tab:
			return Token{
				Kind:    TokenKindIllegal,
				Literal: fmt.Sprintf("invalid character within string: %q", r),
				Column:  startLPos,
				Line:    startLine,
			}

		case r == bsl:
			hasEscape = true

			r, w = l.read()

			// No need to increment bc here, if we hit backslash, we should already have incremented
			// the counter by 1. That one byte increment should satisfy the width of any escape
			// sequence other than unicode escape sequences when decoded as a rune. We handle the
			// unicode escape sequence case further down.
			//bc += w

			if r == 'u' {
				_, _ = l.read()
				_, _ = l.read()
				_, _ = l.read()
				_, _ = l.read()

				// Increment bc by 3, because we've already incremented by 1 above at the start of
				// this loop iteration. We increment by 3 here because we want to have incremented
				// by 4 in total. 4 bytes being the maximum width of a valid unicode escape sequence
				// supported by GraphQL.
				bc += 3
			}
		}
	}

	if !hasEscape {
		return Token{
			Kind:    TokenKindStringValue,
			Literal: btos(l.input[startPos : l.pos-1]),
			Column:  startLPos,
			Line:    startLine,
		}
	}

	l.pos = startPos
	l.lpos = startLPos
	l.line = startLine

	// Sadly, allocations cannot be avoided here unless we modify the input byte slice to make
	// string scanning work. This is because we have to replace the escape sequences with their
	// actual rune counterparts and use that as the token's literal value. To store that data, we
	// need bytes to be allocated.
	bs := make([]byte, 0, bc)
	for {
		r, _ = l.read()

		switch {
		case r == '"' || r == eof:
			return Token{
				Kind:    TokenKindStringValue,
				Literal: btos(bs),
				Column:  startLPos,
				Line:    startLine,
			}

		case r == bsl:
			r, err := escapedChar(l)
			if err != nil {
				return Token{
					Kind:    TokenKindIllegal,
					Literal: err.Error(),
					Column:  startLPos,
					Line:    startLine,
				}
			}

			encodeRune(r, func(b byte) {
				bs = append(bs, b)
			})

		default:
			encodeRune(r, func(b byte) {
				bs = append(bs, b)
			})
		}
	}
}

// scanBlockString scans a valid GraphQL block string.
func (l *Lexer) scanBlockString(r rune) Token {
	startPos := l.pos
	startLPos := l.lpos
	startLine := l.line

	var bc int
	var w int

	var hasEscape bool
	var hasCR bool

	// This loop counts up the number of bytes we should need to store the block string, meaning we
	// avoid many unnecessary allocations. It does slow us down overall, but cuts down memory usage
	// quite a lot.
	for {
		r, w = l.read()
		bc += w

		// Check that escape sequences in this block string are valid.
		if r == bsl {
			r, w = l.read()
			bc += w

			if r == '"' && isTripQuotes(l) {
				bc += 2
				hasEscape = true
			}
		} else if r == cr {
			hasCR = true
		} else if r == eof {
			// Check for end of input.
			return Token{
				Kind:    TokenKindIllegal,
				Literal: "unexpected eof, probably unclosed block string",
				Column:  startLPos - 2,
				Line:    startLine,
			}
		} else if r < ws && r != tab && r != lf && r != cr {
			// Check for invalid characters in the block string, bail early.
			return Token{
				Kind:    TokenKindIllegal,
				Literal: fmt.Sprintf("invalid character within block string: %q", r),
				Column:  startLPos - 2,
				Line:    startLine,
			}
		} else if r == '"' && isTripQuotes(l) {
			// Escape from the loop as early as possible, if we've hit the end of the block string.
			break
		}
	}

	if !hasCR && !hasEscape {
		return Token{
			Kind:    TokenKindStringValue,
			Literal: l.blockStringLiteral(btos(l.input[startPos : l.pos-3])),
			Column:  startLPos - 2,
			Line:    startLine,
		}
	}

	l.pos = startPos
	l.lpos = startLPos
	l.line = startLine

	buf := bytes.NewBuffer(make([]byte, 0, bc))

	for {
		r, _ = l.read()

		if r == lf {
			buf.WriteRune(r)
		} else if r == cr {
			if l.input[l.pos+1] != byte(lf) {
				// Replaces all CR characters with LF characters, only if they're not followed by a
				// LF that we'll write next loop anyway (we don't want a double newline).
				buf.WriteRune(lf)
			}
		} else if r == bsl {
			// Handle escaped sequences appropriately.
			r, _ = l.read()

			if r == '"' && isTripQuotes(l) {
				buf.WriteString(`"""`)
			} else {
				buf.WriteRune(bsl)
				buf.WriteRune(r)
			}
		} else if r == eof {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: "unexpected eof, probably unclosed block string",
				Column:  startLPos - 2,
				Line:    startLine,
			}
		} else if r == '"' && isTripQuotes(l) {
			// Escape from the loop if we've hit the end of the block string.
			break
		} else if r < ws && r != tab && r != lf && r != cr {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: fmt.Sprintf("invalid character within block string: %q", r),
				Column:  startLPos - 2,
				Line:    startLine,
			}
		} else {
			// Write everything else to the buffer.
			buf.WriteRune(r)
		}
	}

	return Token{
		Kind:    TokenKindStringValue,
		Literal: l.blockStringLiteral(btos(buf.Bytes())),
		Column:  startLPos - 2,
		Line:    startLine,
	}
}

// blockStringLiteral takes a raw block string value, trims empty lines of the start and end, and
// removes the common indent from the start of each line.
func (l *Lexer) blockStringLiteral(raw string) string {
	lines := strings.Split(raw, "\n")
	lineCount := len(lines)

	l.line += lineCount - 1

	commonIndent := math.MaxInt64
	for i, line := range lines {
		if i == 0 && lineCount > 1 {
			continue
		}

		indent := leadingWhitespace(line)
		if indent < len(line) && indent < commonIndent {
			commonIndent = indent
			if commonIndent == 0 {
				break
			}
		}
	}

	if commonIndent != math.MaxInt64 && lineCount > 0 {
		for i := 1; i < lineCount; i++ {
			if len(lines[i]) < commonIndent {
				lines[i] = ""
			} else {
				lines[i] = lines[i][commonIndent:]
			}
		}
	}

	start := 0
	end := lineCount

	for start < end && leadingWhitespace(lines[start]) == math.MaxInt64 {
		start++
	}

	for start < end && leadingWhitespace(lines[end-1]) == math.MaxInt64 {
		end--
	}

	return strings.Join(lines[start:end], "\n")
}

// leadingWhitespace returns the index in the given string where the first non-whitespace character
// is found, or returns math.MaxInt64 if there are no non-whitespace characters found.
func leadingWhitespace(str string) int {
	for i, r := range str {
		if r != ' ' && r != '\t' {
			return i
		}
	}

	return math.MaxInt64
}

// isTripQuotes is used if we've just scanned a double-quote and want to test if the next 2
// characters are also double-quotes, returning true if that means we've scanned 3 triple quotes in
// a row, or false otherwise.
func isTripQuotes(l *Lexer) bool {
	r1, w1 := l.read()
	r2, w2 := l.read()
	if r1 == '"' && r2 == '"' {
		return true
	}

	if r1 != eof {
		l.unread(w2)
	}
	if r2 != eof {
		l.unread(w1)
	}
	return false
}

// escapedChar returns the rune that corresponds to an escape sequence that is scanned.
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

		r := unicodeCodePointToRune(r1, r2, r3, r4)
		if r < 0 {
			return 0, fmt.Errorf("invalid character escape sequence: %s", "\\u"+string([]rune{r1, r2, r3, r4}))
		}
		return r, nil
	}

	return 0, fmt.Errorf("invalid character escape sequence: %s", "\\"+string(r))
}

// encodeRune is a copy of the utf8.EncodeRune function, but instead of passing in a byte slice as
// the first argument, a callback is given. This callback may be called multiple times. This allows
// individual bytes to be passed back to the caller, one at a time. This enables the caller to do
// things like encode a rune into an existing byte slice, instead of allocating a new one.
func encodeRune(r rune, cb func(a byte)) {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		cb(byte(r))
		return
	case i <= rune2Max:
		cb(t2 | byte(r>>6))
		cb(tx | byte(r)&maskx)
		return
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = runeError
		fallthrough
	case i <= rune3Max:
		cb(t3 | byte(r>>12))
		cb(tx | byte(r>>6)&maskx)
		cb(tx | byte(r)&maskx)
		return
	default:
		cb(t4 | byte(r>>18))
		cb(tx | byte(r>>12)&maskx)
		cb(tx | byte(r>>6)&maskx)
		cb(tx | byte(r)&maskx)
		return
	}
}

// unicodeCodePointToRune converts 4 hexadecimal characters represented as runes (from read) to a
// single rune that has the value of the unicode code point represented by the 4 hexadecimal
// characters.
//
// See: https://github.com/graphql/graphql-js/blob/84d05fc5c288f2c20df20cf7f60ee356fa6a2cdb/src/language/lexer.js#L689
func unicodeCodePointToRune(ar, br, cr, dr rune) rune {
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
func (l *Lexer) scanComment(r rune) Token {
	var wasCR bool
	var w int

	for {
		r, w = l.read()
		if r == eof {
			return l.Scan()
		}

		// If on the last iteration we saw a CR, then we should check if we just read an LF on this
		// iteration. If we did, reset line position as the next character is still the start of the
		// next line, then scan.
		if wasCR && r == lf {
			l.lpos = 0

			return l.Scan()
		}

		// Otherwise, if we saw a CR, and this rune isn't an LF, then we have started reading the
		// next line's runes, so unread the rune we read, and scan the next token.
		if wasCR && r != lf {
			l.unread(w)

			return l.Scan()
		}

		// If we encounter a CR at any point, this will be true.
		if r == cr {
			// Carriage return, i.e. '\r'.
			l.line++
			l.lpos = 0
			wasCR = true
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
func (l *Lexer) scanName(r rune) Token {
	byteStart := l.pos - 1
	runeStart := l.lpos

Loop:
	for {
		r, w := l.read()

		switch {
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_':
			continue
		case r == eof:
			break Loop
		default:
			l.unread(w)
			break Loop
		}
	}

	return Token{
		Kind:    TokenKindName,
		Literal: btos(l.input[byteStart:l.pos]),
		Column:  runeStart,
		Line:    l.line,
	}
}

// scanPunctuator scans valid GraphQL punctuation tokens.
func (l *Lexer) scanPunctuator(r rune, w int) Token {
	byteStart := l.pos
	runeStart := l.lpos

	if r == '.' {
		r2, _ := l.read()
		r3, _ := l.read()

		rs := []rune{r, r2, r3}
		if rs[1] != '.' || rs[2] != '.' {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: fmt.Sprintf("invalid punctuator, expected \"...\" but got: %q", string(rs)),
				Column:  runeStart,
				Line:    l.line,
			}
		}

		return Token{
			Kind:    TokenKindPunctuator,
			Literal: "...",
			Column:  runeStart,
			Line:    l.line,
		}
	}

	return Token{
		Kind:    TokenKindPunctuator,
		Literal: btos(l.input[byteStart-w : byteStart]),
		Column:  runeStart,
		Line:    l.line,
	}
}

// scanNumber scans valid GraphQL integer and float value tokens.
func (l *Lexer) scanNumber(r rune) Token {
	byteStart := l.pos - 1
	runeStart := l.lpos

	var float bool // If true, number is float.
	var err error  // So no shadowing of r.

	var w int

	// Check for preceding minus sign
	if r == '-' {
		r, w = l.read()
	}

	// Check if digits begins with zero
	if r == '0' {
		r, w = l.read()

		// If there is another digit after zero, error.
		if r >= '0' && r <= '9' {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: fmt.Sprintf("invalid number, unexpected digit after 0: %q", r),
				Column:  runeStart,
				Line:    l.line,
			}
		}

		// If number does not begin with zero, read the digits.
		// If the first character is not a digit, error.
	} else {
		r, w, err = l.readDigits(r)
		if err != nil {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: err.Error(),
				Column:  runeStart,
				Line:    l.line,
			}
		}
	}

	// Check for a decimal place, if there is a decimal place this number is a float.
	if r == '.' {
		float = true

		r, w = l.read()

		// Read the digits after the decimal place if the first character is not a digit, error.
		r, w, err = l.readDigits(r)
		if err != nil {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: err.Error(),
				Column:  runeStart,
				Line:    l.line,
			}
		}
	}

	// Check for exponent sign, if there is an exponent sign this number is a float.
	if r == 'e' || r == 'E' {
		float = true

		r, w = l.read()

		// Check for positive or negative symbol in front of the value.
		if r == '+' || r == '-' {
			r, w = l.read()
		}

		// Read the exponent digits, if the first character is not a digit, error.
		r, w, err = l.readDigits(r)
		if err != nil {
			return Token{
				Kind:    TokenKindIllegal,
				Literal: err.Error(),
				Column:  runeStart,
				Line:    l.line,
			}
		}
	}

	if r != eof {
		l.unread(w)
	}

	t := Token{
		Literal: btos(l.input[byteStart:l.pos]),
		Line:    l.line,
		Column:  runeStart,
	}

	t.Kind = TokenKindIntValue
	if float {
		t.Kind = TokenKindFloatValue
	}

	return t
}

// readDigits reads up until the next non-numeric character in the input.
func (l *Lexer) readDigits(r rune) (rune, int, error) {
	if !(r >= '0' && r <= '9') {
		return eof, 0, fmt.Errorf("invalid number, expected digit but got: %q", r)
	}

	var w int

Loop:
	for {
		r, w = l.read()

		switch {
		case r >= '0' && r <= '9':
			continue
		default:
			// No need to unread here. We actually want to read the character after the numbers.
			break Loop
		}
	}

	return r, w, nil
}

// readNextSignificant reads runes until a "significant" rune is read, i.e. a rune that could be a
// significant token (not whitespace, not tabs, not newlines, not commas, not encoding-specific
// characters, etc.). It also does part of the work for identifying when new lines are encountered
// to increment the line counter.
func (l *Lexer) readNextSignificant() (rune, int) {
	var wasCR bool

	r := er
	w := 0

Loop:
	for r != eof {
		r, w = l.read()

		switch {
		case r == cr:
			// Carriage return, i.e. '\r'.
			l.line++
			l.lpos = 0
			wasCR = true
		case r == lf:
			// Line feed, i.e. '\n'.
			if !wasCR {
				// \r\n is not 2 newlines, so we must check what the last rune was.
				l.line++
				l.lpos = 0
			}
		case r == tab || r == ws || r == com || r == bom:
			// Skip!
		default:
			// Done, this run was significant.
			break Loop
		}
	}

	return r, w
}

// read moves forward in the input, and returns the next rune available. This function also updates
// the position(s) that the lexer keeps track of in the input so the next read continues from where
// the last left off. Returns the EOF rune if we hit the end of the input.
func (l *Lexer) read() (rune, int) {
	if l.pos >= l.inputLen {
		return eof, 0
	}

	r, w := rune(l.input[l.pos]), 1
	if r >= utf8.RuneSelf {
		r, w = utf8.DecodeRune(l.input[l.pos:])
	}

	l.pos += w
	l.lpos++

	return r, w
}

// unread goes back one rune's worth of bytes in the input, changing the
// positions we keep track of.
// Does not currently go back a line.
func (l *Lexer) unread(width int) {
	l.pos -= width

	if l.lpos > 0 {
		l.lpos--
	}
}
