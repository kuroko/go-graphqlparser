package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueOperationTypes(t *testing.T) {
	tt := []ruleTestCase{
		{
			msg: "no schema definition",
			query: `
				type Foo
			`,
		},
		{
			msg: "schema definition with all types",
			query: `
				type Foo

				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}
			`,
		},
		{
			msg: "schema definition with single extension",
			query: `
				type Foo

				schema {
					query: Query
				}

				extend schema {
					mutation: Mutation
					subscription: Subscription
				}
			`,
		},
		{
			msg: "schema definition with separate extensions",
			query: `
				type Foo

				schema {
					query: Query
				}

				extend schema {
					mutation: Mutation
				}

				extend schema {
					subscription: Subscription
				}
			`,
		},
		{
			msg: "extend schema before definition",
			query: `
				type Foo

				extend schema {
					mutation: Mutation
				}

				extend schema {
					subscription: Subscription
				}

				schema {
					query: Query
				}
			`,
		},
		{
			msg: "duplicate operation types inside single schema definition",
			query: `
				type Foo

				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription

					query: Query
					mutation: Mutation
					subscription: Subscription
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.DuplicateOperationTypeError("query", 0, 0)).
				Add(validation.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(validation.DuplicateOperationTypeError("subscription", 0, 0)),
		},
		{
			msg: "duplicate operation types inside schema extension",
			query: `
				type Foo

				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}

				extend schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.DuplicateOperationTypeError("query", 0, 0)).
				Add(validation.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(validation.DuplicateOperationTypeError("subscription", 0, 0)),
		},
		{
			msg: "duplicate operation types inside schema extension twice",
			query: `
				type Foo

				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}

				extend schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}

				extend schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.DuplicateOperationTypeError("query", 0, 0)).
				Add(validation.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(validation.DuplicateOperationTypeError("subscription", 0, 0)).
				Add(validation.DuplicateOperationTypeError("query", 0, 0)).
				Add(validation.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(validation.DuplicateOperationTypeError("subscription", 0, 0)),
		},
		{
			msg: "duplicate operation types inside second schema extension",
			query: `
				type Foo

				schema {
					query: Query
				}

				extend schema {
					mutation: Mutation
					subscription: Subscription
				}

				extend schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.DuplicateOperationTypeError("query", 0, 0)).
				Add(validation.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(validation.DuplicateOperationTypeError("subscription", 0, 0)),
		},
		// These tests provide invalid schemas, and part of our validation occurs before these rules
		// are run, so we can't include these tests.
		//{
		//	msg:    "define schema inside extension SDL",
		//	schema: &types.Schema{},
		//	query: `
		//		schema {
		//			query: Query
		//			mutation: Mutation
		//			subscription: Subscription
		//		}
		//	`,
		//},
		//{
		//	msg:    "define and extend schema inside extension SDL",
		//	schema: &types.Schema{},
		//	query: `
		//		schema {
		//			query: Query
		//		}
		//
		//		extend schema {
		//			mutation: Mutation
		//		}
		//
		//		extend schema {
		//			subscription: Subscription
		//		}
		//	`,
		//},
		{
			msg:    "adding new operation types to existing schema",
			schema: &types.Schema{},
			query: `
				extend schema {
					mutation: Mutation
				}

				extend schema {
					subscription: Subscription
				}
			`,
		},
		{
			msg: "adding conflicting operation types to existing schema",
			schema: graphql.MustBuildSchema(nil, []byte(`
				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}

				type Query
				type Mutation
				type Subscription
			`)),
			query: `
				extend schema {
					query: Foo
					mutation: Foo
					subscription: Foo
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.ExistedOperationTypeError("query", 0, 0)).
				Add(validation.ExistedOperationTypeError("mutation", 0, 0)).
				Add(validation.ExistedOperationTypeError("subscription", 0, 0)),
		},
		{
			msg: "adding conflicting operation types to existing schema twice",
			schema: graphql.MustBuildSchema(nil, []byte(`
				schema {
					query: Query
					mutation: Mutation
					subscription: Subscription
				}

				type Query
				type Mutation
				type Subscription
			`)),
			query: `
				extend schema {
					query: Foo
					mutation: Foo
					subscription: Foo
				}

				extend schema {
					query: Foo
					mutation: Foo
					subscription: Foo
				}
			`,
			errs: (*types.Errors)(nil).
				Add(validation.ExistedOperationTypeError("query", 0, 0)).
				Add(validation.ExistedOperationTypeError("mutation", 0, 0)).
				Add(validation.ExistedOperationTypeError("subscription", 0, 0)).
				Add(validation.ExistedOperationTypeError("query", 0, 0)).
				Add(validation.ExistedOperationTypeError("mutation", 0, 0)).
				Add(validation.ExistedOperationTypeError("subscription", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueOperationTypes)
}
