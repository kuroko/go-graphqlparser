package rules

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
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
			errs: (*graphql.Errors)(nil).
				Add(duplicateEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
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
			errs: (*graphql.Errors)(nil).
				Add(duplicateEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
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
			errs: (*graphql.Errors)(nil).
				Add(duplicateEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
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
			errs: (*graphql.Errors)(nil).
				Add(duplicateEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "adding new value to the type inside existing schema",
			schema: &validation.Schema{
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
			schema: &validation.Schema{
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
			errs: (*graphql.Errors)(nil).
				Add(existedEnumValueNameMessage("SomeEnum", "FOO", 0, 0)).
				Add(existedEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
		},
		{
			msg: "adding conflicting value to existing schema twice",
			schema: &validation.Schema{
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
			errs: (*graphql.Errors)(nil).
				Add(duplicateEnumValueNameMessage("SomeEnum", "FOO", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, uniqueEnumValueNames)
}
