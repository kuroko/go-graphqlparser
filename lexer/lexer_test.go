package lexer

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/token"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
)

var query = []byte("query 0.001 foo { name 12.42e-10 }")

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

func BenchmarkLexer_ScanExt(b *testing.B) {
	src := source.NewSource(&source.Source{
		Body: query,
	})

	for i := 0; i < b.N; i++ {
		lxr := lexer.Lex(src)

		for {
			tok, _ := lxr(0)
			if tok.Kind == lexer.EOF {
				break
			}

			_ = tok
		}
	}
}

func TestLexerScanNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantToken Token
		wantErr   bool
	}{
		// Happy inputs.
		{
			name:  "lone zero is valid",
			input: "0",
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
			input: "-x",
			wantToken: Token{
				Type: token.Illegal,
			},
			wantErr: true,
		},
		{
			name:  "non digit after decimal point is invalid",
			input: "1.x",
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
}

func TestLexerScanComment(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantToken Token
		wantErr   bool
	}{
		// {
		// 	name: "single line comment valid",
		// 	input: `# comment
		// 	query
		// 	`,
		// 	wantToken: Token{
		// 		Type:    token.Name,
		// 		Literal: "query",
		// 	},
		// 	wantErr: false,
		// },
		{
			name:  "single line comment valid",
			input: "# comment" + string(lf) + "foo",
			wantToken: Token{
				Type:     token.Name,
				Literal:  "foo",
				Position: 4,
				Line:     2,
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
				t.Errorf("Lexer.scanComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			isTypeMatch := gotT.Type == tt.wantToken.Type
			isLiteralMatch := gotT.Literal == tt.wantToken.Literal
			if !isTypeMatch || !isLiteralMatch {
				t.Errorf("Lexer.scanComment() = %+v, want %+v", gotT, tt.wantToken)
			}
		})
	}
}
