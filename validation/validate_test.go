package validation_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
	_ "github.com/vektah/gqlparser/validator/rules"
)

func BenchmarkValidateSDL(b *testing.B) {
	// Default query rules.
	walker := validation.NewWalker(rules.SpecifiedSDL)

	var errs *types.Errors

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// When testing Vektah's parser, we can't re-use this, so let's do the same here, even
		// though in our case we can re-use this document.
		doc, err := graphql.Parse(schemaDoc)
		require.NoError(b, err)

		b.StartTimer()

		errs = validation.ValidateSDL(doc, nil, walker)
	}

	b.StopTimer()

	require.Nil(b, errs)
}

func BenchmarkVektahValidateSchema(b *testing.B) {
	var gerr *gqlerror.Error

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		doc, err := parser.ParseSchemas(validator.Prelude, &ast.Source{
			Name:  "test.graphqls",
			Input: string(schemaDoc),
		})

		require.Nil(b, err)

		b.StartTimer()

		_, gerr = validator.ValidateSchemaDocument(doc)
	}

	b.StopTimer()

	require.Nil(b, gerr)
}

func TestValidate(t *testing.T) {
	// TODO: Move into golden test format.
	parser := language.NewParser([]byte(`
		query { hello }
		query { goodbye }
	`))

	doc, err := parser.Parse()
	if err != nil {
		require.NoError(t, err)
	}

	walker := validation.NewWalker(rules.Specified)
	errs := validation.Validate(doc, &types.Schema{}, walker)

	t.Log(spew.Sdump(errs))

	assert.Equal(t, 2, errs.Len())
}
