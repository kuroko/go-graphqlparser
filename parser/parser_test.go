package parser

import (
	"bytes"
	"strings"
	"testing"

	goparser "github.com/graphql-go/graphql/language/parser"
	gosource "github.com/graphql-go/graphql/language/source"
	ast2 "github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"

	gophersquery "github.com/bucketd/go-graphqlparser/benchutil/graphql-gophers/query"
)

func BenchmarkParser(b *testing.B) {
	ultraMegaQuery := bytes.Buffer{}
	for i := 0; i < 100; i++ {
		ultraMegaQuery.Write(monsterQuery)
	}

	tt := []struct {
		name  string
		query []byte
	}{
		{name: "ultraMegaQuery", query: ultraMegaQuery.Bytes()},
		{name: "monsterQuery", query: monsterQuery},
		{name: "normalQuery", query: normalQuery},
		{name: "tinyQuery", query: tinyQuery},
	}

	for _, t := range tt {
		b.Run(t.name, func(b *testing.B) {
			b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
				runBucketdParser(b, t.query)
			})

			b.Run("github.com/graphql-go/graphql", func(b *testing.B) {
				runGraphQLGoParser(b, t.query)
			})

			b.Run("github.com/graphql-gophers/graphql-go", func(b *testing.B) {
				runGraphQLGophersParser(b, t.query)
			})

			b.Run("github.com/vektah/gqlparser", func(b *testing.B) {
				runVektahGQLParser(b, t.query)
			})
		})
	}
}

func runBucketdParser(b *testing.B, query []byte) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		psr := New(query)

		ast, err := psr.Parse()
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}

func runGraphQLGoParser(b *testing.B, query []byte) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		params := goparser.ParseParams{
			Source: gosource.NewSource(&gosource.Source{
				Body: query,
			}),
		}

		ast, err := goparser.Parse(params)
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}

func runGraphQLGophersParser(b *testing.B, query []byte) {
	// No multi-line string support in this parser...
	qry := string(query)
	qry = strings.Replace(qry, `"""
        Hello,
            World!

        Yours,
            GraphQL
    """`, `"Hello,\n    World!\n\nYours,    GraphQL"`, -1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ast, err := gophersquery.Parse(qry)
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}

func runVektahGQLParser(b *testing.B, query []byte) {
	qry := string(query)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		source := ast2.Source{
			Name:  "bench",
			Input: qry,
		}

		ast, err := parser.ParseQuery(&source)
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}
