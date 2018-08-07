package ast

// Sdump prints the given AST document as a GraphQL document string, allowing the parser to be
// easily validated against some given input. In fact, if formatted the same, the output of this
// function should match the input query given to the parser to produce the AST.
func Sdump(doc Document) string {
	return ""
}
