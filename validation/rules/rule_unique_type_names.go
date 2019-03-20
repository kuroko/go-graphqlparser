package rules

import (
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

// uniqueTypeNames ...
func uniqueTypeNames(w *validation.Walker) {

}

// duplicateTypeNameMessage ...
func duplicateTypeNameMessage(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// existedTypeNameMessage ...
func existedTypeNameMessage(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}
