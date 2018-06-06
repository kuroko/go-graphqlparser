package main

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	query := `query foo { bar baz }`
	psr := parser.New([]byte(query))

	doc, err := psr.Parse()
	spew.Dump(doc)
	if err != nil {
		fmt.Println(err)
	}
}
