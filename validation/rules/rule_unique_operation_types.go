package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueOperationTypes ...
func uniqueOperationTypes(w *validation.Walker) {
	w.AddSchemaDefinitionLeaveEventHandler(func(ctx *validation.Context, def *ast.SchemaDefinition) {
		def.RootOperationTypeDefinitions.ForEach(func(rotd ast.RootOperationTypeDefinition, i int) {
			// NOTE: Can't be extending here.
			switch rotd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if ctx.Schema.QueryTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("query", 0, 0))
				}

				ctx.Schema.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if ctx.Schema.MutationTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("mutation", 0, 0))
				}

				ctx.Schema.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if ctx.Schema.SubscriptionTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("subscription", 0, 0))
				}

				ctx.Schema.SubscriptionTypeDefined = true
			}
		})
	})

	w.AddSchemaExtensionLeaveEventHandler(func(ctx *validation.Context, ext *ast.SchemaExtension) {
		ext.OperationTypeDefinitions.ForEach(func(otd ast.OperationTypeDefinition, i int) {
			switch otd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if ctx.Schema.QueryType != nil {
					ctx.AddError(existedOperationTypeMessage("query", 0, 0))
				} else if ctx.Schema.QueryTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("query", 0, 0))
				}

				ctx.Schema.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if ctx.Schema.MutationType != nil {
					ctx.AddError(existedOperationTypeMessage("mutation", 0, 0))
				} else if ctx.Schema.MutationTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("mutation", 0, 0))
				}

				ctx.Schema.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if ctx.Schema.SubscriptionType != nil {
					ctx.AddError(existedOperationTypeMessage("subscription", 0, 0))
				} else if ctx.Schema.SubscriptionTypeDefined {
					ctx.AddError(duplicateOperationTypeMessage("subscription", 0, 0))
				}

				ctx.Schema.SubscriptionTypeDefined = true
			}
		})
	})
}

func duplicateOperationTypeMessage(operation string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one " + operation + " type in schema.",
		// TODO: Location.
	)
}

func existedOperationTypeMessage(operation string, line, col int) graphql.Error {
	return graphql.NewError(
		"Type for " + operation + " already defined in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}
