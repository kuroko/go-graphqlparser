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
		setSchemaOperationTypes,
		setSchemaDirectiveDefinitions,
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

// setSchemaOperationTypes ...
func setSchemaOperationTypes(w *validation.Walker) {
	// NOTE: This handles both schema definitions and extensions.
	w.AddOperationTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def ast.OperationTypeDefinition) {
		switch def.OperationType {
		case ast.OperationDefinitionKindQuery:
			ctx.Schema.QueryType = &def.NamedType
		case ast.OperationDefinitionKindMutation:
			ctx.Schema.MutationType = &def.NamedType
		case ast.OperationDefinitionKindSubscription:
			ctx.Schema.SubscriptionType = &def.NamedType
		}
	})
}

// setSchemaDirectiveDefinitions ...
func setSchemaDirectiveDefinitions(w *validation.Walker) {
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.DirectiveDefinition) {
		if ctx.Schema.Directives == nil {
			ctx.Schema.Directives = make(map[string]*ast.DirectiveDefinition)
		}

		ctx.Schema.Directives[def.Name] = def
	})
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
