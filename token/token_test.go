package token

import "testing"

func TestFoo(t *testing.T) {
	t.Log(WhiteSpace)
	t.Log(Comma)
	t.Log(UnicodeBOM)
	t.Error(FloatValue)
}
