package language

import (
	"bytes"
	"testing"

	goparser "github.com/graphql-go/graphql/language/parser"
	gosource "github.com/graphql-go/graphql/language/source"
	ast2 "github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

var tsQuery = []byte(`
schema @foo @bar {
	query: Query
}

directive @foo on SCHEMA
directive @bar on SCHEMA
`)

func BenchmarkTypeSystemParser(b *testing.B) {
	tt := []struct {
		name  string
		query []byte
	}{
		{name: "tsQuery", query: tsQuery},
	}

	for _, t := range tt {
		b.Run(t.name, func(b *testing.B) {
			b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
				runBucketdParser(b, t.query)
			})

			b.Run("github.com/vektah/gqlparser", func(b *testing.B) {
				runVektahGQLSchemaParser(b, t.query)
			})
		})
	}
}

func BenchmarkParser(b *testing.B) {
	biggerQuery := bytes.Buffer{}
	for i := 0; i < 10; i++ {
		biggerQuery.Write([]byte(bigQuery))
	}
	ultraMegaQuery := bytes.Buffer{}
	for i := 0; i < 100; i++ {
		ultraMegaQuery.Write(biggerQuery.Bytes())
	}

	tt := []struct {
		name  string
		query []byte
	}{
		//{name: "ultraMegaQuery", query: ultraMegaQuery.Bytes()},
		//{name: "biggerQuery", query: biggerQuery.Bytes()},
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

			b.Run("github.com/vektah/gqlparser", func(b *testing.B) {
				runVektahGQLParser(b, t.query)
			})
		})
	}
}

func runBucketdParser(b *testing.B, query []byte) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		psr := NewParser(query)

		doc, err := psr.Parse()
		if err != nil {
			b.Fatal(err)
		}

		_ = doc
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

func runVektahGQLSchemaParser(b *testing.B, query []byte) {
	qry := string(query)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		source := ast2.Source{
			Name:  "bench",
			Input: qry,
		}

		ast, err := parser.ParseSchema(&source)
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}
