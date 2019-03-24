package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestLoneAnonymousOperation(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no operations",
			query: `
				fragment fragA on Type {
					field
				}
			`,
		},
		{
			msg: "one anon operation",
			query: `
				{
					field
				}
			`,
		},
		{
			msg: "multiple named operations",
			query: `
				query Foo {
					field
				}

				query Bar {
					field
				}
			`,
		},
		{
			msg: "anon operation with fragment",
			query: `
				{
					...Foo
				}

				fragment Foo on Type {
					field
				}
			`,
		},
		{
			msg: "multiple anon operations",
			query: `
				{
					fieldA
				}

				{
					fieldB
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.AnonOperationNotAloneError(0, 0)).
				Add(rules.AnonOperationNotAloneError(0, 0)),
		},
		{
			msg: "anon operation with a mutation",
			query: `
				{
					fieldA
				}

				mutation Foo {
					fieldB
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.AnonOperationNotAloneError(0, 0)),
		},
		{
			msg: "anon operation with a subscription",
			query: `
				{
					fieldA
				}

				subscription Foo {
					fieldB
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.AnonOperationNotAloneError(0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.LoneAnonymousOperation)
}
