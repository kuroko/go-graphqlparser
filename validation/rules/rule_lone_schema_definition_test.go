package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
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
		{
			msg: "schema definition in extending schema document",
			query: `
			schema { query: Foo }
			`,
			schema: &validation.Schema{},
			errs: (*graphql.Errors).
				Add(nil, canNotDefineSchemaWithinExtensionError(0, 0)),
		},
	}

	sdlRuleTester(t, tt, loneSchemaDefinition)
}
