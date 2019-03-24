package rules_test

import (
	"fmt"
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func mockNonExecutableDefinitionError(name string, line, col int) types.Error {
	return types.NewError(
		mockNonExecutableDefinitionMessage(name),
		// TODO: Location.
	)
}

func mockNonExecutableDefinitionMessage(name string) string {
	return fmt.Sprintf("The \"%s\" definition is not executable.", name)
}

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
				Add(mockNonExecutableDefinitionError("Cow", 0, 0)).
				Add(mockNonExecutableDefinitionError("Dog", 0, 0)),
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
				Add(mockNonExecutableDefinitionError("schema", 0, 0)).
				Add(mockNonExecutableDefinitionError("Query", 0, 0)).
				Add(mockNonExecutableDefinitionError("schema", 0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.ExecutableDefinitions)
}
