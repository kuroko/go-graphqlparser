package language

import (
	"bytes"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/parser"

	goparser "github.com/graphql-go/graphql/language/parser"
	gosource "github.com/graphql-go/graphql/language/source"
	vast "github.com/vektah/gqlparser/ast"
)

var tsQuery = []byte(`
schema @foo @bar {
	query: Query
}

directive @foo on SCHEMA
directive @bar on SCHEMA

"UUID is a scalar that represents the string form (36 characters) of a UUID"
scalar UUID

"""
Field4Input's description, using a block string.
This one spans multiple lines.
"""
input Field4Input {
	input1: ID!
}

"Field4Payload's description, using a single-line string"
type Field4Payload {
	field1: UUID!
	field2: String!
}

type Query {
	"Fields may also have comments"
	field1: String
	field2(arg1: Int): String
	field3: [String!]
	field4(in: Field4Input!): Field4Payload
}
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
			b.Run("bucketd", func(b *testing.B) {
				runBucketdParser(b, t.query)
			})

			b.Run("vektah", func(b *testing.B) {
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
			b.Run("bucketd", func(b *testing.B) {
				runBucketdParser(b, t.query)
			})

			b.Run("graphql-go", func(b *testing.B) {
				runGraphQLGoParser(b, t.query)
			})

			b.Run("vektah", func(b *testing.B) {
				runVektahGQLParser(b, t.query)
			})
		})
	}
}

func TestParser_Parse(t *testing.T) {
	t.Run("should set Location on Definition", func(t *testing.T) {
		query := []byte(`
			# Really testing this location stuff here...
			{ hello }
		`)

		psr := NewParser(query)

		doc, err := psr.Parse()
		require.NoError(t, err)

		var found bool

		doc.Definitions.ForEach(func(d ast.Definition, i int) {
			assert.Equal(t, 3, d.Location.Line)
			assert.Equal(t, 4, d.Location.Column)

			found = true
		})

		assert.True(t, found)
	})
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
		source := vast.Source{
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
		source := vast.Source{
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
