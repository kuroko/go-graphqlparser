package graphql

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
	"github.com/bucketd/go-graphqlparser/validation"
)

// BuildASTSchema ...
func BuildASTSchema(schema *types.Schema, doc ast.Document) (*types.Schema, *types.Errors) {
	_, errs := ValidateSDLAST(nil, doc)
	if errs.Len() > 0 {
		return nil, errs
	}

	if schema == nil {
		// We only add in the built-in types if we're not extending a schema, we assume an existing
		// schema will have already gone through this process. If not, it's probably the developers
		// fault, and they'll probably encounter a panic.
		schema = &types.Schema{
			Directives: types.SpecifiedDirectives(),
			Types:      make(map[string]*ast.TypeDefinition),
		}

		// TODO: Implement the other built-in scalars, and move into funcs.
		schema.Types["ID"] = &ast.TypeDefinition{
			Name: "ID",
			Kind: ast.TypeDefinitionKindScalar,
		}
	}

	// Build symbol table entries from the schema document.
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
		ctx.Schema.Directives[def.Name] = def
	})
}

// setSchemaTypeDefinitions ...
func setSchemaTypeDefinitions(w *validation.Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *validation.Context, def *ast.TypeDefinition) {
		ctx.Schema.Types[def.Name] = def
	})
}
