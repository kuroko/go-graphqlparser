package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueTypeNames ...
func uniqueTypeNames(w *validation.Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		if ctx.SDLContext.KnownTypeNames == nil {
			ctx.SDLContext.KnownTypeNames = make(map[string]struct{})
		}

		// TODO: Can this ever be nil?
		if _, ok := ctx.Schema.Types[def.Name]; ok {
			ctx.AddError(existedTypeNameMessage(def.Name, 0, 0))
			return
		}

		if _, ok := ctx.SDLContext.KnownTypeNames[def.Name]; ok {
			ctx.AddError(duplicateTypeNameMessage(def.Name, 0, 0))
		} else {
			ctx.SDLContext.KnownTypeNames[def.Name] = struct{}{}
		}
	})
}

// duplicateTypeNameMessage ...
func duplicateTypeNameMessage(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// existedTypeNameMessage ...
func existedTypeNameMessage(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}
