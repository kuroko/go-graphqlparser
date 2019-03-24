package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// Lone Schema definition
//
// A GraphQL document is only valid if it contains only one schema definition.
func LoneSchemaDefinition(w *validation.Walker) {
	w.AddSchemaDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.SchemaDefinition) {
		// NOTE: If ctx.SDLContext is nil here, we should panic, this is an SDL only rule.

		if ctx.SDLContext.IsExtending {
			ctx.AddError(CanNotDefineSchemaWithinExtensionError(0, 0))
			return
		}

		if ctx.SDLContext.HasSeenSchemaDefinition {
			ctx.AddError(SchemaDefinitionNotAloneError(0, 0))
			return
		}

		ctx.SDLContext.HasSeenSchemaDefinition = true
	})
}

func SchemaDefinitionNotAloneError(line, col int) types.Error {
	return types.NewError(
		"Must provide only one schema definition.",
		// TODO: Location.
	)
}

func CanNotDefineSchemaWithinExtensionError(line, col int) types.Error {
	return types.NewError(
		"Cannot define a new schema within a schema extension.",
		// TODO: Location.
	)
}
