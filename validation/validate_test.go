package validation_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

var query = []byte(`
query Foo {
  bar
  baz
  qux
}
`)

func BenchmarkValidate(b *testing.B) {
	walker := validation.NewWalker(rules.Specified)
	parser := language.NewParser(query)

	doc, err := parser.Parse()
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validation.Validate(doc, walker)
	}
}
