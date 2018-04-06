package lexer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bucketd/go-graphqlparser/token"
	"github.com/seeruk/assert"
)

var query = []byte("query \"\\u4e16\" 0.001 foo { name 12.42e-10 }")

func BenchmarkLexer_Scan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lxr := New(query)

		for {
			tok, _ := lxr.Scan()
			if tok.Type == token.EOF {
				break
			}

			_ = tok
		}
	}
}

func TestLexer_Scan(t *testing.T) {
	t.Run("scanNumber()", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantToken Token
			wantErr   bool
		}{
			// Happy inputs.
			{
				name:  "lone zero is valid",
				input: "0 ", // Q: padding ws for coverage?
				wantToken: Token{
					Type:    token.IntValue,
					Literal: "0",
				},
				wantErr: false,
			},
			{
				name:  "positive int is valid",
				input: "123456789",
				wantToken: Token{
					Type:    token.IntValue,
					Literal: "123456789",
				},
				wantErr: false,
			},
			{
				name:  "negative int is valid",
				input: "-123456789",
				wantToken: Token{
					Type:    token.IntValue,
					Literal: "-123456789",
				},
				wantErr: false,
			},
			{
				name:  "positive float is valid",
				input: "1.1",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "1.1",
				},
				wantErr: false,
			},
			{
				name:  "negative float is valid",
				input: "-1.1",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "-1.1",
				},
				wantErr: false,
			},
			{
				name:  "exponent is valid",
				input: "10E99",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "10E99",
				},
				wantErr: false,
			},
			{
				name:  "negative exponent is valid",
				input: "-10E99",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "-10E99",
				},
				wantErr: false,
			},
			{
				name:  "float, positive exponent is valid",
				input: "1.1e99",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "1.1e99",
				},
				wantErr: false,
			},
			{
				name:  "float, negative exponent is valid",
				input: "-1.1e-99",
				wantToken: Token{
					Type:    token.FloatValue,
					Literal: "-1.1e-99",
				},
				wantErr: false,
			},
			// Errorful inputs.
			{
				name:  "negative symbol with no following digits is invalid",
				input: "-",
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "negative symbol with non digit after is invalid",
				input: "-界",
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "non digit after decimal point is invalid",
				input: "1.界",
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "zero followed by number is invalid",
				input: "01",
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "non digit after exponent",
				input: "-1.1e界",
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bs := []byte(tt.input)
				l := New(bs)
				gotT, err := l.Scan()
				if !tt.wantErr && err != nil {
					t.Errorf("Lexer.scanNumber() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				isTypeMatch := gotT.Type == tt.wantToken.Type
				isLiteralMatch := gotT.Literal == tt.wantToken.Literal
				if !isTypeMatch || !isLiteralMatch {
					t.Errorf("Lexer.scanNumber() = %+v, want %+v", gotT, tt.wantToken)
				}
			})
		}
	})

	t.Run("scanName()", func(t *testing.T) {

	})

	t.Run("scanPunctuator()", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantToken Token
			wantErr   bool
		}{
			// Happy inputs.
			{
				name:  "standard punctuation is valid",
				input: " { ", // Q: padding ws for coverage?
				wantToken: Token{
					Type:     token.Punctuator,
					Literal:  "{",
					Position: 2,
					Line:     1,
				},

				wantErr: false,
			},
			{
				name:  "ellipsis is valid",
				input: " ... ", // Q: padding ws for coverage?
				wantToken: Token{
					Type:     token.Punctuator,
					Literal:  "...",
					Position: 2,
					Line:     1,
				},

				wantErr: false,
			},
			// Errorful inputs.
			{
				name:  " period followed by not two more periods is invalid",
				input: " .界界 ", // Q: padding ws for coverage?
				wantToken: Token{
					Type:    token.Illegal,
					Literal: "",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bs := []byte(tt.input)
				l := New(bs)
				gotT, err := l.Scan()
				if !tt.wantErr && err != nil {
					t.Errorf("Lexer.scanPunctuator() error = %+v, wantErr %+v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotT, tt.wantToken) {
					t.Errorf("Lexer.scanPunctuator() = %+v, want %+v", gotT, tt.wantToken)
				}
			})
		}
	})

	t.Run("scanComment()", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantToken Token
			wantErr   bool
		}{
			{
				name:  "only comment, no newlines",
				input: "# only comment",
				wantToken: Token{
					Type:     token.EOF,
					Literal:  "",
					Position: 14,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "single line comment valid + lf",
				input: "# comment" + string(lf) + "foo",
				wantToken: Token{
					Type:     token.Name,
					Literal:  "foo",
					Position: 1,
					Line:     2,
				},
				wantErr: false,
			},
			{
				name:  "single line comment valid + cr",
				input: "# comment" + string(cr) + "foo",
				wantToken: Token{
					Type:     token.Name,
					Literal:  "foo",
					Position: 1,
					Line:     2,
				},
				wantErr: false,
			},
			{
				name:  "multi-line comment valid",
				input: "# line 1" + string(lf) + "# line 2" + string(lf) + "query",
				wantToken: Token{
					Type:     token.Name,
					Literal:  "query",
					Position: 1,
					Line:     3,
				},
			},
			{
				name:  "cr + lf only one extra line",
				input: "# comment" + string(cr) + string(lf) + "foo",
				wantToken: Token{
					Type:     token.Name,
					Literal:  "foo",
					Position: 1,
					Line:     2,
				},
			},
			{
				name:  "lf + cr two extra lines",
				input: "# comment" + string(lf) + string(cr) + "foo",
				wantToken: Token{
					Type:     token.Name,
					Literal:  "foo",
					Position: 1,
					Line:     3,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bs := []byte(tt.input)
				l := New(bs)
				gotT, err := l.Scan()
				if !tt.wantErr && err != nil {
					t.Errorf("Lexer.scanComment() error = %+v, wantErr %+v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotT, tt.wantToken) {
					t.Errorf("Lexer.scanComment() = %+v, want %+v", gotT, tt.wantToken)
				}
			})
		}
	})

	t.Run("scanString()", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantToken Token
			wantErr   bool
		}{
			// Happy
			{
				name:  "simple string",
				input: `"simple"`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  "simple",
					Position: 1,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "escaped lf",
				input: `"line\nfeed"`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  "line" + string(lf) + "feed",
					Position: 1,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "world",
				input: `"\u4e16world"`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  "世world",
					Position: 1,
					Line:     1,
				},
				wantErr: false,
			},
			// Errorful
			{
				name:  "no closing quote",
				input: `"foo`,
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "invalid unicode",
				input: `"\uAZ"`,
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "invalid escape",
				input: `"\z"`,
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
			{
				name:  "newline errors",
				input: `"foo` + string(lf) + `"`,
				wantToken: Token{
					Type: token.Illegal,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bs := []byte(tt.input)
				l := New(bs)
				gotT, err := l.Scan()
				if !tt.wantErr && err != nil {
					t.Errorf("Lexer.scanString() error = %+v, wantErr %+v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotT, tt.wantToken) {
					t.Errorf("Lexer.scanString() = %+v, want %+v", gotT, tt.wantToken)
				}
			})
		}
	})

	t.Run("scanBlockString()", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantToken Token
			wantErr   bool
		}{
			// Happy
			{
				name:  "simple block string",
				input: `"""simple"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  "simple",
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "block string with nested quotes",
				input: `"""nested "" quotes"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `nested "" quotes`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "block string with nested quotes",
				input: `"""nested "" quotes"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `nested "" quotes`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "block string escaped triple quotes",
				input: `"""nested \\""" quotes"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `nested """ quotes`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bs := []byte(tt.input)
				l := New(bs)
				gotT, err := l.Scan()
				if !tt.wantErr && err != nil {
					t.Errorf("Lexer.scanString() error = %+v, wantErr %+v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotT, tt.wantToken) {
					t.Errorf("Lexer.scanString() = %+v, want %+v", gotT, tt.wantToken)
				}
			})
		}
	})
}

func TestLexerReadUnread(t *testing.T) {
	bs := []byte("世h界e界l界l界o")
	l := New(bs)

	r, _ := l.read()
	assert.Equal(t, fmt.Sprintf("%q", r), fmt.Sprintf("%q", '世'))

	r, _ = l.read()
	assert.Equal(t, fmt.Sprintf("%q", r), fmt.Sprintf("%q", 'h'))

	l.unread()
	r, _ = l.read()
	assert.Equal(t, fmt.Sprintf("%q", r), fmt.Sprintf("%q", 'h'))

	l.unread()
	l.unread()
	r, _ = l.read()
	assert.Equal(t, fmt.Sprintf("%q", r), fmt.Sprintf("%q", '世'))
}
