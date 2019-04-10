package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
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
				Add(validation.DuplicateDirectiveNameError("foo", 0, 0)),
		},
		{
			msg: "adding new directive to existing schema",
			schema: graphql.MustBuildSchema(nil, []byte(`
				directive @foo on SCHEMA
			`)),
			query: `
				directive @bar on SCHEMA
			`,
		},
		{
			msg: "adding new directive with standard name to existing schema",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type foo
			`)),
			query: `
				directive @skip on SCHEMA
			`,
			errs: (*types.Errors)(nil).
				Add(validation.ExistedDirectiveNameError("skip", 0, 0)),
		},
		{
			msg: "adding new directive to existing schema with same-named type",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type foo
			`)),
			query: `
				directive @foo on SCHEMA
			`,
		},
		{
			msg: "adding conflicting directives to existing schema",
			schema: graphql.MustBuildSchema(nil, []byte(`
				directive @foo on SCHEMA
			`)),
			query: `
				directive @foo on SCHEMA
			`,
			errs: (*types.Errors)(nil).
				Add(validation.ExistedDirectiveNameError("foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, func(w *validation.Walker) {})
}
