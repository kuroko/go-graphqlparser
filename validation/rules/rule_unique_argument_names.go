package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueArgumentNames ...
func UniqueArgumentNames(w *validation.Walker) {
	w.AddFieldSelectionEnterEventHandler(func(ctx *validation.Context, sel ast.Selection) {
		ctx.KnownArgNames = make(map[string]struct{})
	})

	w.AddDirectiveEnterEventHandler(func(ctx *validation.Context, dir ast.Directive) {
		ctx.KnownArgNames = make(map[string]struct{})
	})

	// NOTE: This is not validated in graphql-js, but will silently cause issues.
	w.AddFieldDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.FieldDefinition) {
		ctx.KnownArgNames = make(map[string]struct{})
	})

	// NOTE: This is not validated in graphql-js, but will silently cause issues.
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
		ctx.KnownArgNames = make(map[string]struct{})
	})

	w.AddArgumentEnterEventHandler(func(ctx *validation.Context, arg ast.Argument) {
		argName := arg.Name

		if _, ok := ctx.KnownArgNames[argName]; ok {
			ctx.AddError(DuplicateArgError(argName, 0, 0))
		} else {
			ctx.KnownArgNames[argName] = struct{}{}
		}
	})

	w.AddInputValueDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.InputValueDefinition) {
		argName := def.Name

		if _, ok := ctx.KnownArgNames[argName]; ok {
			ctx.AddError(DuplicateArgError(argName, 0, 0))
		} else {
			ctx.KnownArgNames[argName] = struct{}{}
		}
	})
}

// DuplicateArgError ...
func DuplicateArgError(argName string, line, col int) types.Error {
	return types.NewError("There can be only one argument named \"" + argName + "\".")
}
