package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
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
				Add(rules.DuplicateOperationTypeError("query", 0, 0)).
				Add(rules.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(rules.DuplicateOperationTypeError("subscription", 0, 0)),
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
				Add(rules.DuplicateOperationTypeError("query", 0, 0)).
				Add(rules.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(rules.DuplicateOperationTypeError("subscription", 0, 0)),
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
				Add(rules.DuplicateOperationTypeError("query", 0, 0)).
				Add(rules.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(rules.DuplicateOperationTypeError("subscription", 0, 0)).
				Add(rules.DuplicateOperationTypeError("query", 0, 0)).
				Add(rules.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(rules.DuplicateOperationTypeError("subscription", 0, 0)),
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
				Add(rules.DuplicateOperationTypeError("query", 0, 0)).
				Add(rules.DuplicateOperationTypeError("mutation", 0, 0)).
				Add(rules.DuplicateOperationTypeError("subscription", 0, 0)),
		},
		{
			msg:    "define schema inside extension SDL",
			schema: &types.Schema{},
			query: `
			schema {
				query: Query
				mutation: Mutation
				subscription: Subscription
			}
			`,
		},
		{
			msg:    "define and extend schema inside extension SDL",
			schema: &types.Schema{},
			query: `
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
			// TODO: Maybe replace with something that builds this from a schema string?
			schema: &types.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Query",
				},
				MutationType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Mutation",
				},
				SubscriptionType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Subscription",
				},
			},
			query: `
			extend schema {
				query: Foo
				mutation: Foo
				subscription: Foo
			}
			`,
			errs: (*types.Errors)(nil).
				Add(rules.ExistedOperationTypeError("query", 0, 0)).
				Add(rules.ExistedOperationTypeError("mutation", 0, 0)).
				Add(rules.ExistedOperationTypeError("subscription", 0, 0)),
		},
		{
			msg: "adding conflicting operation types to existing schema twice",
			// TODO: Maybe replace with something that builds this from a schema string?
			schema: &types.Schema{
				QueryType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Query",
				},
				MutationType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Mutation",
				},
				SubscriptionType: &ast.Type{
					Kind:      ast.TypeKindNamed,
					NamedType: "Subscription",
				},
			},
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
				Add(rules.ExistedOperationTypeError("query", 0, 0)).
				Add(rules.ExistedOperationTypeError("mutation", 0, 0)).
				Add(rules.ExistedOperationTypeError("subscription", 0, 0)).
				Add(rules.ExistedOperationTypeError("query", 0, 0)).
				Add(rules.ExistedOperationTypeError("mutation", 0, 0)).
				Add(rules.ExistedOperationTypeError("subscription", 0, 0)),
		},
	}

	sdlRuleTester(t, tt, rules.UniqueOperationTypes)
}
