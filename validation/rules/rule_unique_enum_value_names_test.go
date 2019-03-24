package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueEnumValueNames(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no values",
			query: `
				enum SomeEnum
			`,
		},
		{
			msg: "one value",
			query: `
				enum SomeEnum {
					FOO
				}
			`,
		},
		{
			msg: "multiple values",
			query: `
				enum SomeEnum {
					FOO
					BAR
				}
			`,
		},
		{
			msg: "duplicate values inside the same enum definition",
			query: `
				enum SomeEnum {
					FOO
					BAR
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "extend enum with new value",
			query: `
				enum SomeEnum {
					FOO
				}

				extend enum SomeEnum {
					BAR
				}

				extend enum SomeEnum {
					BAZ
				}
			`,
		},
		{
			msg: "extend enum with duplicate value",
			query: `
				enum SomeEnum {
					FOO
				}

				extend enum SomeEnum {
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "duplicate value inside extension",
			query: `
				enum SomeEnum

				extend enum SomeEnum {
					FOO
					BAR
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "duplicate value inside different extensions",
			query: `
				enum SomeEnum

				extend enum SomeEnum {
					FOO
				}

				extend enum SomeEnum {
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "adding new value to the type inside existing schema",
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeEnum": {
						Name: "SomeEnum",
						Kind: ast.TypeDefinitionKindEnum,
					},
				},
			},
			query: `
				extend enum SomeEnum {
					FOO
				}
			`,
		},
		{
			msg: "adding conflicting value to existing schema twice",
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeEnum": {
						Name: "SomeEnum",
						Kind: ast.TypeDefinitionKindEnum,
						EnumValuesDefinition: (*ast.EnumValueDefinitions)(nil).
							Add(ast.EnumValueDefinition{
								EnumValue: "FOO",
							}),
					},
				},
			},
			query: `
				extend enum SomeEnum {
					FOO
				}

				extend enum SomeEnum {
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.ExistedEnumValueNameError("SomeEnum", "FOO", 0, 0)).
				Add(rules.ExistedEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "adding conflicting value to existing schema twice",
			schema: &types.Schema{
				Types: map[string]*ast.TypeDefinition{
					"SomeEnum": {
						Name: "SomeEnum",
						Kind: ast.TypeDefinitionKindEnum,
					},
				},
			},
			query: `
				extend enum SomeEnum {
					FOO
				}

				extend enum SomeEnum {
					FOO
				}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateEnumValueNameError("SomeEnum", "FOO", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueEnumValueNames)
}
