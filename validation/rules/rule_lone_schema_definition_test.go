package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
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
			errs: (*types.Errors)(nil).
				Add(rules.SchemaDefinitionNotAloneError(0, 0)).
				Add(rules.SchemaDefinitionNotAloneError(0, 0)),
		},
		{
			msg:    "define schema in schema extension",
			schema: &types.Schema{},
			query: `
			schema {
				query: Foo
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.CanNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "redefine schema in schema extension",
			// TODO: Maybe replace with something that builds this from a schema string?
			schema: &types.Schema{
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
			errs: (*types.Errors)(nil).
				Add(rules.CanNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "redefine implicit schema in schema extension",
			// TODO: Maybe replace with something that builds this from a schema string?
			// TODO: This isn't "valid" really, we're not testing the implicit schema definition.
			schema: &types.Schema{
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
			errs: (*types.Errors)(nil).
				Add(rules.CanNotDefineSchemaWithinExtensionError(0, 0)),
		},
		{
			msg: "extend schema in schema extension",
			schema: &types.Schema{
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

	sdlRuleTester(t, tt, rules.LoneSchemaDefinition)
}
