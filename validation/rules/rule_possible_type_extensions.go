package rules

import (
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// PossibleTypeExtensions ...
func PossibleTypeExtensions(w *validation.Walker) {}

// ExtendingUnknownTypeError ...
func ExtendingUnknownTypeError(typeName string, line, col int) types.Error {
	return types.NewError("Cannot extend type \"" + typeName + "\" because it is not defined.")
}

// ExtendingDifferentTypeKindError ...
func ExtendingDifferentTypeKindError(typeName, kind string, line, col int) types.Error {
	return types.NewError("Cannot extend non-" + kind + "type \"" + typeName + "\".")
}
