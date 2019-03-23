package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueDirectiveNames ...
func uniqueDirectiveNames(w *validation.Walker) {
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
		if ctx.SDLContext.KnownDirectiveNames == nil {
			ctx.SDLContext.KnownDirectiveNames = make(map[string]struct{})
		}

		directiveName := def.Name

		if _, ok := ctx.Schema.Directives[directiveName]; ok {
			ctx.AddError(existedDirectiveNameMessage(directiveName, 0, 0))
			return
		}

		if _, ok := ctx.SDLContext.KnownDirectiveNames[directiveName]; ok {
			ctx.AddError(duplicateDirectiveNameMessage(directiveName, 0, 0))
		} else {
			ctx.SDLContext.KnownDirectiveNames[directiveName] = struct{}{}
		}
	})
}

// duplicateDirectiveNameMessage ...
func duplicateDirectiveNameMessage(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one directive named \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// existedDirectiveNameMessage ...
func existedDirectiveNameMessage(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Directive \"" + directiveName + "\" already exists in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}
