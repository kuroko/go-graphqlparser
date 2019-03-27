package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestKnownTypeNames(t *testing.T) {
	t.Run("query document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "known type names are valid",
				query: `
					query Foo($var: String, $required: [String!]!) {
						user(id: 4) {
							pets { ... on Pet { name }, ...PetFields, ... { name } }
						}
					}

					fragment PetFields on Pet {
						name
					}
				`,
			},
			{
				msg: "unknown type names are invalid",
				query: `
					query Foo($var: JumbledUpLetters) {
						user(id: 4) {
							name
							pets { ... on Badger { name }, ...PetFields }
						}
					}

					fragment PetFields on Peettt {
						name
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownTypeError("JumbledUpLetters", []string{}, 0, 0)).
					Add(rules.UnknownTypeError("Badger", []string{}, 0, 0)).
					Add(rules.UnknownTypeError("Peettt", []string{"Pet"}, 0, 0)),
			},
			// TODO: It's not possible to use our parser and have a schema without the built-in
			// scalar types included. It's part of the spec, and a server would be pretty useless
			// without them, so we're not going to include this test, but for reference, it's name
			// is included below:
			//{
			//	msg: "references to standard scalars that are missing in schema",
			//},
		}

		queryRuleTester(t, tt, rules.KnownTypeNames)
	})

	t.Run("sdl document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "use standard scalars",
				query: `
					type Query {
						string: String
						int: Int
						float: Float
						boolean: Boolean
						id: ID
					}
				`,
			},
			{
				msg: "reference types defined inside the same document",
				query: `
					union SomeUnion = SomeObject | AnotherObject

					type SomeObject implements SomeInterface {
						someScalar(arg: SomeInputObject): SomeScalar
					}

					type AnotherObject {
						foo(arg: SomeInputObject): String
					}

					type SomeInterface {
						someScalar(arg: SomeInputObject): SomeScalar
					}

					input SomeInputObject {
						someScalar: SomeScalar
					}

					scalar SomeScalar

					type RootQuery {
						someInterface: SomeInterface
						someUnion: SomeUnion
						someScalar: SomeScalar
						someObject: SomeObject
					}

					schema {
						query: RootQuery
					}
				`,
			},
			{
				msg: "unknown type references",
				query: `
					type A
					type B

					type SomeObject implements C {
						e(d: D): E
					}

					union SomeUnion = F | G

					interface SomeInterface {
						i(h: H): I
					}

					input SomeInput {
						j: J
					}

					directive @SomeDirective(k: K) on QUERY

					schema {
						query: L
						mutation: M
						subscription: N
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownTypeError("C", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("D", []string{"ID", "A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("E", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("F", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("G", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("H", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("I", []string{"ID", "A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("J", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("K", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("L", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("M", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("N", []string{"A", "B"}, 0, 0)),
			},
			{
				msg: "doesn't consider non-type definitions",
				query: `
					query Foo { __typename }
					fragment Foo on Query { __typename }
					directive @Foo on QUERY

					type Query {
						foo: Foo
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownTypeError("Foo", []string{}, 0, 0)),
			},
			{
				msg:    "reference standard scalars inside extension document",
				schema: graphql.MustBuildSchema(nil, []byte(`type Foo`)),
				query: `
					type SomeType {
						string: String
						int: Int
						float: Float
						boolean: Boolean
						id: ID
					}
				`,
			},
			{
				msg:    "reference types inside extension document",
				schema: graphql.MustBuildSchema(nil, []byte(`type Foo`)),
				query: `
					type QueryRoot {
						foo: Foo
						bar: Bar
					}

					scalar Bar

					schema {
						query: QueryRoot
					}
				`,
			},
			{
				msg:    "unknown type references inside extension document",
				schema: graphql.MustBuildSchema(nil, []byte(`type A`)),
				query: `
					type B

					type SomeObject implements C {
						e(d: D): E
					}

					union SomeUnion = F | G

					interface SomeInterface {
						i(h: H): I
					}

					input SomeInput {
						j: J
					}

					directive @SomeDirective(k: K) on QUERY

					schema {
						query: L
						mutation: M
						subscription: N
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.UnknownTypeError("C", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("D", []string{"ID", "A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("E", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("F", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("G", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("H", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("I", []string{"ID", "A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("J", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("K", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("L", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("M", []string{"A", "B"}, 0, 0)).
					Add(rules.UnknownTypeError("N", []string{"A", "B"}, 0, 0)),
			},
		}

		sdlRuleTester(t, tt, rules.KnownTypeNames)
	})
}
