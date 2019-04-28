package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueArgumentNames ...
func UniqueArgumentNames(w *validation.Walker) {
	w.AddFieldSelectionEnterEventHandler(func(ctx *validation.Context, sel ast.Selection) {
		validateArguments(ctx, sel.Arguments)
	})

	w.AddDirectiveEnterEventHandler(func(ctx *validation.Context, dir ast.Directive) {
		validateArguments(ctx, dir.Arguments)
	})

	// NOTE: This is not validated in graphql-js, but will silently cause issues.
	w.AddFieldDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.FieldDefinition) {
		validateInputValueDefinitions(ctx, def.ArgumentsDefinition)
	})

	// NOTE: This is not validated in graphql-js, but will silently cause issues.
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
		validateInputValueDefinitions(ctx, def.ArgumentsDefinition)
	})
}

// validateArguments ...
func validateArguments(ctx *validation.Context, arguments *ast.Arguments) {
	knownArgNames := make(map[string]struct{})

	arguments.ForEach(func(a ast.Argument, i int) {
		argName := a.Name

		if _, ok := knownArgNames[argName]; ok {
			ctx.AddError(validation.DuplicateArgError(argName, 0, 0))
		} else {
			knownArgNames[argName] = struct{}{}
		}
	})
}

// validateInputValueDefinitions ...
func validateInputValueDefinitions(ctx *validation.Context, ivds *ast.InputValueDefinitions) {
	knownArgNames := make(map[string]struct{})

	ivds.ForEach(func(ivd ast.InputValueDefinition, i int) {
		argName := ivd.Name

		if _, ok := knownArgNames[argName]; ok {
			ctx.AddError(validation.DuplicateArgError(argName, 0, 0))
		} else {
			knownArgNames[argName] = struct{}{}
		}
	})
}
