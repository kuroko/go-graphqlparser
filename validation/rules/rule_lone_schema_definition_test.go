package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
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
			msg:    "define schema in schema extension",
			schema: &validation.Schema{},
			query: `
			schema {
				query: Foo
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, canNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "redefine schema in schema extension",
			// TODO: Maybe replace with something that builds this from a schema string?
			schema: &validation.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Foo",
				},
				QueryTypeDefined: true,
			},
			query: `
			schema {
				mutation: Foo
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, canNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "redefine implicit schema in schema extension",
			// TODO: Maybe replace with something that builds this from a schema string?
			// TODO: This isn't "valid" really, we're not testing the implicit schema definition.
			schema: &validation.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Foo",
				},
				QueryTypeDefined: true,
			},
			query: `
			schema {
				mutation: Foo
			}
			`,
			errs: (*graphql.Errors).
				Add(nil, canNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "extend schema in schema extension",
			schema: &validation.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Foo",
				},
				QueryTypeDefined: true,
			},
			query: `
			extend schema {
				mutation: Foo
			}
			`,
			errs: nil,
		},
	}

	sdlRuleTester(t, tt, loneSchemaDefinition)
}
