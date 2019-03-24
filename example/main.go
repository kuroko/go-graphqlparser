package main

import (
	"fmt"
	"os"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/davecgh/go-spew/spew"
)

var rawSchema = []byte(`
type Query {
  bar: String
  baz: String
}
`)

var rawQuery = []byte(`
query Foo {
  bar
  baz
}
`)

func main() {
	schema, errs, err := graphql.BuildSchema(nil, rawSchema)
	if err != nil {
		panic(err)
	}

	if errs.Len() > 0 {
		fmt.Println("Failed to validate schema")
		errs.ForEach(func(e types.Error, i int) {
			fmt.Println(e.Message)
		})

		os.Exit(1)
	}

	errs, err = graphql.Validate(schema, rawQuery)
	if err != nil {
		panic(err)
	}

	spew.Dump(errs)
}
