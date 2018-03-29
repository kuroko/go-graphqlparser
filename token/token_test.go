package token

import (
	"testing"

	"github.com/seeruk/assert"
)

func TestToken_String(t *testing.T) {
	t.Run("should return an appropriate string", func(t *testing.T) {
		tests := []struct {
			expected string
			token    Type
		}{
			{"Illegal", Illegal},
			{"EOF", EOF},
			{"Punctuator", Punctuator},
			{"Name", Name},
			{"IntValue", IntValue},
			{"FloatValue", FloatValue},
			{"StringValue", StringValue},
			{"UnicodeBOM", UnicodeBOM},
			{"WhiteSpace", WhiteSpace},
			{"LineTerminator", LineTerminator},
			{"Comment", Comment},
			{"Comma", Comma},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.token.String())
		}
	})
}
