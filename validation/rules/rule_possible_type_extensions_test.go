package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestPossibleTypeExtensions(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no extensions",
			query: `
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
			`,
		},
		{
			msg: "one extension per type",
			query: `
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
				extend scalar FooScalar @dummy
				extend type FooObject @dummy
				extend interface FooInterface @dummy
				extend union FooUnion @dummy
				extend enum FooEnum @dummy
				extend input FooInputObject @dummy
			`,
		},
		{
			msg: "many extensions per type",
			query: `
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
				extend scalar FooScalar @dummy
				extend type FooObject @dummy
				extend interface FooInterface @dummy
				extend union FooUnion @dummy
				extend enum FooEnum @dummy
				extend input FooInputObject @dummy
				extend scalar FooScalar @dummy
				extend type FooObject @dummy
				extend interface FooInterface @dummy
				extend union FooUnion @dummy
				extend enum FooEnum @dummy
				extend input FooInputObject @dummy
			`,
		},
		{
			msg: "extending unknown type",
			query: `
				type Known
				extend scalar Unknown @dummy
				extend type Unknown @dummy
				extend interface Unknown @dummy
				extend union Unknown @dummy
				extend enum Unknown @dummy
				extend input Unknown @dummy
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)),
		},
		{
			msg: "does not consider non-type definitions",
			query: `
				query Foo { __typename }
				fragment Foo on Query { __typename }
				directive @Foo on SCHEMA
				extend scalar Foo @dummy
				extend type Foo @dummy
				extend interface Foo @dummy
				extend union Foo @dummy
				extend enum Foo @dummy
				extend input Foo @dummy
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Foo", 0, 0)),
		},
		{
			msg: "extending with different kinds",
			query: `
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
				extend type FooScalar @dummy
				extend interface FooObject @dummy
				extend union FooInterface @dummy
				extend enum FooUnion @dummy
				extend input FooEnum @dummy
				extend scalar FooInputObject @dummy
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.ExtendingDifferentTypeKindError("FooScalar", "scalar", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooObject", "object", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooInterface", "interface", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooUnion", "union", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooEnum", "enum", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooInputObject", "input object", 0, 0)),
		},
		{
			msg: "extending types within existing schema",
			schema: mustBuildSchema(nil, []byte(`
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
			`)),
			query: `
				extend scalar FooScalar @dummy
				extend type FooObject @dummy
				extend interface FooInterface @dummy
				extend union FooUnion @dummy
				extend enum FooEnum @dummy
				extend input FooInputObject @dummy
			`,
		},
		{
			msg:    "extending unknown types within existing schema",
			schema: mustBuildSchema(nil, []byte(`type Known`)),
			query: `
				extend scalar Unknown @dummy
				extend type Unknown @dummy
				extend interface Unknown @dummy
				extend union Unknown @dummy
				extend enum Unknown @dummy
				extend input Unknown @dummy
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)).
				Add(validation.ExtendingUnknownTypeError("Unknown", 0, 0)),
		},
		{
			msg: "extending types with different kinds within existing schema",
			schema: mustBuildSchema(nil, []byte(`
				scalar FooScalar
				type FooObject
				interface FooInterface
				union FooUnion
				enum FooEnum
				input FooInputObject
			`)),
			query: `
				extend type FooScalar @dummy
				extend interface FooObject @dummy
				extend union FooInterface @dummy
				extend enum FooUnion @dummy
				extend input FooEnum @dummy
				extend scalar FooInputObject @dummy
			`,
			errs: (*graphql.Errors)(nil).
				Add(validation.ExtendingDifferentTypeKindError("FooScalar", "scalar", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooObject", "object", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooInterface", "interface", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooUnion", "union", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooEnum", "enum", 0, 0)).
				Add(validation.ExtendingDifferentTypeKindError("FooInputObject", "input object", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.PossibleTypeExtensions)
}
