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
		directiveArguments := make(map[string]struct{}, d.Arguments.Len())

		d.Arguments.ForEach(func(a ast.Argument, _ int) {
			directiveArguments[a.Name] = struct{}{}
		})

		var dd *ast.DirectiveDefinition
		var ok bool

		dd, ok = ctx.Schema.Directives[d.Name]
		if !ok {
			dd, ok = ctx.SDLContext.DirectiveDefinitions[d.Name]
		}
		if !ok {
			dd, ok = types.SpecifiedDirectives()[d.Name]
		}
		if !ok {
			return
		}

		dd.ArgumentsDefinition.ForEach(func(ivd ast.InputValueDefinition, _ int) {
			_, provided := directiveArguments[ivd.Name]

			if !provided && ivd.DefaultValue == nil && ivd.Type.NonNullable {
				ctx.AddError(validation.MissingDirectiveArgError(d.Name, ivd.Name, ivd.Type.String(), 0, 0))
			}
		})
	})
}
