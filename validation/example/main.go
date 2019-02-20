package main

import (
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/davecgh/go-spew/spew"

	_ "github.com/bucketd/go-graphqlparser/validation/rules"
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

	// Validate the result, returning GraphQL errors.
	errs := validation.ValidateWithWalker(doc, schema, validation.DefaultQueryWalker)

	spew.Dump(errs)
}
