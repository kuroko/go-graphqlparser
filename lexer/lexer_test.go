package lexer

import (
	"bytes"
	"testing"

	"github.com/bucketd/go-graphqlparser/token"
)

func BenchmarkLexer_Scan(b *testing.B) {
	input := bytes.NewReader([]byte("query foo { name model }"))

	for i := 0; i < b.N; i++ {
		lxr := New(input)

		for {
			tok := lxr.Scan()
			if tok.Type == token.EOF {
				break
			}

			_ = tok
		}
	}
}
