package graphql

import (
	"errors"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

// Validate ...
func Validate(schema *types.Schema, doc []byte) (*types.Errors, error) {
	parsed, err := Parse(doc)
	if err != nil {
		return nil, errors.New("graphql: failed to parse document")
	}

	return ValidateAST(schema, parsed), nil
}

// ValidateSDL ...
func ValidateSDL(schema *types.Schema, doc []byte) (*types.Errors, error) {
	parsed, err := Parse(doc)
	if err != nil {
		return nil, errors.New("graphql: failed to parse document")
	}

	return ValidateSDLAST(schema, parsed), nil
}

// ValidateAST ...
func ValidateAST(schema *types.Schema, doc ast.Document) *types.Errors {
	return validation.Validate(doc, schema, validation.NewWalker(rules.Specified))
}

// ValidateSDLAST ...
func ValidateSDLAST(schema *types.Schema, doc ast.Document) *types.Errors {
	return validation.ValidateSDL(doc, schema, validation.NewWalker(rules.SpecifiedSDL))
}
