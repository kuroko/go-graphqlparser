package graphql

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
)

// BuildASTSchema ...
func BuildASTSchema(schema *types.Schema, doc ast.Document) (*types.Schema, *types.Errors) {
	_, errs := ValidateSDLAST(nil, doc)
	if errs.Len() > 0 {
		return nil, errs
	}

	return schema, nil
}

// BuildSchema ...
func BuildSchema(schema *types.Schema, doc []byte) (*types.Schema, *types.Errors, error) {
	schemaDoc, err := Parse(doc)
	if err != nil {
		// TODO: Error wrapping. Maybe some kind of context?
		return nil, nil, err
	}

	schema, errs := BuildASTSchema(schema, schemaDoc)

	return schema, errs, nil
}

// MustBuildSchema ...
func MustBuildSchema(schema *types.Schema, doc []byte) *types.Schema {
	schema, errs, err := BuildSchema(schema, doc)
	if err != nil {
		panic(fmt.Sprintf("graphql: error building schema: %v", err))
	}

	if errs.Len() > 0 {
		panic(fmt.Sprintf("graphql: validation failed whilst building schema: %v", errs))
	}

	return schema
}
