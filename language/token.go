package language

// All different lexical token types.
const (
	TokenKindIllegal TokenKind = iota - 1
	TokenKindEOF

	// Lexical Tokens.
	TokenKindPunctuator
	TokenKindName
	TokenKindIntValue
	TokenKindFloatValue
	TokenKindStringValue

	// Ignored tokens.
	TokenKindUnicodeBOM
	TokenKindWhiteSpace
	TokenKindLineTerminator
	TokenKindComment
	TokenKindComma
)

// TokenKindNames is a map of token types to their names as strings.
var TokenKindNames = map[TokenKind]string{
	TokenKindIllegal:        "Illegal",
	TokenKindEOF:            "EOF",
	TokenKindPunctuator:     "Punctuator",
	TokenKindName:           "Name",
	TokenKindIntValue:       "IntValue",
	TokenKindFloatValue:     "FloatValue",
	TokenKindStringValue:    "StringValue",
	TokenKindUnicodeBOM:     "UnicodeBOM",
	TokenKindWhiteSpace:     "WhiteSpace",
	TokenKindLineTerminator: "LineTerminator",
	TokenKindComment:        "Comment",
	TokenKindComma:          "Comma",
}

// Token represents a small, easily categorisable data structure that is fed to the parser to
// produce the abstract syntax tree (AST).
type Token struct {
	Kind     TokenKind // The token type.
	Literal  string    // The literal value consumed.
	Position int       // The starting position, in runes, of this token in the input.
	Line     int       // The line number at the start of this item.
}

// TokenKind represents a type of token. The types are predefined as constants.
type TokenKind int

// String returns the name of this type.
func (t TokenKind) String() string {
	return TokenKindNames[t]
}
