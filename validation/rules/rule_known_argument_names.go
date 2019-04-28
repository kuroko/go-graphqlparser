package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/validation"
)

// KnownArgumentNames ...
func KnownArgumentNames(w *validation.Walker) {}

// KnownArgumentNamesOnDirectives ...
func KnownArgumentNamesOnDirectives(w *validation.Walker) {
	w.AddDirectiveEnterEventHandler(func(ctx *validation.Context, dir ast.Directive) {
		if dir.Arguments.Len() == 0 {
			return
		}

		directiveName := dir.Name
		// TODO: This probably won't handle built-in directives yet.
		directiveDef, _ := ctx.DirectiveDefinition(directiveName)

		if directiveDef == nil {
			// TODO: Should we do something else here? Or is this handled elsewhere?
			return
		}

		dirGen := dir.Arguments.Generator()
		defGen := directiveDef.ArgumentsDefinition.Generator()

		for dirArg, i := dirGen.Next(); i >= 0; dirArg, i = dirGen.Next() {
			var found bool
			for defArg, j := defGen.Next(); j >= 0; defArg, j = defGen.Next() {
				if dirArg.Name == defArg.Name {
					found = true
					break
				}
			}

			defGen.Reset()

			if !found {
				ctx.AddError(validation.UnknownDirectiveArgError(dirArg.Name, directiveName, 0, 0))
			}
		}
	})
}
