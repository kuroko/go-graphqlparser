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

		_, existsInSchema := ctx.Schema.Types[typeName]
		_, existsInDocument := ctx.TypeDefinitions()[typeName]

		if !existsInSchema && !existsInDocument && !isSpecifiedScalarName(typeName) {
			ctx.AddError(UnknownTypeError(typeName, []string{}, 0, 0))
		}
	})
}

// UnknownTypeMessage ...
func UnknownTypeError(typeName string, suggestedTypes []string, line, col int) types.Error {
	// TODO: Implement this kind of logic, ish.
	//if len(suggestedTypes) > 0 {
	//	message += ` Did you mean ${quotedOrList(suggestedTypes)}?`;
	//}

	return types.NewError(
		"Unknown type \"" + typeName + "\".",
		// TODO: Location.
	)
}

// isSpecifiedScalarName ...
func isSpecifiedScalarName(typeName string) bool {
	switch typeName {
	case "String", "Int", "Float", "Boolean", "ID":
		return true
	}

	return false
}
