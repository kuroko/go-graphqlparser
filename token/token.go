// Package token contains the tokens used for lexical analysis of GraphQL documents. Many of the
// tokens names are based on the official GraphQL specification (Working Draft, October 2016), which
// can be found at http://facebook.github.io/graphql/October2016.
package token

// All different lexical token types.
const (
	Illegal Type = iota - 1
	EOF

	// Lexical Tokens.
	Punctuator
	Name
	IntValue
	FloatValue
	StringValue

	// Ignored tokens.
	UnicodeBOM
	WhiteSpace
	LineTerminator
	Comment
	Comma
)

// TypeNames is a map of token types to their names as strings.
var TypeNames = map[Type]string{
	Illegal:        "Illegal",
	EOF:            "EOF",
	Punctuator:     "Punctuator",
	Name:           "Name",
	IntValue:       "IntValue",
	FloatValue:     "FloatValue",
	StringValue:    "StringValue",
	UnicodeBOM:     "UnicodeBOM",
	WhiteSpace:     "WhiteSpace",
	LineTerminator: "LineTerminator",
	Comment:        "Comment",
	Comma:          "Comma",
}

// Type represents a type of token. The types are predefined as constants.
type Type int

// String returns the name of this type.
func (t Type) String() string {
	return TypeNames[t]
}
