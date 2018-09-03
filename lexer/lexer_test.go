package lexer

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/bucketd/go-graphqlparser/token"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/ast"
	lexer2 "github.com/vektah/gqlparser/lexer"
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
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

						Welcome to GraphQL.
						Let's make this string a little bigger then. Because the larger this string
						becomes, the more efficient our lexer should look...

					From, Bucketd
				"""
				readTime: 2.742
			)
		}
	`

	scr = string(cr)
	slf = string(lf)
)

func BenchmarkLexer(b *testing.B) {
	qry := []byte(query)

	b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lxr := New(qry)

			for {
				tok := lxr.Scan()
				if tok.Type == token.Illegal {
					b.Fatal(tok.Literal)
				}

				if tok.Type == token.EOF || tok.Type == token.Illegal {
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

	b.Run("github.com/vektah/gqlparser", func(b *testing.B) {
		input := string(qry)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			lxr := lexer2.New(&ast.Source{
				Name:  "bench",
				Input: input,
			})

			for {
				tok, err := lxr.ReadToken()
				if err != nil {
					b.Fatal(err)
				}

				if tok.Kind == lexer2.EOF {
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
		Tokens []Token
		Errors []string
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.index), func(t *testing.T) {
			lxr := New([]byte(test.input))
			actual := record{
				Input: test.input,
			}

			for {
				tok := lxr.Scan()

				actual.Tokens = append(actual.Tokens, tok)
				if tok.Type == token.Illegal {
					actual.Errors = append(actual.Errors, tok.Literal)
				}

				if tok.Type == token.EOF || tok.Type == token.Illegal {
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

func TestLexerReadUnread(t *testing.T) {
	// We need to test what happens when we have bytes containing runes of different lengths when we
	// do reads and unreads, so we know we go backwards and forwards the right number of bytes.
	bs := []byte("世h界e界l界l界o")
	l := New(bs)

	r, w1 := l.read()
	assert.Equal(t, fmt.Sprintf("%q", '世'), fmt.Sprintf("%q", r))

	r, w2 := l.read()
	assert.Equal(t, fmt.Sprintf("%q", 'h'), fmt.Sprintf("%q", r))

	l.unread(w2)
	r, w2 = l.read()
	assert.Equal(t, fmt.Sprintf("%q", 'h'), fmt.Sprintf("%q", r))

	l.unread(w2)
	l.unread(w1)
	r, _ = l.read()
	assert.Equal(t, fmt.Sprintf("%q", '世'), fmt.Sprintf("%q", r))
}

func TestEncodeRune(t *testing.T) {
	thing := '😀'

	var counter int
	var bs []byte

	encodeRune(thing, func(b byte) {
		counter++
		bs = append(bs, b)
	})
	if counter != 4 {
		t.Errorf("expected emoji triggers default case, counter should be 4, actual: %d\n", counter)
	}

	tbs := []byte(string(thing))
	if !reflect.DeepEqual(bs, tbs) {
		t.Errorf(
			"\nexpected: %x %x %x %x\ngot: %x %x %x %x\n",
			tbs[0], tbs[1], tbs[2], tbs[3],
			bs[0], bs[1], bs[2], bs[3],
		)
	}
}
