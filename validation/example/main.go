package main

import (
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
	"github.com/davecgh/go-spew/spew"
)

var query = []byte(`
query Foo {
  bar
  baz
}
`)

func main() {
	// Later on, probably in a request, create a parser, and parse a query.
	parser := language.NewParser(query)

	doc, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	// TODO: Parse schema first.
	schema := &validation.Schema{}

	// This walker instance should be shared for the lifetime of a whole application, not just a
	// single request. It is stateless, but not cheap to create.
	walker := validation.NewWalker(rules.Specified)

	// Validate the result, returning GraphQL errors.
	errs := validation.Validate(doc, schema, walker)

	spew.Dump(errs)
}
