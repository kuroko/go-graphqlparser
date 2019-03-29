package graphql_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
)

// BenchmarkValidate should give a good idea of how quick our library is overall for all
// pre-execution phases of working with GraphQL documents.
func BenchmarkValidate(b *testing.B) {
	schema := graphql.MustBuildSchema(nil, schemaDoc)

	b.ReportAllocs()
	b.ResetTimer()

	var gerrs *types.Errors
	var err error

	for i := 0; i < b.N; i++ {
		gerrs, err = graphql.Validate(schema, queryDoc)
	}

	b.StopTimer()

	require.NoError(b, err)
	require.Nil(b, gerrs)
}

func BenchmarkValidateSDL(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var gerrs *types.Errors
	var err error

	for i := 0; i < b.N; i++ {
		gerrs, err = graphql.ValidateSDL(nil, schemaDoc)
	}

	b.StopTimer()

	require.NoError(b, err)
	require.Nil(b, gerrs)
}

// BenchmarkBuildSchema tests how quickly we can prepare a schema for use in extension, validation
// and execution of other GraphQL documents.
func BenchmarkBuildSchema(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var gerrs *types.Errors
	var err error

	for i := 0; i < b.N; i++ {
		_, gerrs, err = graphql.BuildSchema(nil, schemaDoc)
	}

	b.StopTimer()

	require.NoError(b, err)
	require.Nil(b, gerrs)
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

	var gerr gqlerror.List

	input := string(queryDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, gerr = gqlparser.LoadQuery(schema, input)
	}

	b.StopTimer()

	require.Nil(b, gerr)
}

// BenchmarkVektahLoadSchema tests how quickly Vektah's gqlparser library can produce it's
// SDL-related output. This probably matches closest with our ValidateSDL function, but might be
// more similar to our BuildSchema function depending on what extra information is found in Vektah's
// AST. We use our Schema table as a symbol table.
func BenchmarkVektahLoadSchema(b *testing.B) {
	var gerr *gqlerror.Error

	input := string(schemaDoc)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, gerr = gqlparser.LoadSchema(&ast.Source{
			Name:  "test.grahpqls",
			Input: input,
		})
	}

	b.StopTimer()

	require.Nil(b, gerr)
}
