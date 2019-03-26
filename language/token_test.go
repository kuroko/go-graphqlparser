package language_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/language"
	"github.com/stretchr/testify/assert"
)

func TestToken_String(t *testing.T) {
	t.Run("should return an appropriate string", func(t *testing.T) {
		tests := []struct {
			expected string
			token    language.TokenKind
		}{
			{"Illegal", language.TokenKindIllegal},
			{"EOF", language.TokenKindEOF},
			{"Punctuator", language.TokenKindPunctuator},
			{"Name", language.TokenKindName},
			{"IntValue", language.TokenKindIntValue},
			{"FloatValue", language.TokenKindFloatValue},
			{"StringValue", language.TokenKindStringValue},
			{"UnicodeBOM", language.TokenKindUnicodeBOM},
			{"WhiteSpace", language.TokenKindWhiteSpace},
			{"LineTerminator", language.TokenKindLineTerminator},
			{"Comment", language.TokenKindComment},
			{"Comma", language.TokenKindComma},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.token.String())
		}
	})
}
