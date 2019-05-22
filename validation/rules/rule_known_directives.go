package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// KnownDirectives ...
func KnownDirectives(w *validation.Walker) {
	w.AddDirectiveEnterEventHandler(func(ctx *validation.Context, dir ast.Directive) {
		var def *ast.DirectiveDefinition
		var ok bool

		// First check if an existing schema exists, and the directive has already been defined. If
		// we're extending a schema, this will include the built-in directives.
		if ctx.Schema != nil && ctx.Schema.Directives != nil {
			def, ok = ctx.Schema.Directives[dir.Name]
		}

		// If we're in an SDL document, we need to consider directives that are defined.
		if !ok && ctx.SDLContext != nil {
			def, ok = ctx.SDLContext.DirectiveDefinitions[dir.Name]
		}

		// Lastly, check for built-in directives, these may have been overridden, which is why they
		// should be checked last.
		if !ok {
			def, ok = graphql.SpecifiedDirectives()[dir.Name]
		}

		if def == nil || !ok {
			ctx.AddError(validation.UnknownDirectiveError(dir.Name, 0, 0))
			return
		}

		// The directive definition doesn't contain the location this directive is currently being
		// used on it.
		if def.DirectiveLocations&dir.Location == 0 {
			ctx.AddError(validation.MisplacedDirectiveError(dir.Name, dir.Location, 0, 0))
		}
	})
}
