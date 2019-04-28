package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestKnownArgumentNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, rules.KnownArgumentNames)
}

func TestKnownArgumentNamesOnDirectives(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "known arg on directive defined inside SDL",
			query: `
				type Query {
					foo: String @test(arg: "")
				}

				directive @test(arg: String) on FIELD_DEFINITION
			`,
		},
		{
			msg: "unknown arg on directive defined inside SDL",
			query: `
				type Query {
					foo: String @test(unknown: "")
				}

				directive @test(arg: String) on FIELD_DEFINITION
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnknownDirectiveArgError("unknown", "test", 0, 0)),
		},
		// NOTE: We don't do suggestions in this library, it's a bit pointless, and uses way too
		// many resources to bother with it. Tooling can do this for you without the server.
		//{
		//	msg: "misspelled arg name is reported on directive defined inside SDL",
		//},
		{
			msg: "unknown arg on standard directive",
			query: `
				type Query {
					foo: String @deprecated(unknown: "")
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnknownDirectiveArgError("unknown", "deprecated", 0, 0)),
		},
		// TODO: Check if this is behaving as intended...
		{
			msg: "unknown arg on overridden standard directive",
			query: `
				type Query {
					foo: String @deprecated(reason: "")
				}

				directive @deprecated(arg: String) on FIELD
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnknownDirectiveArgError("reason", "deprecated", 0, 0)),
		},
		{
			msg: "unknown arg on directive defined in schema extension",
			schema: graphql.MustBuildSchema(nil, []byte(`
				type Query {
					foo: String
				}
			`)),
			query: `
				directive @test(arg: String) on OBJECT

				extend type Query @test(unknown: "")
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnknownDirectiveArgError("unknown", "test", 0, 0)),
		},
		{
			msg: "unknown arg on directive used in schema extension",
			schema: graphql.MustBuildSchema(nil, []byte(`
				directive @test(arg: String) on OBJECT

				type Query {
					foo: String
				}
			`)),
			query: `
				extend type Query @test(unknown: "")
			`,
			errs: (*types.Errors)(nil).
				Add(validation.UnknownDirectiveArgError("unknown", "test", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.KnownArgumentNamesOnDirectives)
}
