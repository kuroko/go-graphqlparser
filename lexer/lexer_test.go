package lexer

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/token"
	"github.com/seeruk/assert"
)

func TestLexer(t *testing.T) {
	t.Run("read()", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []rune
		}{
			{input: "Hello, 世界", expected: []rune{'H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界'}},
		}

		for _, test := range tests {
			lexer := New(test.input)

			var runes []rune

			for {
				r := lexer.read()
				if r == eof {
					break
				}

				runes = append(runes, r)
			}

			assert.Equal(t, string(test.expected), string(runes))
		}
	})
}

func TestLexer_Scan(t *testing.T) {
	lxr := New("query foo { name model }")
	tok := lxr.Scan()

	assert.Equal(t, tok.Type, token.Name)
}

func BenchmarkLexer(b *testing.B) {
	input := "Hello, 世界"

	for i := 0; i < b.N; i++ {
		lexer := New(input)

		for {
			r := lexer.read()
			if r == eof {
				break
			}

			_ = r
		}
	}
}

func BenchmarkLexer_Scan(b *testing.B) {
	input := "query foo { name model }"

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
