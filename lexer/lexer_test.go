package lexer

import (
	"bytes"
	"testing"
)

//func TestLexer(t *testing.T) {
//	t.Run("read()", func(t *testing.T) {
//		tests := []struct {
//			input    string
//			expected []rune
//		}{
//			{input: "Hello, 世界", expected: []rune{'H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界'}},
//		}
//
//		for _, test := range tests {
//			lexer := New(strings.NewReader(test.input))
//
//			var runes []rune
//
//			for {
//				r, _ := lexer.read()
//				if r == eof {
//					break
//				}
//
//				runes = append(runes, r)
//			}
//
//			assert.Equal(t, test.expected, runes)
//		}
//	})
//}

func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lexer := New(bytes.NewReader([]byte("Hello, 世界")))

		for {
			r, _ := lexer.read()
			if r == eof {
				break
			}

			_ = r
		}
	}
}
