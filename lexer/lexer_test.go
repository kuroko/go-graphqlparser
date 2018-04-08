package lexer

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/lexer/benchutil/graphql-gophers"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/graphql-go/graphql/language/lexer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/stretchr/testify/assert"
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

	scr = string(cr)
	slf = string(lf)
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

		// Scan
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
				tok, err := lxr.Scan()

				actual.Tokens = append(actual.Tokens, tok)
				if err != nil {
					actual.Errors = append(actual.Errors, err.Error())
				}

				if err != nil || tok.Type == token.EOF {
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
