package rules_test

import (
	"github.com/bucketd/go-graphqlparser"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/language"
	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/davecgh/go-spew/spew"
)

// mustBuildSchema ...
func mustBuildSchema(schema *graphql.Schema, doc []byte) *graphql.Schema {
	schema, errs, err := buildSchema(schema, doc)
	if err != nil {
		panic(err)
	}

	if errs.Len() > 0 {
		panic(spew.Sdump(errs))
	}

	return schema
}

// buildSchema ...
func buildSchema(schema *graphql.Schema, doc []byte) (*graphql.Schema, *graphql.Errors, error) {
	parser := language.NewParser(doc)

	sdlAST, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	ctx := validation.ValidateSDL(sdlAST, schema, graphqlparser.DefaultValidationWalkerSDL)
	if ctx.Errors.Len() > 0 {
		return nil, ctx.Errors, nil
	}

	schema, err = validation.BuildSchema(ctx)
	if err != nil {
		return nil, nil, err
	}

	return schema, nil, nil
}
