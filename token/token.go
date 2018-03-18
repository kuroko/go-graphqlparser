// Package token contains the tokens used for lexical analysis of GraphQL documents. Many of the
// tokens names are based on the official GraphQL specification (Working Draft, October 2016), which
// can be found at http://facebook.github.io/graphql/October2016.
package token

// TODO(elliot): Go strings are UTF-8, but reading a file actually allows us to read bytes. For now,
// let's just assume that everything is UTF-8, otherwise it's going to be far more complex. Given
// that each character should be a UTF-8 character, runes are an easy choice for the tokens below.
// Each `\uXXXX` value can easily be made into a rune by using `rune(0xXXXX)`.

const (
	Illegal = iota
	EOF

	// 2.1.6: Lexical Tokens.
	Punctuator
	Name
	IntValue
	FloatValue
	StringValue

	// 2.1.7: Ignored Tokens
	UnicodeBOM // "\uFEFF"
	WhiteSpace // "\u0009" || "\u0020"
	LineTerminator
	Comment
	Comma
)
