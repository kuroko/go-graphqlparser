package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	runtime.GOMAXPROCS(1)

	query := []byte(`
query {
	foo(content: """
		Hello,
	
			Welcome to GraphQL. \""" \t
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
		From, Bucketd
	""")
}
	`)

	query = []byte(`
schema @foo @bar {
	query: Query
}

directive @foo on SCHEMA | SCALAR
directive @bar on SCHEMA | QUERY

directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT

directive @example on
  | FIELD
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT
`)

	start := time.Now()

	var doc ast.Document
	var err error

	for i := 0; i < 1; i++ {
		psr := language.NewParser(query)

		doc, err = psr.Parse()
		if err != nil {
			fmt.Println(err)
		}

		_ = doc
	}

	fmt.Println(time.Since(start))

	spew.Dump(doc)

	//doc.Definitions.ForEach(func(definition ast.Definition, _ int) {
	//	definition.ExecutableDefinition.SelectionSet.ForEach(func(selection ast.Selection, _ int) {
	//		selection.Arguments.ForEach(func(argument ast.Argument, _ int) {
	//			fmt.Println(argument.Value.StringValue)
	//		})
	//	})
	//})
}
