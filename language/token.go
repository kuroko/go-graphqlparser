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

// TokenKind represents a type of token. The types are predefined as constants.
type TokenKind int

// String returns the name of this type.
func (t TokenKind) String() string {
	return TokenKindNames[t]
}
