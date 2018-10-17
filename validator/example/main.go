package main

import (
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/bucketd/go-graphqlparser/validator"
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
	prsr := parser.New(query)
	doc, err := prsr.Parse()
	if err != nil {
		panic(err)
	}

	errs := validator.Validate(doc)

	spew.Dump(errs)
}
