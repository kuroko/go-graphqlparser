package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueInputFieldNames(t *testing.T) {
	t.Run("query documents", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "input object with fields",
				query: `
					{
						field1(arg: { f: true })
						field2 @foo(arg: { f: true })
					}
				`,
			},
			{
				msg: "same input object within two args",
				query: `
					{
						field1(arg1: { f: true }, arg2: { f: true })
						field2 @foo(arg1: { f: true }, arg2: { f: true })
					}
				`,
			},
			{
				msg: "multiple input object fields",
				query: `
					{
						field1(arg: { f1: "value", f2: "value", f3: "value" })
						field2 @foo(arg: { f1: "value", f2: "value", f3: "value" })
					}
				`,
			},
			{
				msg: "allows for nested input objects with similar fields",
				query: `
					{
						field1(arg: {
							deep: {
								deep: {
									id: 1
								}
								id: 1
							}
							id: 1
						})
						field2 @foo(arg: {
							deep: {
								deep: {
									id: 1
								}
								id: 1
							}
							id: 1
						})
					}
				`,
			},
			{
				msg: "duplicate input object fields",
				query: `
					{
						field1(arg: { f1: "value", f1: "value" })
						field2 @foo(arg: { f1: "value", f1: "value" })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)),
			},
			{
				msg: "many duplicate input object fields",
				query: `
					{
						field1(arg: { f1: "value", f1: "value", f1: "value" })
						field2 @foo(arg: { f1: "value", f1: "value", f1: "value" })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)),
			},
			{
				msg: "nested duplicate input object fields",
				query: `
					{
						field1(arg: { f1: { f2: "value", f2: "value" } })
						field2 @foo(arg: { f1: { f2: "value", f2: "value" } })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f2", 0, 0)).
					Add(validation.DuplicateInputFieldError("f2", 0, 0)),
			},
		}

		queryRuleTester(t, tt, rules.UniqueInputFieldNames)
	})

	t.Run("sdl documents", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "input object with fields",
				query: `
					type SomeType {
						field: String @foo(arg: { f: true })
					}
				`,
			},
			{
				msg: "same input object within two args",
				query: `
					type SomeType {
						field: String @foo(arg1: { f: true }, arg2: { f: true })
					}
				`,
			},
			{
				msg: "multiple input object fields",
				query: `
					type SomeType {
						field: String @foo(arg: { f1: "value", f2: "value", f3: "value" })
					}
				`,
			},
			{
				msg: "allows for nested input objects with similar fields",
				query: `
					type SomeType {
						field: String @foo(arg: {
							deep: {
								deep: {
									id: 1
								}
								id: 1
							}
							id: 1
						})
					}
				`,
			},
			{
				msg: "duplicate input object fields",
				query: `
					type SomeType {
						field: String @foo(arg: { f1: "value", f1: "value" })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)),
			},
			{
				msg: "many duplicate input object fields",
				query: `
					type SomeType {
						field: String @foo(arg: { f1: "value", f1: "value", f1: "value" })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)).
					Add(validation.DuplicateInputFieldError("f1", 0, 0)),
			},
			{
				msg: "nested duplicate input object fields",
				query: `
					type SomeType {
						field: String @foo(arg: { f1: { f2: "value", f2: "value" } })
					}
				`,
				errs: (*types.Errors)(nil).
					Add(validation.DuplicateInputFieldError("f2", 0, 0)),
			},
		}

		sdlRuleTester(t, tt, rules.UniqueInputFieldNames)
	})
}
