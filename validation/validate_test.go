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
)

var query = []byte(`
query Foo {
  bar
  baz
  qux
}
`)

func BenchmarkValidate(b *testing.B) {
	parser := language.NewParser(query)

	doc, err := parser.Parse()
	if err != nil {
		b.Error(err)
	}

	walker := validation.NewWalker(rules.Specified)

	schema := graphql.MustBuildSchema(nil, []byte("type Foo"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validation.Validate(doc, schema, walker)
	}
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
