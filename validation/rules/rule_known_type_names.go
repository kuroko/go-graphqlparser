package rules

import (
	"sort"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// KnownTypeNames ...
func KnownTypeNames(w *validation.Walker) {
	w.AddNamedTypeEnterEventHandler(func(ctx *validation.Context, t ast.Type) {
		typeName := t.NamedType

		_, existsInSchema := ctx.Schema.Types[typeName]
		_, existsInDocument := ctx.TypeDefinitions()[typeName]

		if !existsInSchema && !existsInDocument && !isSpecifiedScalarName(typeName) {
			// TODO: Optimise, we can have a type names slice on the Schema at least.
			// TODO: Move into own function?
			// TODO: Need to include built-in scalar types if we're not extending, and it's an SDL
			// document. The test for both of those things is simple at least.
			typeNames := make([]string, 0, len(ctx.Schema.Types)+len(ctx.TypeDefinitions()))

			for k := range ctx.Schema.Types {
				typeNames = append(typeNames, k)
			}

			for k := range ctx.TypeDefinitions() {
				typeNames = append(typeNames, k)
			}

			if ctx.SDLContext != nil && ctx.SDLContext.IsExtending == false {
				typeNames = append(typeNames, []string{"String", "Int", "Float", "Boolean", "ID"}...)
			}

			// This is important.
			sort.Strings(typeNames)

			ctx.AddError(UnknownTypeError(typeName, validation.SuggestionList(typeName, typeNames), 0, 0))
		}
	})
}

// UnknownTypeMessage ...
func UnknownTypeError(typeName string, suggestedTypes []string, line, col int) types.Error {
	message := "Unknown type \"" + typeName + "\"."
	if len(suggestedTypes) > 0 {
		message += " Did you mean " + validation.QuotedOrList(suggestedTypes) + "?"
	}

	return types.NewError(
		message,
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
