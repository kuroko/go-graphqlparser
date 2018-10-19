package main

import (
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
	"github.com/davecgh/go-spew/spew"
)

var query = []byte(`
query Foo {
  selectionA
  selectionB
}

enum Bar {
  BAZ
  QUX
}
`)

func main() {
	parser := language.NewParser(query)
	doc, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	errs := validation.Validate(doc, rules.Specified)

	spew.Dump(errs)
}
