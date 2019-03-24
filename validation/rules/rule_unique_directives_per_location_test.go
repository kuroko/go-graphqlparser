package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueDirectivesPerLocation(t *testing.T) {
	t.Run("query document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "no directives",
				query: `
					fragment Test on Type {
						field
					}
				`,
			},
			{
				msg: "unique directives in different locations",
				query: `
					fragment Test on Type @directiveA {
						field @directiveB
					}
				`,
			},
			{
				msg: "unique directives in same locations",
				query: `
					fragment Test on Type @directiveA @directiveB {
						field @directiveA @directiveB
					}
				`,
			},
			{
				msg: "same directives in different locations",
				query: `
					fragment Test on Type @directiveA {
						field @directiveA
					}
				`,
			},
			{
				msg: "same directives in similar locations",
				query: `
					fragment Test on Type {
						field @directive
						field @directive
					}
				`,
			},
			{
				msg: "duplicate directives in one location",
				query: `
					fragment Test on Type {
						field @directive @directive
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)),
			},
			{
				msg: "many duplicate directives in one location",
				query: `
					fragment Test on Type {
						field @directive @directive @directive
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)),
			},
			{
				msg: "different duplicate directives in one location",
				query: `
					fragment Test on Type {
						field @directiveA @directiveB @directiveA @directiveB
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.DuplicateDirectiveError("directiveA", 0, 0)).
					Add(rules.DuplicateDirectiveError("directiveB", 0, 0)),
			},
			{
				msg: "duplicate directives in many locations",
				query: `
					fragment Test on Type @directive @directive {
						field @directive @directive
					}
				`,
				errs: (*types.Errors)(nil).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)),
			},
		}

		queryRuleTester(t, tt, rules.UniqueDirectivesPerLocation)
	})

	t.Run("sdl document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "duplicate directives in many locations",
				query: `
					schema @directive @directive { query: Dummy }
					extend schema @directive @directive

					scalar TestScalar @directive @directive
					extend scalar TestScalar @directive @directive

					type TestObject @directive @directive
					extend type TestObject @directive @directive

					interface TestInterface @directive @directive
					extend interface TestInterface @directive @directive

					union TestUnion @directive @directive
					extend union TestUnion @directive @directive

					input TestInput @directive @directive
					extend input TestInput @directive @directive
				`,
				errs: (*types.Errors)(nil).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)).
					Add(rules.DuplicateDirectiveError("directive", 0, 0)),
			},
		}

		sdlRuleTester(t, tt, rules.UniqueDirectivesPerLocation)
	})
}
