package graphqlparser

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/bucketd/go-graphqlparser/validation/rules"
)

var (
	// DefaultValidationWalker ...
	DefaultValidationWalker = validation.NewWalker(rules.Specified)
	// DefaultValidationWalkerSDL ...
	DefaultValidationWalkerSDL = validation.NewWalker(rules.SpecifiedSDL)
)

// ParseDoc ...
func ParseDoc(doc []byte, schema *graphql.Schema) (*ast.Document, *graphql.Errors, error) {
	parser := language.NewParser(doc)

	queryAST, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	ctx := validation.Validate(queryAST, schema, DefaultValidationWalker)
	if ctx.Errors.Len() > 0 {
		return nil, ctx.Errors, nil
	}

	return &queryAST, nil, nil
}

// ParseSDLDoc ...
func ParseSDLDoc(doc []byte, schema *graphql.Schema) (*graphql.Schema, *graphql.Errors, error) {
	parser := language.NewParser(doc)

	sdlAST, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	ctx := validation.ValidateSDL(sdlAST, schema, DefaultValidationWalkerSDL)
	if ctx.Errors.Len() > 0 {
		return nil, ctx.Errors, nil
	}

	schema, err = validation.BuildSchema(ctx)
	if err != nil {
		return nil, nil, err
	}

	errs := validation.ValidateSchema(ctx, schema)
	if errs.Len() > 0 {
		return nil, errs, nil
	}

	return schema, nil, nil
}
