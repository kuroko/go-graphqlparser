package validation_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
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

	b.ResetTimer()

	walker := validation.NewWalker(rules.Specified)

	for i := 0; i < b.N; i++ {
		validation.Validate(doc, walker)
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
	errs := validation.Validate(doc, walker)

	assert.Equal(t, 4, errs.Len())
}
