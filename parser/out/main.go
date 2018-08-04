package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	runtime.GOMAXPROCS(1)

	query := []byte(`
query foo($foo: Boolean = 2) {
	hello @foo(bar: "baz") {
		foo
		bar
	}
	world
}
	`)

	start := time.Now()

	var doc ast.Document
	var err error

	for i := 0; i < 1; i++ {
		psr := parser.New(query)

		doc, err = psr.Parse()
		if err != nil {
			fmt.Println(err)
		}

		_ = doc
	}

	fmt.Println(time.Since(start))

	for _, def := range doc.Definitions {
		selections := def.ExecutableDefinition.SelectionSet
		spew.Dump(selections)

		for {
			fmt.Println(selections.Data.Name)
			if selections.Next == nil {
				break
			}

			selections = selections.Next
		}
	}
}
