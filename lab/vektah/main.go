package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

func main() {
	schema, err := gqlparser.LoadSchema(&ast.Source{
		Input: `
			type Query {
				bar: String @deprecated(foo: "bar") @deprecated(foo: "bar")
			}
		`,
	})

	if err != nil {
		panic(err)
	}

	spew.Dump(schema)
}
