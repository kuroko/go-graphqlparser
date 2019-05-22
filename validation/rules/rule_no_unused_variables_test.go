package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"

	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestNoUnusedVariables(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "uses all variables",
			query: `
				query ($a: String, $b: String, $c: String) {
					field(a: $a, b: $b, c: $c)
				}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables deeply",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					field(a: $a) {
						field(b: $b) {
							field(c: $c)
						}
					}
				}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables deeply in inline fragments",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					... on Type {
						field(a: $a) {
							field(b: $b) {
								... on Type {
									field(c: $c)
								}
							}
						}
					}
				}
			`,
			errs: nil,
		},
		{
			msg: "uses all variables in fragments",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					...FragA
				}
				fragment FragA on Type {
					field(a: $a) {
						...FragB
					}
				}
				fragment FragB on Type {
					field(b: $b) {
						...FragC
					}
				}
				fragment FragC on Type {
					field(c: $c)
				}
			`,
			errs: nil,
		},
		{
			msg: "variable used by fragment in multiple operations",
			query: `
				query Foo($a: String) {
					...FragA
				}
				query Bar($b: String) {
					...FragB
				}
				fragment FragA on Type {
					field(a: $a)
				}
				fragment FragB on Type {
					field(b: $b)
				}
			`,
			errs: nil,
		},
		{
			msg: "variable used by recursive fragment",
			query: `
				query Foo($a: String) {
					...FragA
				}
				fragment FragA on Type {
					field(a: $a) {
						...FragA
					}
				}
			`,
			errs: nil,
		},
		{
			msg: "variable not used",
			query: `
				query ($a: String, $b: String, $c: String) {
					field(a: $a, b: $b)
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("c", "", 0, 0)),
		},
		{
			msg: "multiple variables not used",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					field(b: $b)
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("a", "Foo", 0, 0)).
				Add(validation.UnusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "variable not used in fragments",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					...FragA
				}
				fragment FragA on Type {
					field(a: $a) {
						...FragB
					}
				}
				fragment FragB on Type {
					field(b: $b) {
						...FragC
					}
				}
				fragment FragC on Type {
					field
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "multiple variables not used in fragments",
			query: `
				query Foo($a: String, $b: String, $c: String) {
					...FragA
				}
				fragment FragA on Type {
					field {
						...FragB
					}
				}
				fragment FragB on Type {
					field(b: $b) {
						...FragC
					}
				}
				fragment FragC on Type {
					field
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("a", "Foo", 0, 0)).
				Add(validation.UnusedVariableError("c", "Foo", 0, 0)),
		},
		{
			msg: "variable not used by unreferenced fragment",
			query: `
				query Foo($b: String) {
					...FragA
				}
				fragment FragA on Type {
					field(a: $a)
				}
				fragment FragB on Type {
					field(b: $b)
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("b", "Foo", 0, 0)),
		},
		{
			msg: "variable not used by fragment used by other operation",
			query: `
				query Foo($b: String) {
					...FragA
				}
				query Bar($a: String) {
					...FragB
				}
				fragment FragA on Type {
					field(a: $a)
				}
				fragment FragB on Type {
					field(b: $b)
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.UnusedVariableError("b", "Foo", 0, 0)).
				Add(validation.UnusedVariableError("a", "Bar", 0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.NoUnusedVariables)
}
