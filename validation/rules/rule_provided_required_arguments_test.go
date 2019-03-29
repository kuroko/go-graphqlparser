package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestProvidedRequiredArguments(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "ignores unknown arguments",
			query: `
			{
				dog {
				  isHousetrained(unknownArgument: true)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "arg on optional arg",
			query: `
			{
				dog {
				  isHousetrained(atOtherHomes: true)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "no arg on optional arg",
			query: `
			{
				dog {
				  isHousetrained
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "no arg on non-null field with default",
			query: `
			{
				complicatedArgs {
				  nonNullFieldWithDefault
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "multiple args",
			query: `
			{
				complicatedArgs {
				  multipleReqs(req1: 1, req2: 2)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "multiple args reverse order",
			query: `
			{
				complicatedArgs {
				  multipleReqs(req2: 2, req1: 1)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "no args on multiple optional",
			query: `
			{
				complicatedArgs {
				  multipleReqs(req2: 2, req1: 1)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "one arg on multiple optional",
			query: `
			{
				complicatedArgs {
				  multipleOpts(opt1: 1)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "second arg on multiple optional",
			query: `
			{
				complicatedArgs {
				  multipleOpts(opt2: 1)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "multiple reqs on mixedList",
			query: `
			{
				complicatedArgs {
				  multipleOptAndReq(req1: 3, req2: 4)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "multiple reqs and one opt on mixedList",
			query: `
			{
				complicatedArgs {
				  multipleOptAndReq(req1: 3, req2: 4, opt1: 5)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "all reqs and opts on mixedList",
			query: `
			{
				complicatedArgs {
				  multipleOptAndReq(req1: 3, req2: 4, opt1: 5, opt2: 6)
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "missing one non-nullable argument",
			query: `
			{
				complicatedArgs {
				  multipleReqs(req2: 2)
				}
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingFieldArgMessage("multipleReqs", "req1", "Int!", 0, 0)),
		},
		{
			msg: "missing multiple non-nullable arguments",
			query: `
			{
				complicatedArgs {
				  multipleReqs
				}
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingFieldArgMessage("multipleReqs", "req1", "Int!", 0, 0)).
				Add(rules.MissingFieldArgMessage("multipleReqs", "req2", "Int!", 0, 0)),
		},
		{
			msg: "incorrect value and missing argument",
			query: `
			{
				complicatedArgs {
				  multipleReqs(req1: "one")
				}
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingFieldArgMessage("multipleReqs", "req2", "Int!", 0, 0)),
		},
		{
			msg: "ignores unknown directives",
			query: `
			{
				dog @unknown
			}
			`,
			errs: nil,
		},
		{
			msg: "with directives of valid types",
			query: `
			{
				dog @include(if: true) {
				  name
				}
				human @skip(if: false) {
				  name
				}
			}
			`,
			errs: nil,
		},
		{
			msg: "with directive with missing types",
			query: `
			{
				dog @include {
				  name @skip
				}
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("include", "if", "Boolean!", 0, 0)).
				Add(rules.MissingDirectiveArgMessage("skip", "if", "Boolean!", 0, 0)),
		},
	}

	queryRuleTester(t, tt, rules.ProvidedRequiredArguments)

	sdlTT := []ruleTestCase{
		{
			msg: "missing optional args on directive defined inside SDL",
			query: `
			type Query {
				foo: String @test
			}
			directive @test(arg1: String, arg2: String! = "") on FIELD_DEFINITION
			`,
			errs: nil,
		},
		{
			msg: "missing arg on directive defined inside SDL",
			query: `
			type Query {
				foo: String @test
			}
			directive @test(arg: String!) on FIELD_DEFINITION
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("test", "arg", "String!", 0, 0)),
		},
		{
			msg: "missing arg on standard directive",
			query: `
			type Query {
				foo: String @include
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("include", "if", "Boolean!", 0, 0)),
		},
		{
			msg: "missing arg on overridden standard directive",
			query: `
			type Query {
				foo: String @deprecated
			  }
			directive @deprecated(reason: String!) on FIELD
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("deprecated", "reason", "String!", 0, 0)),
		},
		{
			msg: "missing arg on directive defined in schema extension",
			schema: graphql.MustBuildSchema(nil, []byte(`
			type Query {
				foo: String
			}
			`)),
			query: `
			directive @test(arg: String!) on OBJECT
			extend type Query  @test
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("test", "arg", "String!", 0, 0)),
		},
		{
			msg: "missing arg on directive used in schema extension",
			schema: graphql.MustBuildSchema(nil, []byte(`
			directive @test(arg: String!) on OBJECT
			type Query {
				foo: String
			}
			`)),
			query: `
			extend type Query @test
			`,
			errs: (*types.Errors)(nil).
				Add(rules.MissingDirectiveArgMessage("test", "arg", "String!", 0, 0)),
		},
	}

	queryRuleTester(t, sdlTT, rules.ProvidedRequiredArgumentsOnDirectives)
}
