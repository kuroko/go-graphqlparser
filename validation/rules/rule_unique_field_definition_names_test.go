package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

func TestUniqueFieldDefinitionNames(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no fields",
			query: `
				type SomeObject
				interface SomeInterface
				input SomeInputObject
			`,
		},
		{
			msg: "one field",
			query: `
				type SomeObject {
					foo: String
				}

				interface SomeInterface {
					foo: String
				}

				input SomeInputObject {
					foo: String
				}
			`,
		},
		{
			msg: "multiple fields",
			query: `
				type SomeObject {
					foo: String
					bar: String
				}

				interface SomeInterface {
					foo: String
					bar: String
				}

				input SomeInputObject {
					foo: String
					bar: String
				}
			`,
		},
		{
			msg: "duplicate fields inside the same type definition",
			query: `
				type SomeObject {
					foo: String
					bar: String
					foo: String
				}

				interface SomeInterface {
					foo: String
					bar: String
					foo: String
				}

				input SomeInputObject {
					foo: String
					bar: String
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(duplicateFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "extend type with new field",
			query: `
				type SomeObject {
					foo: String
				}
				extend type SomeObject {
					bar: String
				}
				extend type SomeObject {
					baz: String
				}

				interface SomeInterface {
					foo: String
				}
				extend interface SomeInterface {
					bar: String
				}
				extend interface SomeInterface {
					baz: String
				}

				input SomeInputObject {
					foo: String
				}
				extend input SomeInputObject {
					bar: String
				}
				extend input SomeInputObject {
					baz: String
				}
			`,
		},
		{
			msg: "extend type with duplicate field",
			query: `
				extend type SomeObject {
					foo: String
				}
				type SomeObject {
					foo: String
				}

				extend interface SomeInterface {
					foo: String
				}
				interface SomeInterface {
					foo: String
				}

				extend input SomeInputObject {
					foo: String
				}
				input SomeInputObject {
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(duplicateFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "duplicate field inside extension",
			query: `
				type SomeObject
				extend type SomeObject {
					foo: String
					bar: String
					foo: String
				}

				interface SomeInterface
				extend interface SomeInterface {
					foo: String
					bar: String
					foo: String
				}

				input SomeInputObject
				extend input SomeInputObject {
					foo: String
					bar: String
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(duplicateFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "duplicate field inside different extensions",
			query: `
				type SomeObject
				extend type SomeObject {
					foo: String
				}
				extend type SomeObject {
					foo: String
				}

				interface SomeInterface
				extend interface SomeInterface {
					foo: String
				}
				extend interface SomeInterface {
					foo: String
				}

				input SomeInputObject
				extend input SomeInputObject {
					foo: String
				}
				extend input SomeInputObject {
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(duplicateFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "adding new field to the type inside existing schema",
			schema: &validation.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeObject": {
						Name: "SomeObject",
						Kind: ast.TypeDefinitionKindObject,
					},
					"SomeInterface": {
						Name: "SomeInterface",
						Kind: ast.TypeDefinitionKindInterface,
					},
					"SomeInputObject": {
						Name: "SomeInputObject",
						Kind: ast.TypeDefinitionKindInputObject,
					},
				},
			},
			query: `
				extend type SomeObject {
					foo: String
				}

				extend interface SomeInterface {
					foo: String
				}

				extend input SomeInputObject {
					foo: String
				}
			`,
		},
		{
			msg: "adding conflicting fields to existing schema twice",
			schema: &validation.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeObject": {
						Name: "SomeObject",
						Kind: ast.TypeDefinitionKindObject,
						FieldsDefinition: (*ast.FieldDefinitions)(nil).
							Add(ast.FieldDefinition{
								Name: "foo",
								Type: ast.Type{
									NamedType: "String",
									Kind:      ast.TypeKindNamed,
								},
							}),
					},
					"SomeInterface": {
						Name: "SomeInterface",
						Kind: ast.TypeDefinitionKindInterface,
						FieldsDefinition: (*ast.FieldDefinitions)(nil).
							Add(ast.FieldDefinition{
								Name: "foo",
								Type: ast.Type{
									NamedType: "String",
									Kind:      ast.TypeKindNamed,
								},
							}),
					},
					"SomeInputObject": {
						Name: "SomeInputObject",
						Kind: ast.TypeDefinitionKindInputObject,
						InputFieldsDefinition: (*ast.InputValueDefinitions)(nil).
							Add(ast.InputValueDefinition{
								Name: "foo",
								Type: ast.Type{
									NamedType: "String",
									Kind:      ast.TypeKindNamed,
								},
							}),
					},
				},
			},
			query: `
				extend type SomeObject {
					foo: String
				}
				extend interface SomeInterface {
					foo: String
				}
				extend input SomeInputObject {
					foo: String
				}

				extend type SomeObject {
					foo: String
				}
				extend interface SomeInterface {
					foo: String
				}
				extend input SomeInputObject {
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(existedFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(existedFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(existedFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)).
				Add(existedFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(existedFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(existedFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "adding fields to existing schema twice",
			schema: &validation.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeObject": {
						Name: "SomeObject",
						Kind: ast.TypeDefinitionKindObject,
					},
					"SomeInterface": {
						Name: "SomeInterface",
						Kind: ast.TypeDefinitionKindInterface,
					},
					"SomeInputObject": {
						Name: "SomeInputObject",
						Kind: ast.TypeDefinitionKindInputObject,
					},
				},
			},
			query: `
				extend type SomeObject {
					foo: String
				}
				extend type SomeObject {
					foo: String
				}

				extend interface SomeInterface {
					foo: String
				}
				extend interface SomeInterface {
					foo: String
				}

				extend input SomeInputObject {
					foo: String
				}
				extend input SomeInputObject {
					foo: String
				}
			`,
			errs: (*graphql.Errors)(nil).
				Add(duplicateFieldDefinitionNameMessage("SomeObject", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInterface", "foo", 0, 0)).
				Add(duplicateFieldDefinitionNameMessage("SomeInputObject", "foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, uniqueFieldDefinitionNames)
}
