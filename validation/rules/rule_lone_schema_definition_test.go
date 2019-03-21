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
			type Foo { checkEnumValueUniqueness: String }
			`,
			errs: nil,
		},
		{
			msg: "one schema definition",
			query: `
			schema { query: Foo }
			type Foo { checkEnumValueUniqueness: String }
			`,
			errs: nil,
		},
		{
			msg: "multiple schema definitions",
			query: `
			schema { query: Foo }
			type Foo { checkEnumValueUniqueness: String }
			schema { mutation: Foo }
			schema { subscription: Foo }
			`,
			errs: (*graphql.Errors)(nil).
				Add(schemaDefinitionNotAloneError(0, 0)).
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
			errs: (*graphql.Errors)(nil).
				Add(canNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "redefine schema in schema extension",
			// TODO: Maybe replace with something that builds this from a schema string?
			schema: &validation.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Foo",
				},
			},
			query: `
			schema {
				mutation: Foo
			}
			`,
			errs: (*graphql.Errors)(nil).
				Add(canNotDefineSchemaWithinExtensionError(0, 0)),
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
			},
			query: `
			schema {
				mutation: Foo
			}
			`,
			errs: (*graphql.Errors)(nil).
				Add(canNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "extend schema in schema extension",
			schema: &validation.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Foo",
				},
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
