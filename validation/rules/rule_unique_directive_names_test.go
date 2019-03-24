package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueDirectiveNames(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no directive",
			query: `
				type Foo
			`,
		},
		{
			msg: "one directive",
			query: `
				directive @foo on SCHEMA
			`,
		},
		{
			msg: "many directives",
			query: `
				directive @foo on SCHEMA
				directive @bar on SCHEMA
				directive @baz on SCHEMA
			`,
		},
		{
			msg: "directive and non-directive definitions named the same",
			query: `
				query foo { __typename }
				fragment foo on foo { __typename }
				type foo

				directive @foo on SCHEMA
			`,
		},
		{
			msg: "directives named the same",
			query: `
				directive @foo on SCHEMA
				directive @foo on SCHEMA
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateDirectiveNameError("foo", 0, 0)),
		},
		{
			msg: "adding new directive to existing schema",
			schema: &types.Schema{
				Directives: map[string]*ast.DirectiveDefinition{
					"foo": {
						Name:               "foo",
						DirectiveLocations: ast.DirectiveLocationKindSchema,
					},
				},
			},
			query: `
				directive @bar on SCHEMA
			`,
		},
		{
			msg: "adding new directive with standard name to existing schema",
			// TODO: If we implement `buildSchema`, we should be adding `skip`, `include`, and
			// `deprecated` to the directives map.
			schema: &types.Schema{
				Directives: map[string]*ast.DirectiveDefinition{
					"skip": {
						Name: "skip",
						DirectiveLocations: ast.DirectiveLocationKindField |
							ast.DirectiveLocationKindFragmentSpread |
							ast.DirectiveLocationKindInlineFragment,
					},
				},
			},
			query: `
				directive @skip on SCHEMA
			`,
			errs: (*types.Errors)(nil).
				Add(rules.ExistedDirectiveNameError("skip", 0, 0)),
		},
		{
			msg: "adding new directive to existing schema with same-named type",
			// TODO: If we implement `buildSchema`, we should be adding `skip`, `include`, and
			// `deprecated` to the directives map.
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"foo": {
						Name: "foo",
					},
				},
			},
			query: `
				directive @foo on SCHEMA
			`,
		},
		{
			msg: "adding conflicting directives to existing schema",
			// TODO: If we implement `buildSchema`, we should be adding `skip`, `include`, and
			// `deprecated` to the directives map.
			schema: &types.Schema{
				Directives: map[string]*ast.DirectiveDefinition{
					"foo": {
						Name:               "foo",
						DirectiveLocations: ast.DirectiveLocationKindSchema,
					},
				},
			},
			query: `
				directive @foo on SCHEMA
			`,
			errs: (*types.Errors)(nil).
				Add(rules.ExistedDirectiveNameError("foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueDirectiveNames)
}
