package main

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	query := `
		} foo($bar: Boolean = {foo: [ENUM_VALUE, false, null], bar: 123, baz: $toplel}, $baz: [String!]) {
			# ...
		}
	`

	psr := parser.New([]byte(query))

	doc, err := psr.Parse()
	spew.Dump(doc)
	if err != nil {
		fmt.Println(err)
	}
}
