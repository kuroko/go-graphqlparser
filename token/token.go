// Package token contains the tokens used for lexical analysis of GraphQL documents. Many of the
// tokens names are based on the official GraphQL specification (Working Draft, October 2016), which
// can be found at http://facebook.github.io/graphql/October2016.
package token

const (
	Illegal Type = iota - 1
	EOF

	// 2.1.6: Lexical Tokens.
	Punctuator  // One of: !, $, (, ), ..., :, =, @, [, ], {, |, }
	Name        // Regex: /[_A-Za-z][_0-9A-Za-z]*/
	IntValue    // Regex: /^[+-]?(0|[1-9]+)$/
	FloatValue  // Regex: /^[+-]?(0|[1-9]+)(\.[0-9]+)?([eE][+-]?[0-9]+)?$/
	StringValue // Trickier: "", "some string", "\u1234", "\b\f\n\r\t", "this\nis_a%$string\u0020"

	// 2.1.7: Ignored Tokens
	UnicodeBOM     // Unicode: "\uFEFF"
	WhiteSpace     // Unicode: "\u0009" || "\u0020". Only significant in strings
	LineTerminator // Unicode: "\u000A", "\u000D". Only significant in strings
	Comment        // Literal: "#". Consume up to next LineTerminator
	Comma          // Literal: ","
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
