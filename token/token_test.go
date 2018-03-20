package token

import (
	"math"
	"testing"

	"github.com/seeruk/assert"
)

func TestToken_String(t *testing.T) {
	t.Run("should return an appropriate string", func(t *testing.T) {
		tests := []struct {
			expected string
			token    Token
		}{
			{"Illegal", Token{Illegal, ""}},
			{"EOF", Token{EOF, ""}},
			{"Punctuator(...)", Token{Punctuator, "..."}},
			{"Name(vehicles)", Token{Name, "vehicles"}},
			{"IntValue(4249)", Token{IntValue, "4249"}},
			{"FloatValue(3.141)", Token{FloatValue, "3.141"}},
			{"StringValue(\"Bucketd\")", Token{StringValue, "\"Bucketd\""}},
			{"UnicodeBOM", Token{UnicodeBOM, ""}},
			{"WhiteSpace", Token{WhiteSpace, ""}},
			{"LineTerminator", Token{LineTerminator, ""}},
			{"Comment(# An example)", Token{Comment, "# An example"}},
			{"Comma", Token{Comma, ""}},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.token.String())
		}
	})

	t.Run("should panic if an invalid token is given", func(t *testing.T) {
		var recovered bool

		test := func() {
			defer func() {
				if err := recover(); err != nil {
					recovered = true
				}
			}()

			tok := Token{math.MaxInt64, ""}
			tok.String()
		}

		test()

		assert.True(t, recovered, "expected a panic to recover from")
	})
}
