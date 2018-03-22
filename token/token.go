// Package token contains the tokens used for lexical analysis of GraphQL documents. Many of the
// tokens names are based on the official GraphQL specification (Working Draft, October 2016), which
// can be found at http://facebook.github.io/graphql/October2016.
package token

// TODO(elliot): Go strings are UTF-8, but reading a file actually allows us to read bytes. For now,
// let's just assume that everything is UTF-8, otherwise it's going to be far more complex. Given
// that each character should be a UTF-8 character, runes are an easy choice for the tokens below.
// Each `\uXXXX` value can easily be made into a rune by using `rune(0xXXXX)`.

// TODO(elliot): Each token will likely start off just reading the first character. If for example,
// there was just a `$` on it's own, we'd instantly know that a punctuator was being read. On the
// other hand, if a `4` was read, then we'd have to keep reading it until we knew if it was an int
// value or a float value. This will not use regular expressions though, like the comments below
// might suggest, as that might be too slow.

const (
	Illegal Type = iota
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
func (t *Type) String() string {
	return TypeNames[*t]
}
