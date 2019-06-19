package main

import (
	"log"

	"github.com/bucketd/go-graphqlparser"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	schema, errs, err := graphqlparser.ParseSDLDoc([]byte(`
		schema {
			query: Foo
		}

		type Foo {
			bar: String!
		}

		extend type Foo {
			baz: Int!
		}

		directive @f9oo on QUERY
	`), nil)

	if err != nil {
		panic(err)
	}

	if errs.Len() > 0 {
		log.Println("Error(s) encountered whilst parsing SDL document:")

		errs.ForEach(func(e graphql.Error, i int) {
			log.Printf("- %s", e.Message)
		})

		log.Fatalln("Exiting...")
	}

	spew.Dump(schema)
}
