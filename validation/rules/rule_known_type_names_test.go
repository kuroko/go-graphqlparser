package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestKnownTypeNames(t *testing.T) {
	t.Run("query document", func(t *testing.T) {
		tt := []ruleTestCase{
			{
				msg: "known type names are valid",
				query: `
					query Foo($var: String, $required: [String!]!) {
						user(id: 4) {
							pets { ... on Pet { name }, ...PetFields, ... { name } }
						}
					}

					fragment PetFields on Pet {
						name
					}
				`,
			},
		}

		queryRuleTester(t, tt, rules.KnownTypeNames)
	})

	t.Run("sdl document", func(t *testing.T) {
		tt := []ruleTestCase{}

		sdlRuleTester(t, tt, rules.KnownTypeNames)
	})
}
