package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestPossibleNames(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "valid names",
			query: `
				type FooObject {
					bar(arg1: String, arg2: Int): String
				}

				input FooInputObject {
					bar: String
				}

				scalar FooScalar
				interface FooInterface
				union FooUnion
				enum FooEnum
				directive @fooDirective on FIELD
			`,
		},
		{
			msg: "invalid names",
			query: `
				type __FooObject {
					__bar(__arg1: String, __arg2: Int): String
				}

				input __FooInputObject {
					__bar: String
				}

				scalar __FooScalar
				interface __FooInterface
				union __FooUnion
				enum __FooEnum
				directive @__fooDirective on FIELD
			`,
			errs: (*types.Errors)(nil).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooObject", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooInputObject", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__bar", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__arg1", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__arg2", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__bar", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooScalar", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooInterface", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooUnion", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__FooEnum", 0, 0)).
				Add(validation.NameStartsWithTwoUnderscoresError("__fooDirective", 0, 0)),
		},
		{
			msg: "invalid names in extensions",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type FooObject {
					bar: String
				}
			`)),
			query: `
				extend type FooObject {
					__baz: String
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.NameStartsWithTwoUnderscoresError("__baz", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.PossibleNames)
}
