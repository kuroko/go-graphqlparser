package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueTypeNames(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no types",
			query: `
				directive @test on SCHEMA
			`,
		},
		{
			msg: "one type",
			query: `
				type Foo
			`,
		},
		{
			msg: "many types",
			query: `
				type Foo
				type Bar
				type Baz
			`,
		},
		{
			msg: "type and non-type definitions named the same",
			query: `
				query Foo { __typename }
				fragment Foo on Query { __typename }
				directive @Foo on SCHEMA

				type Foo
			`,
		},
		{
			msg: "types named the same",
			query: `
				type Foo

				scalar Foo
				type Foo
				interface Foo
				union Foo
				enum Foo
				input Foo
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)).
				Add(rules.DuplicateTypeNameError("Foo", 0, 0)),
		},
		{
			msg: "adding new types to existing schema",
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"Foo": {},
				},
			},
			query: `
				type Bar
			`,
		},
		// TODO: We currently don't have directives stored on this type, revisit this test. Maybe we
		// could also do with something like `buildSchema`?
		//{
		//	msg:    "adding new type to existing schema with same-named directive",
		//	schema: &types.Schema{},
		//	query: `
		//		type Foo
		//	`,
		//},
		{
			msg: "adding conflicting types to existing schema",
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"Foo": {},
				},
			},
			query: `
				scalar Foo
				type Foo
				interface Foo
				union Foo
				enum Foo
				input Foo
			`,
			errs: (*types.Errors)(nil).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)).
				Add(rules.ExistedTypeNameError("Foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueTypeNames)
}
