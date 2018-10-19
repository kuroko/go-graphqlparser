package main

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/ast/validation"
	"github.com/bucketd/go-graphqlparser/language"
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

	errs := ast.Validate(doc, validation.Rules)

	spew.Dump(errs)
}
