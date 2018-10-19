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
	// A walker with a given set of rules can be re-used as many times as is necessary.
	walker := validation.NewWalker(rules.Specified)

	// Later on, probably in a request, create a parser, and parse a query.
	parser := language.NewParser(query)
	doc, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	// Validate the result, returning GraphQL errors.
	errs := validation.Validate(doc, walker)

	spew.Dump(errs)
}
