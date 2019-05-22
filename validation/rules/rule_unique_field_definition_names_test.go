package rules_test

//func TestUniqueFieldDefinitionNames(t *testing.T) {
//	tt := []ruleTestCase{
//		{
//			msg: "no fields",
//			query: `
//				type SomeObject
//				interface SomeInterface
//				input SomeInputObject
//			`,
//		},
//		{
//			msg: "one field",
//			query: `
//				type SomeObject {
//					foo: String
//				}
//
//				interface SomeInterface {
//					foo: String
//				}
//
//				input SomeInputObject {
//					foo: String
//				}
//			`,
//		},
//		{
//			msg: "multiple fields",
//			query: `
//				type SomeObject {
//					foo: String
//					bar: String
//				}
//
//				interface SomeInterface {
//					foo: String
//					bar: String
//				}
//
//				input SomeInputObject {
//					foo: String
//					bar: String
//				}
//			`,
//		},
//		{
//			msg: "duplicate fields inside the same type definition",
//			query: `
//				type SomeObject {
//					foo: String
//					bar: String
//					foo: String
//				}
//
//				interface SomeInterface {
//					foo: String
//					bar: String
//					foo: String
//				}
//
//				input SomeInputObject {
//					foo: String
//					bar: String
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//		{
//			msg: "extend type with new field",
//			query: `
//				type SomeObject {
//					foo: String
//				}
//				extend type SomeObject {
//					bar: String
//				}
//				extend type SomeObject {
//					baz: String
//				}
//
//				interface SomeInterface {
//					foo: String
//				}
//				extend interface SomeInterface {
//					bar: String
//				}
//				extend interface SomeInterface {
//					baz: String
//				}
//
//				input SomeInputObject {
//					foo: String
//				}
//				extend input SomeInputObject {
//					bar: String
//				}
//				extend input SomeInputObject {
//					baz: String
//				}
//			`,
//		},
//		{
//			msg: "extend type with duplicate field",
//			query: `
//				extend type SomeObject {
//					foo: String
//				}
//				type SomeObject {
//					foo: String
//				}
//
//				extend interface SomeInterface {
//					foo: String
//				}
//				interface SomeInterface {
//					foo: String
//				}
//
//				extend input SomeInputObject {
//					foo: String
//				}
//				input SomeInputObject {
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//		{
//			msg: "duplicate field inside extension",
//			query: `
//				type SomeObject
//				extend type SomeObject {
//					foo: String
//					bar: String
//					foo: String
//				}
//
//				interface SomeInterface
//				extend interface SomeInterface {
//					foo: String
//					bar: String
//					foo: String
//				}
//
//				input SomeInputObject
//				extend input SomeInputObject {
//					foo: String
//					bar: String
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//		{
//			msg: "duplicate field inside different extensions",
//			query: `
//				type SomeObject
//				extend type SomeObject {
//					foo: String
//				}
//				extend type SomeObject {
//					foo: String
//				}
//
//				interface SomeInterface
//				extend interface SomeInterface {
//					foo: String
//				}
//				extend interface SomeInterface {
//					foo: String
//				}
//
//				input SomeInputObject
//				extend input SomeInputObject {
//					foo: String
//				}
//				extend input SomeInputObject {
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//		{
//			msg: "adding new field to the type inside existing schema",
//			schema: mustBuildSchema(nil, []byte(`
//				type SomeObject
//				interface SomeInterface
//				input SomeInputObject
//			`)),
//			query: `
//				extend type SomeObject {
//					foo: String
//				}
//
//				extend interface SomeInterface {
//					foo: String
//				}
//
//				extend input SomeInputObject {
//					foo: String
//				}
//			`,
//		},
//		{
//			msg: "adding conflicting fields to existing schema twice",
//			schema: mustBuildSchema(nil, []byte(`
//				type SomeObject {
//					foo: String
//				}
//
//				interface SomeInterface {
//					foo: String
//				}
//
//				input SomeInputObject {
//					foo: String
//				}
//			`)),
//			query: `
//				extend type SomeObject {
//					foo: String
//				}
//				extend interface SomeInterface {
//					foo: String
//				}
//				extend input SomeInputObject {
//					foo: String
//				}
//
//				extend type SomeObject {
//					foo: String
//				}
//				extend interface SomeInterface {
//					foo: String
//				}
//				extend input SomeInputObject {
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//		{
//			msg: "adding fields to existing schema twice",
//			schema: mustBuildSchema(nil, []byte(`
//				type SomeObject
//				interface SomeInterface
//				input SomeInputObject
//			`)),
//			query: `
//				extend type SomeObject {
//					foo: String
//				}
//				extend type SomeObject {
//					foo: String
//				}
//
//				extend interface SomeInterface {
//					foo: String
//				}
//				extend interface SomeInterface {
//					foo: String
//				}
//
//				extend input SomeInputObject {
//					foo: String
//				}
//				extend input SomeInputObject {
//					foo: String
//				}
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeObject", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInterface", "foo", 0, 0)).
//				Add(validation.DuplicateFieldDefinitionNameError("SomeInputObject", "foo", 0, 0)),
//		},
//	}
//
//	sdlRuleTester(t, tt, func(w *validation.Walker) {})
//}
