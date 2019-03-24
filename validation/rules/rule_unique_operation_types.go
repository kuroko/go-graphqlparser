package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueOperationTypes ...
func UniqueOperationTypes(w *validation.Walker) {
	w.AddSchemaDefinitionLeaveEventHandler(func(ctx *validation.Context, def *ast.SchemaDefinition) {
		def.RootOperationTypeDefinitions.ForEach(func(rotd ast.RootOperationTypeDefinition, i int) {
			// NOTE: Can't be extending here.
			switch rotd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if ctx.SDLContext.QueryTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if ctx.SDLContext.MutationTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if ctx.SDLContext.SubscriptionTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.SubscriptionTypeDefined = true
			}
		})
	})

	w.AddSchemaExtensionLeaveEventHandler(func(ctx *validation.Context, ext *ast.SchemaExtension) {
		ext.OperationTypeDefinitions.ForEach(func(otd ast.OperationTypeDefinition, i int) {
			switch otd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if ctx.Schema.QueryType != nil {
					ctx.AddError(ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.QueryTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if ctx.Schema.MutationType != nil {
					ctx.AddError(ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.MutationTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if ctx.Schema.SubscriptionType != nil {
					ctx.AddError(ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.SubscriptionTypeDefined {
					ctx.AddError(DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.SubscriptionTypeDefined = true
			}
		})
	})
}

// DuplicateOperationTypeError ...
func DuplicateOperationTypeError(operation string, line, col int) types.Error {
	return types.NewError(
		"There can be only one " + operation + " type in schema.",
		// TODO: Location.
	)
}

// ExistedOperationTypeError ...
func ExistedOperationTypeError(operation string, line, col int) types.Error {
	return types.NewError(
		"Type for " + operation + " already defined in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}
