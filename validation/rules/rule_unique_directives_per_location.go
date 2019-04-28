package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// UniqueDirectivesPerLocation ...
func UniqueDirectivesPerLocation(w *validation.Walker) {
	w.AddDirectivesEnterEventHandler(func(ctx *validation.Context, directives *ast.Directives) {
		if directives.Len() == 0 {
			return
		}

		knownDirectives := make(map[string]struct{}, directives.Len())

		directives.ForEach(func(directive ast.Directive, i int) {
			directiveName := directive.Name

			if _, ok := knownDirectives[directiveName]; ok {
				ctx.AddError(validation.DuplicateDirectiveError(directiveName, 0, 0))
			} else {
				knownDirectives[directiveName] = struct{}{}
			}
		})
	})
}
