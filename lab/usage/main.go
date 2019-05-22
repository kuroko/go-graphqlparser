package main

import (
	"github.com/bucketd/go-graphqlparser"
	"github.com/davecgh/go-spew/spew"
)

var doc = []byte(``)

func main() {
	// Parse an SDL document.
	schema, errs, err := graphqlparser.ParseSDLDoc(doc, nil)
	if err != nil {
		panic(err)
	}

	if errs.Len() > 0 {
		// ...
	}

	// Parse a query document.
	queryAST, errs, err := graphqlparser.ParseDoc(doc, schema)
	if err != nil {
		panic(err)
	}

	if errs.Len() > 0 {
		// ...
	}

	spew.Dump(queryAST)
}
