package lexer

import "testing"

func TestLexer_Foo(t *testing.T) {
	lexer := Lexer{}
	lexer.Foo("Hello, 世界")

	t.Fail()
}
