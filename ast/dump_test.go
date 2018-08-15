package ast_test

import (
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/stretchr/testify/assert"
)

var (

	wipTest = strings.TrimSpace(`
query Var($v: Int! = $var) {
  selection
}

query Vars($v: Int! = $var, $i: Int! = 123, $f: Float! = 1.23e+10, $s: String! = "string") {
  selection
}

query Vars2($b: Boolean! = true, $b2: Boolean! = false, $n: Int = null, $e: Enum = ENUM_VALUE) {
  selection
}

query Vars3($l: [Int!]! = [1, 2, 3], $o: Point2D = { x: 1.2, y: 3.4 }) {
  selection
}

query Directive @foo {
  selection
}

query Directives @foo(bar: $baz, baz: "qux") @bar @baz(foo: 123) {
  selection
}

query Selection {
  selection
}

query Selections {
  selection1
  selection2
  selection3 @foo
  selection4 @bar(baz: "qux")
  selection5 @baz(qux: 123) @foo @bar {
    nested {
      aliased: selections
    }
  }
}

query Frags {
  ...userFields
  ... on User @exclude(not: $noInfo) {
    pals: friends {
      count
    }
  }
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
  sendEmail2(message: "Hello, World!\tYours, \u0080 \u754c 😀 \" \\ \/ \b \f \t GraphQL.")
  intVal(foo: 12345)
  floatVal(bar: 123.456)
  floatVal2(bar: 1.23456e+10)
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

fragment postFields on Subscription @skip(if: $bar, do: $not) {
  title
  subtitle
  slug
  category {
    name
    slug
  }
  content
}
`)
	wipTest2 = strings.TrimSpace(`
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
)

func TestSdump(t *testing.T) {
	tt := []struct {
		descr string
		query string
	}{
		{descr: "wip test", query: wipTest},
		//{descr: "wip test2", query: wipTest2},
	}

	for _, tc := range tt {
		psr := parser.New([]byte(tc.query))

		doc, err := psr.Parse()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.query, ast.Sdump(doc))
	}
}
