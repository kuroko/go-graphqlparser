package parser

import (
	"testing"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

var (
	normalQuery = []byte(`
query first($foo: Boolean = true, $foo: Boolean = true) {
    # How about some comments too?
    user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
        ...userFields
    }
    # Do comments even slow us down?
    post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
        ...postFields
    }
    # Directives
    fooa: foo @include(if: $foo)
    bara: bar @skip(if: $bar)
    baza: baz @gt(val: $baz)
    # Inline fragments
    ... @include(if: $expandedInfo) {
        firstName
        lastName
        birthday
    }
}

mutation second($variable: String = "test") {
    sendEmail(message: """
        Hello,
            World!

        Yours,
            GraphQL
    """)
    sendEmail2(message: "Hello\n,  World!\n\nYours,\n  GraphQL.")
    intVal(foo: 12345)
    floatVal(bar: 123.456)
    floatVal2(bar: 123.456e10)
    boolVal(bool: false)
    listVal(list: [1, 2, 3])
    variableVal(var: $variable)
}

fragment userFields on User {
    firstName
    lastName
    title
    company {
        name
        slug
    }
    email
    mobile
}

fragment postFields on Subscription {
    title
    subtitle
    slug
    author {
        ...userFields
    }
    category {
        name
        slug
    }
    content
}
	`)

	tinyQuery = []byte(`
		{
			person {
				name
			}
		}
	`)
)

func BenchmarkParser(b *testing.B) {
	tt := []struct {
		name  string
		query []byte
	}{
		{name: "normalQuery", query: normalQuery},
		{name: "tinyQuery", query: tinyQuery},
	}

	for _, t := range tt {
		b.Run(t.name, func(b *testing.B) {
			b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
				runBucketdParser(b, t.query)
			})

			b.Run("github.com/graphql-go/graphql", func(b *testing.B) {
				runGraphQLGoParser(b, t.query)
			})
		})
	}
}

func runBucketdParser(b *testing.B, query []byte) {
	for i := 0; i < b.N; i++ {
		psr := New(query)

		ast, err := psr.Parse()
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}

func runGraphQLGoParser(b *testing.B, query []byte) {
	for i := 0; i < b.N; i++ {
		params := parser.ParseParams{
			Source: source.NewSource(&source.Source{
				Body: query,
			}),
		}

		ast, err := parser.Parse(params)
		if err != nil {
			b.Fatal(err)
		}

		_ = ast
	}
}
