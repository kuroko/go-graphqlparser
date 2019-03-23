package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueDirectivesPerLocation ...
func uniqueDirectivesPerLocation(w *validation.Walker) {
	w.AddDirectivesEnterEventHandler(func(ctx *validation.Context, directives *ast.Directives) {
		if directives.Len() == 0 {
			return
		}

		knownDirectives := make(map[string]struct{}, directives.Len())

		directives.ForEach(func(directive ast.Directive, i int) {
			directiveName := directive.Name

			if _, ok := knownDirectives[directiveName]; ok {
				ctx.AddError(duplicateDirectiveMessage(directiveName, 0, 0))
			} else {
				knownDirectives[directiveName] = struct{}{}
			}
		})
	})
}

// duplicateDirectiveMessage ...
func duplicateDirectiveMessage(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"The directive \"" + directiveName + "\" can only be used once at this location.",
		// TODO: Location.
	)
}
