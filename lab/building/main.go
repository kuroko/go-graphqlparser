package main

import (
	"log"

	"github.com/bucketd/go-graphqlparser"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	schema, errs, err := graphqlparser.ParseSDLDoc([]byte(`
		type Query {
			foo: Foo
		}

		interface Bar {
			sup: String
			bar: ID!
		}

		type Foo implements Bar {
			sup: String
			bar: String!
		}
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
