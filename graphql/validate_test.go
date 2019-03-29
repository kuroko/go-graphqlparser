package graphql_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

// BenchmarkValidate should give a good idea of how quick our library is overall for all
// pre-execution phases of working with GraphQL documents.
func BenchmarkValidate(b *testing.B) {
	schema := graphql.MustBuildSchema(nil, schemaDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		errs, err := graphql.Validate(schema, queryDoc)
		_, _ = errs, err
	}
}

func BenchmarkValidateSDL(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		errs, err := graphql.ValidateSDL(nil, schemaDoc)
		_, _ = errs, err
	}
}

// BenchmarkBuildSchema tests how quickly we can prepare a schema for use in extension, validation
// and execution of other GraphQL documents.
func BenchmarkBuildSchema(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		schema := graphql.MustBuildSchema(nil, schemaDoc)
		_ = schema
	}
}

// BenchmarkVektahLoadSchema tests how quickly Vektah's gqlparser library can produce it's
// SDL-related output. This probably matches closest with our ValidateSDL function, but might be
// more similar to our BuildSchema function depending on what extra information is found in Vektah's
// AST. We use our Schema table as a symbol table.
func BenchmarkVektahLoadQuery(b *testing.B) {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "test.grahpqls",
		Input: string(schemaDoc),
	})

	require.Nil(b, err)

	input := string(queryDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		schema, err := gqlparser.LoadQuery(schema, input)
		_, _ = schema, err
	}
}

// BenchmarkVektahLoadSchema tests how quickly Vektah's gqlparser library can produce it's
// SDL-related output. This probably matches closest with our ValidateSDL function, but might be
// more similar to our BuildSchema function depending on what extra information is found in Vektah's
// AST. We use our Schema table as a symbol table.
func BenchmarkVektahLoadSchema(b *testing.B) {
	input := string(schemaDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		schema, err := gqlparser.LoadSchema(&ast.Source{
			Name:  "test.grahpqls",
			Input: input,
		})

		_, _ = schema, err
	}
}
