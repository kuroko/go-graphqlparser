package language_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/bucketd/go-graphqlparser/language"
	"github.com/stretchr/testify/assert"

	glexer "github.com/graphql-go/graphql/language/lexer"
	gsource "github.com/graphql-go/graphql/language/source"
	vast "github.com/vektah/gqlparser/ast"
	vlexer "github.com/vektah/gqlparser/lexer"
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
                content: """
                    Hello,

                        Welcome to GraphQL.
                        Let's make this string a little bigger then. Because the larger this string
                        becomes, the more we can test our lexer.

                    From, Bucketd
                """
                readTime: 2.742
            ) @async(bar: $baz)
        }
    `

	lf  = rune(0x000A)
	cr  = rune(0x000D)
	scr = string(cr)
	slf = string(lf)
)

func BenchmarkLexer(b *testing.B) {
	qry := []byte(query)

	b.Run("bucketd", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lxr := language.NewLexer(qry)

			for {
				tok := lxr.Scan()
				if tok.Kind == language.TokenKindIllegal {
					b.Fatal(tok.Literal)
				}

				if tok.Kind == language.TokenKindEOF {
					break
				}

				_ = tok
			}
		}
	})

	b.Run("graphql-go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lxr := glexer.Lex(gsource.NewSource(&gsource.Source{
				Body: qry,
			}))

			for {
				tok, err := lxr(0)
				if err != nil {
					b.Fatal(err)
				}

				if tok.Kind == glexer.EOF {
					break
				}

				_ = tok
			}
		}
	})

	b.Run("vektah", func(b *testing.B) {
		input := string(qry)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			lxr := vlexer.New(&vast.Source{
				Name:  "bench",
				Input: input,
			})

			for {
				tok, err := lxr.ReadToken()
				if err != nil {
					b.Fatal(err)
				}

				if tok.Kind == vlexer.EOF {
					break
				}

				_ = tok
			}
		}
	})
}

func TestLexer_ScanGolden(t *testing.T) {
	tests := []struct {
		index string
		input string
	}{
		// scanNumber
		{"001", "0"},
		{"002", "123456789"},
		{"003", "-123456789"},
		{"004", "1.1"},
		{"005", "-1.1"},
		{"006", "10E99"},
		{"007", "-10E99"},
		{"008", "1.1e99"},
		{"009", "-1.1e-99"},
		{"010", "-"},
		{"011", "-界"},
		{"012", "1.界"},
		{"013", "01"},
		{"014", "-1.1e界"},

		// scanPunctuator
		{"101", "!$()...:=@[]{|}"},
		{"102", " ! $ ( ) ... : = @ [ ] { | }"},
		{"103", ".界界"},

		// scanComment
		{"201", "# only comment"},
		{"202", "# line 1" + scr + "foo"},
		{"203", "# line 1" + slf + "foo"},
		{"204", "# line 1" + scr + slf + "foo"},
		{"205", "# line 1" + slf + scr + "foo"},
		{"206", "# line 1" + scr + "# line 2" + scr + "query"},
		{"207", "# line 1" + slf + "# line 2" + slf + "query"},
		{"208", "# line 1" + scr + slf + "# line 2" + scr + slf + "query"},

		// scanString
		{"301", `"foo"`},
		{"302", `"foo \n bar"`},
		{"303", `"foo \u4e16"`},
		{"304", `"foo`},
		{"305", `"\""`},
		{"306", `"\\"`},
		{"307", `"\/"`},
		{"308", `"\b"`},
		{"309", `"\f"`},
		{"310", `"\n"`},
		{"311", `"\r"`},
		{"312", `"\t"`},
		{"313", `"\uAZ"`},
		{"314", `"\z"`},
		{"315", `"foo` + slf + `"`},
		{"316", `"foo` + scr + slf + `"`},
		{"317", `"\uAzAz"`},
		{"318", `"\u0080"`},
		{"319", `"\uFFFF"`},
		{"320", `"😀"`},
		{"321", `"\uD800"`},

		// scanBlockString
		{"401", `""""""`},
		{"402", `"""foo"""`},
		{"403", `"""foo "" bar"""`},
		{"404", `"""\t \u1234"""`},
		{"405", `"""` + scr + `foo` + scr + `"""`},
		{"406", `"""` + slf + `foo` + slf + `"""`},
		{"407", `"""` + scr + slf + `foo` + scr + slf + `"""`},
		{"408", `"""` + slf + scr + `foo` + slf + scr + `"""`},
		{"409", `"""` + scr + `\"""` + scr + `"""`},
		{"410", `"""` + slf + `\"""` + slf + `"""`},
		{"411", `"""` + scr + slf + `\"""` + scr + slf + `"""`},
		{"412", `"""` + slf + scr + `\"""` + slf + scr + `"""`},
		{"413", `"""foo \""" bar"""`},
		{"414", `"""foo \u1234 " \""""""`},
		{"415", `""""`},
		{"416", `"""""`},
		{"417", `"""\"""""`},
		{"418", `"""\u1234""`},
		{"419", `"""😀"""`},
		{"420", `"""\""a"""`},

		// Scan
		{"998", `query foo ` + string(rune(128515)) + ` { bar baz }`},
		{"999", query},
	}

	type record struct {
		Input  string
		Tokens []language.Token
		Errors []string
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.index), func(t *testing.T) {
			lxr := language.NewLexer([]byte(test.input))
			actual := record{
				Input: test.input,
			}

			for {
				tok := lxr.Scan()

				actual.Tokens = append(actual.Tokens, tok)
				if tok.Kind == language.TokenKindIllegal {
					actual.Errors = append(actual.Errors, tok.Literal)
					break
				}

				if tok.Kind == language.TokenKindEOF {
					break
				}
			}

			goldenFileName := fmt.Sprintf("testdata/ScanGolden.%s.json", test.index)

			if *update {
				bs, err := json.MarshalIndent(actual, "", "  ")
				if err != nil {
					t.Error(err)
				}

				err = ioutil.WriteFile(goldenFileName, bs, 0666)
				if err != nil {
					t.Error(err)
				}

				return
			}

			goldenBs, err := ioutil.ReadFile(goldenFileName)
			if err != nil {
				t.Error(err)
			}

			expected := record{}

			err = json.Unmarshal(goldenBs, &expected)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, expected, actual)
		})
	}
}
