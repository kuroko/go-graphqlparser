package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueTypeNames ...
func UniqueTypeNames(w *validation.Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		if ctx.SDLContext.KnownTypeNames == nil {
			ctx.SDLContext.KnownTypeNames = make(map[string]struct{})
		}

		if _, ok := ctx.Schema.Types[def.Name]; ok {
			ctx.AddError(ExistedTypeNameError(def.Name, 0, 0))
			return
		}

		if _, ok := ctx.SDLContext.KnownTypeNames[def.Name]; ok {
			ctx.AddError(DuplicateTypeNameError(def.Name, 0, 0))
		} else {
			ctx.SDLContext.KnownTypeNames[def.Name] = struct{}{}
		}
	})
}

// DuplicateTypeNameError ...
func DuplicateTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// ExistedTypeNameError ...
func ExistedTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}
