package language

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToken_String(t *testing.T) {
	t.Run("should return an appropriate string", func(t *testing.T) {
		tests := []struct {
			expected string
			token    TokenKind
		}{
			{"Illegal", TokenKindIllegal},
			{"EOF", TokenKindEOF},
			{"Punctuator", TokenKindPunctuator},
			{"Name", TokenKindName},
			{"IntValue", TokenKindIntValue},
			{"FloatValue", TokenKindFloatValue},
			{"StringValue", TokenKindStringValue},
			{"UnicodeBOM", TokenKindUnicodeBOM},
			{"WhiteSpace", TokenKindWhiteSpace},
			{"LineTerminator", TokenKindLineTerminator},
			{"Comment", TokenKindComment},
			{"Comma", TokenKindComma},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.token.String())
		}
	})
}
