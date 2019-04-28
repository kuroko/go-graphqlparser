package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestExecutableDefinitions(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "with only operation",
			query: `
				query Foo {
					dog {
						name
					}
				}
			`,
			errs: nil,
		},
		{
			msg: "with operation and fragment",
			query: `
				query Foo {
					dog {
						name
						...Frag
					}
				}

				fragment Frag on Dog {
					name
				}
			`,
			errs: nil,
		},
		{
			msg: "with type definition",
			query: `
				query Foo {
					dog {
						name
					}
				}

				type Cow {
					name: String
				}

				extend type Dog {
					color: String
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.NonExecutableDefinitionError("Cow", 0, 0)).
				Add(validation.NonExecutableDefinitionError("Dog", 0, 0)),
		},
		{
			msg: "with schema definition",
			query: `
				schema {
					query: Query
				}

				type Query {
					test: String
				}

				extend schema @directive
			`,
			errs: (*types.Errors)(nil).
				Add(validation.NonExecutableDefinitionError("schema", 0, 0)).
				Add(validation.NonExecutableDefinitionError("Query", 0, 0)).
				Add(validation.NonExecutableDefinitionError("schema", 0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.ExecutableDefinitions)
}
