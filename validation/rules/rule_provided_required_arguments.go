package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// ProvidedRequiredArguments ...
func ProvidedRequiredArguments(w *validation.Walker) {}

// ProvidedRequiredArgumentsOnDirectives ...
func ProvidedRequiredArgumentsOnDirectives(w *validation.Walker) {
	w.AddDirectiveEnterEventHandler(func(ctx *validation.Context, d ast.Directive) {
		ctx.DirectiveArguments = make(map[string]struct{}, d.Arguments.Len())

		d.Arguments.ForEach(func(a ast.Argument, _ int) {
			ctx.DirectiveArguments[a.Name] = struct{}{}
		})

		dd, ok := ctx.Schema.Directives[d.Name]
		if !ok {
			return
		}

		dd.ArgumentsDefinition.ForEach(func(ivd ast.InputValueDefinition, _ int) {
			_, provided := ctx.DirectiveArguments[ivd.Name]

			if !provided && ivd.DefaultValue == nil && ivd.Type.NonNullable {
				ctx.AddError(MissingDirectiveArgMessage(d.Name, ivd.Name, ivd.Type.NamedType, 0, 0))
			}
		})
	})
}

// MissingFieldArgMessage ...
func MissingFieldArgMessage(fieldName, argName, typeName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + fieldName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// MissingDirectiveArgMessage ...
func MissingDirectiveArgMessage(directiveName, argName, typeName string, line, col int) types.Error {
	return types.NewError(
		"Directive \"" + directiveName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}
