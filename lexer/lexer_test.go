package lexer

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/lexer/benchutil/graphql-gophers"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/seeruk/assert"
)

var update = flag.Bool("update", false, "update golden record files?")

const (
	// query is a valid GraphQL query that contains at least one of each token type. It's re-used to
	// compare against other GraphQL libraries.
	query = `
		# Mutation for testing different token types.
		mutation {
			createPost(
				id: 1024
				title: "String Value"
				content: """Block string value isn't supported by all libs."""
				readTime: 2.742
			)
		}
	`
)

func BenchmarkLexer(b *testing.B) {
	qry := []byte(query)

	b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lxr := New(qry)

			for {
				tok, err := lxr.Scan()
				if err != nil {
					b.Fatal(err)
				}

				if tok.Type == token.EOF {
					break
				}

				_ = tok
			}
		}
	})

	b.Run("github.com/graphql-go/graphql", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lxr := lexer.Lex(source.NewSource(&source.Source{
				Body: qry,
			}))

			for {
				tok, err := lxr(0)
				if err != nil {
					b.Fatal(err)
				}

				if tok.Kind == lexer.EOF {
					break
				}

				_ = tok
			}
		}
	})

	b.Run("github.com/graphql-gophers/graphql-go", func(b *testing.B) {
		// This lib doesn't support block strings currently.
		qry := strings.Replace(query, `"""`, `"`, -1)

		for i := 0; i < b.N; i++ {
			lxr := graphql_gophers.NewLexer(qry)
			lxr.Consume()

			// This lexer is a little more fiddly to bench, we have to know the expected structure
			// of a query to call the right lexer methods in the right order:
			_ = lxr.ConsumeIdent()   // mutation
			lxr.ConsumeToken('{')    //
			_ = lxr.ConsumeIdent()   // createPost
			lxr.ConsumeToken('(')    //
			_ = lxr.ConsumeIdent()   // id
			lxr.ConsumeToken(':')    //
			_ = lxr.ConsumeLiteral() // 1024
			_ = lxr.ConsumeIdent()   // title
			lxr.ConsumeToken(':')    //
			_ = lxr.ConsumeLiteral() // "String Value"
			_ = lxr.ConsumeIdent()   // content
			lxr.ConsumeToken(':')    //
			_ = lxr.ConsumeLiteral() // "Block string value isn't supported by everything."
			_ = lxr.ConsumeIdent()   // readTime
			lxr.ConsumeToken(':')    //
			_ = lxr.ConsumeLiteral() // 2.742
			lxr.ConsumeToken(')')    //
			lxr.ConsumeToken('}')    //
		}
	})
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
				name:  "block string escaped triple quotes",
				input: `"""nested ` + string(bsl) + string(dq) + string(dq) + string(dq) + ` quotes"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `nested """ quotes`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "non trip quote escapes ignored",
				input: `"""\t \u1234"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `\t \u1234`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name: "ignore leading and trailing newlines",
				input: `"""
ignore leading and trailing newlines
"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `ignore leading and trailing newlines`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "escape sequences with escaped triple quotes",
				input: `"""\u1234 " \""""""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `\u1234 " """`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "leading trailing newlines with escaped triple quotes",
				input: `"""` + string(cr) + string(lf) + `\"""` + string(cr) + string(lf) + `"""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  `"""`,
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			{
				name:  "empty triple quotes",
				input: `""""""`,
				wantToken: Token{
					Type:     token.StringValue,
					Literal:  "",
					Position: 3,
					Line:     1,
				},
				wantErr: false,
			},
			// Errorful
			{
				name:  "not closing properly",
				input: `"""""`,
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
					t.Errorf("Lexer.scanBlockString() error = %#v, wantErr %#v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotT, tt.wantToken) {
					t.Errorf("Lexer.scanBlockString() = %#v, want %#v", gotT, tt.wantToken)
				}
			})
		}
	})
}

func TestLexerReadUnread(t *testing.T) {
	// We need to test what happens when we have bytes containing runes of different lengths when we
	// do reads and unreads, so we know we go backwards and forwards the right number of bytes.
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
