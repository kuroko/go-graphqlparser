package graphql

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// BuildASTSchema ...
func BuildASTSchema(schema *types.Schema, doc ast.Document) (*types.Schema, *types.Errors) {
	errs := ValidateSDLAST(nil, doc)
	if errs.Len() > 0 {
		return nil, errs
	}

	if schema == nil {
		schema = &types.Schema{}
	}

	schemaVisitFns := []validation.VisitFunc{
		setSchemaTypeDefinitions,
	}

	// Traverse the schema AST, populating the schema object with relevant information.
	validation.NewWalker(schemaVisitFns).Walk(&validation.Context{Schema: schema}, doc)

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

// setSchemaTypeDefinitions ...
func setSchemaTypeDefinitions(w *validation.Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		if ctx.Schema.Types == nil {
			ctx.Schema.Types = make(map[string]*ast.TypeDefinition)
		}

		ctx.Schema.Types[def.Name] = def
	})
}
