package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
)

func main() {
	runtime.GOMAXPROCS(1)

	query := []byte(`
query first($foo: Boolean = true) {
    # How about some comments too?
    user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
        ...userFields
    }
    # Do comments even slow us down?
    post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
        ...postFields
    }
    # Directives
    fooa: foo @include(if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo) @include(if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo)
    bara: bar @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar)
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

subscription third {
    ...postFields
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

	start := time.Now()

	var doc ast.Document
	var err error

	for i := 0; i < 100; i++ {
		psr := parser.New(query)

		doc, err = psr.Parse()
		if err != nil {
			fmt.Println(err)
		}

		_ = doc
	}

	fmt.Println(time.Since(start))
}
