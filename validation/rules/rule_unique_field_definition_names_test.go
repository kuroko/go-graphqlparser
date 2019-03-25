package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
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
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
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
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
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
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
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
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "adding new field to the type inside existing schema",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type SomeObject
				interface SomeInterface
				input SomeInputObject
			`)),
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
			schema: graphql.MustBuildSchema(nil, []byte(`
				type SomeObject {
					foo: String
				}

				interface SomeInterface {
					foo: String
				}

				input SomeInputObject {
					foo: String
				}
			`)),
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
			errs: (*types.Errors)(nil).
				Add(rules.ExistedFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.ExistedFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.ExistedFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)).
				Add(rules.ExistedFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.ExistedFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.ExistedFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
		},
		{
			msg: "adding fields to existing schema twice",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type SomeObject
				interface SomeInterface
				input SomeInputObject
			`)),
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
			errs: (*types.Errors)(nil).
				Add(rules.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
				Add(rules.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueFieldDefinitionNames)
}
