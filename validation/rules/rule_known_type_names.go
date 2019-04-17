package rules

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// KnownTypeNames ...
func KnownTypeNames(w *validation.Walker) {
	w.AddNamedTypeEnterEventHandler(func(ctx *validation.Context, t ast.Type) {
		typeName := t.NamedType

		var existsInSchema bool
		var existsInDocument bool

		_, existsInSchema = ctx.Schema.Types[typeName]

		// If we're validating an SDL document, types declared in the current document must also be
		// taken into account.
		if ctx.SDLContext != nil {
			_, existsInDocument = ctx.SDLContext.TypeDefinitions[typeName]
		}

		if !existsInSchema && !existsInDocument && !isSpecifiedScalarName(typeName) {
			ctx.AddError(UnknownTypeError(typeName, 0, 0))
		}
	})
}

// UnknownTypeError ...
func UnknownTypeError(typeName string, line, col int) types.Error {
	return types.NewError(
		"Unknown type \"" + typeName + "\".",
		// TODO: Location.
	)
}

// isSpecifiedScalarName ...
//
// NOTE: This function is needed for when a schema is being parsed, and we aren't extending a
// schema. At this point we haven't "registered" the built-in scalar types.
func isSpecifiedScalarName(typeName string) bool {
	switch typeName {
	case "String", "Int", "Float", "Boolean", "ID":
		return true
	}

	return false
}
