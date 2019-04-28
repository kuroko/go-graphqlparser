package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueArgumentNames(t *testing.T) {
	t.Run("query documents", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "no arguments on field",
				query: `
					{
						field
					}
				`,
			},
			{
				msg: "no arguments on directive",
				query: `
					{
						field @directive
					}
				`,
			},
			{
				msg: "argument on field",
				query: `
					{
						field(arg: "value")
					}
				`,
			},
			{
				msg: "argument on directive",
				query: `
					{
						field @directive(arg: "value")
					}
				`,
			},
			{
				msg: "same argument on two fields",
				query: `
					{
						one: field(arg: "value")
						two: field(arg: "value")
					}
				`,
			},
			{
				msg: "same argument on field and directive",
				query: `
					{
						field(arg: "value") @directive(arg: "value")
					}
				`,
			},
			{
				msg: "same argument on two directives",
				query: `
					{
						field @directive1(arg: "value") @directive2(arg: "value")
					}
				`,
			},
			{
				msg: "multiple field arguments",
				query: `
					{
						field(arg1: "value", arg2: "value", arg3: "value")
					}
				`,
			},
			{
				msg: "multiple directive arguments",
				query: `
					{
						field @directive(arg1: "value", arg2: "value", arg3: "value")
					}
				`,
			},
			{
				msg: "duplicate field arguments",
				query: `
					{
						field(arg1: "value", arg1: "value")
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "many duplicate field arguments",
				query: `
					{
						field(arg1: "value", arg1: "value", arg1: "value")
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "duplicate directive arguments",
				query: `
					{
						field @directive(arg1: "value", arg1: "value")
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "many duplicate directive arguments",
				query: `
					{
						field @directive(arg1: "value", arg1: "value", arg1: "value")
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
		}

		queryRuleTester(t, tt, rules.UniqueArgumentNames)
	})

	t.Run("sdl documents", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "no arguments on field",
				query: `
					type SomeType {
						field: String
					}

					interface SomeInterface {
						field: String
					}
				`,
			},
			{
				msg: "no arguments on directive",
				query: `
					directive @foo on SCHEMA
				`,
			},
			{
				msg: "one argument on field",
				query: `
					type SomeType {
						field(arg1: String): String
					}

					interface SomeInterface {
						field(arg1: String): String
					}
				`,
			},
			{
				msg: "one argument on directive",
				query: `
					directive @foo(bar: String) on SCHEMA
				`,
			},
			{
				msg: "same argument on two fields",
				query: `
					type SomeType {
						field1(arg1: String): String
						field2(arg1: String): String
					}

					interface SomeInterface {
						field1(arg1: String): String
						field2(arg1: String): String
					}
				`,
			},
			{
				msg: "same argument on two directives",
				query: `
					directive @foo(baz: String) on SCHEMA
					directive @bar(baz: String) on SCHEMA
				`,
			},
			{
				msg: "multiple field arguments",
				query: `
					type SomeType {
						field(arg1: String, arg2: String): String
					}

					interface SomeInterface {
						field(arg1: String, arg2: String): String
					}
				`,
			},
			{
				msg: "multiple directive arguments",
				query: `
					directive @foo(arg1: String, arg2: String) on SCHEMA
				`,
			},
			{
				msg: "duplicate field arguments",
				query: `
					type SomeType {
						field(arg1: String, arg1: String): String
					}

					interface SomeInterface {
						field(arg1: String, arg1: String): String
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "duplicate directive arguments",
				query: `
					directive @foo(arg1: String, arg1: String) on SCHEMA
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "many duplicate field arguments",
				query: `
					type SomeType {
						field(arg1: String, arg1: String, arg1: String): String
					}

					interface SomeInterface {
						field(arg1: String, arg1: String, arg1: String): String
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
			{
				msg: "many duplicate directive arguments",
				query: `
					directive @foo(arg1: String, arg1: String, arg1: String) on SCHEMA
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateArgError("arg1", 0, 0)).
					Add(validation.DuplicateArgError("arg1", 0, 0)),
			},
		}

		sdlRuleTester(t, tt, rules.UniqueArgumentNames)
	})
}
