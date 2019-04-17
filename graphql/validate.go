package graphql

import (
	"errors"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

var (
	// DefaultValidationWalker ...
	DefaultValidationWalker = validation.NewWalker(rules.Specified)
	// DefaultValidationWalkerSDL ...
	DefaultValidationWalkerSDL = validation.NewWalker(rules.SpecifiedSDL)
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
func ValidateSDL(schema *types.Schema, doc []byte) (*types.Schema, *types.Errors, error) {
	parsed, err := Parse(doc)
	if err != nil {
		return nil, nil, errors.New("graphql: failed to parse document")
	}

	schema, errs := ValidateSDLAST(schema, parsed)

	return schema, errs, nil
}

// ValidateAST ...
func ValidateAST(schema *types.Schema, doc ast.Document) *types.Errors {
	// TODO: Re-use existing walker, stored in global.
	return validation.Validate(doc, schema, DefaultValidationWalker)
}

// ValidateSDLAST ...
func ValidateSDLAST(schema *types.Schema, doc ast.Document) (*types.Schema, *types.Errors) {
	// TODO: Re-use existing walker, stored in global.
	return validation.ValidateSDL(doc, schema, DefaultValidationWalkerSDL)
}
