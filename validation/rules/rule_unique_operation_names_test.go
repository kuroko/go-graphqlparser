// +build ignore

package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueOperationNames(t *testing.T) {
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
			msg: "one named operation",
			query: `
			query Foo {
			  field
			}
			`,
		},
		{
			msg: "multiple operations",
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
			msg: "multiple operations of different types",
			query: `
			query Foo {
			  field
			}
			mutation Bar {
			  field
			}
			subscription Baz {
			  field
			}
			`,
		},
		{
			msg: "fragment and operation named the same",
			query: `
			query Foo {
			  ...Foo
			}
			fragment Foo on Type {
			  field
			}
			`,
		},
		{
			msg: "multiple operations of same name",
			query: `
			query Foo {
			  fieldA
			}
			query Foo {
			  fieldB
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateOperationNameError("Foo", 0, 0)),
		},
		{
			msg: "multiple ops of same name of different types (mutation)",
			query: `
			query Foo {
			  fieldA
			}
			mutation Foo {
			  fieldB
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateOperationNameError("Foo", 0, 0)),
		},
		{
			msg: "multiple ops of same name of different types (subscription)",
			query: `
			query Foo {
			  fieldA
			}
			subscription Foo {
			  fieldB
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateOperationNameError("Foo", 0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.UniqueOperationNames)
}
