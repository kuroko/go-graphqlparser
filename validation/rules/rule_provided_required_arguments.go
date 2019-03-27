package rules

import (
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// ProvidedRequiredArguments ...
func ProvidedRequiredArguments(w *validation.Walker) {}

// MissingFieldArgMessage ...
func MissingFieldArgMessage(fieldName, argName, t string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + fieldName + "\" argument \"" + argName + "\" of type \"" + t + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// MissingDirectiveArgMessage ...
func MissingDirectiveArgMessage(directiveName, argName, t string, line, col int) types.Error {
	return types.NewError(
		"Directive \"" + directiveName + "\" argument \"" + argName + "\" of type \"" + t + "\" is required, but it was not provided",
		// TODO: Location.
	)
}
