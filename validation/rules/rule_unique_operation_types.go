package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueOperationTypes ...
func UniqueOperationTypes(w *validation.Walker) {
	w.AddSchemaDefinitionLeaveEventHandler(func(ctx *validation.Context, def *ast.SchemaDefinition) {
		var qryDefined bool
		var mutDefined bool
		var subDefined bool

		def.OperationTypeDefinitions.ForEach(func(rotd ast.OperationTypeDefinition, i int) {
			// NOTE: Can't be extending here.
			switch rotd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if qryDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				qryDefined = true
				ctx.SDLContext.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if mutDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				mutDefined = true
				ctx.SDLContext.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if subDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(rotd.OperationType.String(), 0, 0))
				}

				subDefined = true
				ctx.SDLContext.SubscriptionTypeDefined = true
			}
		})
	})

	w.AddSchemaExtensionLeaveEventHandler(func(ctx *validation.Context, ext *ast.SchemaExtension) {
		ext.OperationTypeDefinitions.ForEach(func(otd ast.OperationTypeDefinition, i int) {
			switch otd.OperationType {
			case ast.OperationDefinitionKindQuery:
				if ctx.Schema.QueryType != nil {
					ctx.AddError(validation.ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.QueryTypeDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.QueryTypeDefined = true
			case ast.OperationDefinitionKindMutation:
				if ctx.Schema.MutationType != nil {
					ctx.AddError(validation.ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.MutationTypeDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.MutationTypeDefined = true
			case ast.OperationDefinitionKindSubscription:
				if ctx.Schema.SubscriptionType != nil {
					ctx.AddError(validation.ExistedOperationTypeError(otd.OperationType.String(), 0, 0))
				} else if ctx.SDLContext.SubscriptionTypeDefined {
					ctx.AddError(validation.DuplicateOperationTypeError(otd.OperationType.String(), 0, 0))
				}

				ctx.SDLContext.SubscriptionTypeDefined = true
			}
		})
	})
}
