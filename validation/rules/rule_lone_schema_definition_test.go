package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
)

func TestLoneSchemaDefinition(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no schema",
			query: `
			type Foo { foo: String }
			`,
			errs: nil,
		},
		{
			msg: "one schema definition",
			query: `
			schema { query: Foo }
			type Foo { foo: String }
			`,
			errs: nil,
		},
		{
			msg: "multiple schema definitions",
			query: `
			schema { query: Foo }
			type Foo { foo: String }
			schema { mutation: Foo }
			schema { subscription: Foo }
			`,
			errs: (*graphql.Errors).
				Add(nil, schemaDefinitionNotAloneError(0, 0)).
				Add(schemaDefinitionNotAloneError(0, 0)),
		},
	}

	ruleTester(t, tt, loneSchemaDefinition)
}
