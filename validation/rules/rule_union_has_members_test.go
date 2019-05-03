package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUnionHasMembers(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "one member",
			query: `
				union Foo = Bar
			`,
		},
		{
			msg: "many members",
			query: `
				union Foo = Bar | Baz
			`,
		},
		{
			msg: "no members",
			query: `
				union Foo
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnionHasNoMembersError("Foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UnionHasMembers)
}
