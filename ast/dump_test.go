package ast_test

import (
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/stretchr/testify/assert"
)

var (
	shorthandTest = strings.TrimSpace(`
{
  hello
  world
}
`)
	exhaustiveTest = strings.TrimSpace(`
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
  ...userFields @testytesty
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
      World! \"""
    
    \u0080 \u754c 😀 \"
    
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

query first($foo: Boolean = true, $foo: Boolean = true) {
  user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
    ...userFields
  }
  post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
    ...postFields
  }
  fooa: foo @include(if: $foo)
  bara: bar @skip(if: $bar)
  baza: baz @gt(val: $baz)
  ... @include(if: $expandedInfo) {
    firstName
    lastName
    birthday
  }
}

schema @foo @bar {
  query: Query
  mutation: Mutation
}

"""
Documentation blocks should support
...multiple lines

And gaps in lines
"""
directive @foo on QUERY

scalar JSON

type Person implements NamedEntity {
  name: String
  weight(equipmentWeight: Int, ropeStrength: Int): Int
}

type Business implements NamedEntity & ValuedEntity {
  name: String
  value: Int
  employeeCount: Int
}

type Climbing implements
    & Ropes
    & Rocks
    & Chalk {
  name: String
  weight: Int
}

interface NamedEntity {
  name: String
}

interface ValuedEntity {
  value: Int
}

union SmallestResult = DriftKing

union SmallerResult = Donkey | Kong

union BiggerResult =
  | Orange
  | Black
  | RX7

enum Direction {
  NORTH
  EAST
  SOUTH
  WEST
}

input Point2D {
  x: Float
  y: Float
}

directive @example on FIELD

directive @example on FIELD_DEFINITION | ARGUMENT_DEFINITION

directive @example on
  | FIELD
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT

extend schema @tonight

extend schema {
  query: TheBrain
  mutation: Pinky
}

extend scalar IntyThing @add

extend type Car implements Transport {
"""
Should the quotes around this be indented the same as the fields below?
Indentation is not preserved.
"""
  upgrade(component: CarComponent): Cost
}

extend interface NamedEntity @addedDirective

extend interface Transport {
  cost: Int
}

extend union SmallestResult = Nissan350z

extend enum Direction {
  NORTHWEST
  NORTHEAST
}

extend input Point2D {
  distanceFromOrigin: Length
}
`)
)

func TestSdump(t *testing.T) {
	tt := []struct {
		descr string
		query string
	}{
		{descr: "shorthand test", query: shorthandTest},
		{descr: "exhaustive test", query: exhaustiveTest},
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
