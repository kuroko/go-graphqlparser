package graphql

import (
	"fmt"

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

	schema.Directives = make(map[string]*ast.DirectiveDefinition)
	schema.Types = make(map[string]*ast.TypeDefinition)

	// Add built-in directives, required by spec.
	schema.Directives["skip"] = buildGraphQLSkipDirective()
	schema.Directives["include"] = buildGraphQLIncludeDirective()
	schema.Directives["deprecated"] = buildGraphQLDeprecatedDirective()

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

// buildGraphQLSkipDirective ...
func buildGraphQLSkipDirective() *ast.DirectiveDefinition {
	return &ast.DirectiveDefinition{
		Name:        "skip",
		Description: "Directs the executor to skip this field or fragment when the `if` argument is true.",
		DirectiveLocations: ast.DirectiveLocationKindField |
			ast.DirectiveLocationKindFragmentSpread |
			ast.DirectiveLocationKindInlineFragment,
		ArgumentsDefinition: (*ast.InputValueDefinitions)(nil).
			Add(ast.InputValueDefinition{
				Name:        "if",
				Description: "Skipped when true.",
				Type: ast.Type{
					NamedType:   "Boolean",
					Kind:        ast.TypeKindNamed,
					NonNullable: true,
				},
			}),
	}
}

// buildGraphQLIncludeDirective ...
func buildGraphQLIncludeDirective() *ast.DirectiveDefinition {
	return &ast.DirectiveDefinition{
		Name:        "include",
		Description: "Directs the executor to include this field or fragment only when the `if` argument is true.",
		DirectiveLocations: ast.DirectiveLocationKindField |
			ast.DirectiveLocationKindFragmentSpread |
			ast.DirectiveLocationKindInlineFragment,
		ArgumentsDefinition: (*ast.InputValueDefinitions)(nil).
			Add(ast.InputValueDefinition{
				Name:        "if",
				Description: "Included when true.",
				Type: ast.Type{
					NamedType:   "Boolean",
					Kind:        ast.TypeKindNamed,
					NonNullable: true,
				},
			}),
	}
}

// buildGraphQLDeprecatedDirective ...
func buildGraphQLDeprecatedDirective() *ast.DirectiveDefinition {
	return &ast.DirectiveDefinition{
		Name:        "deprecated",
		Description: "Marks an element of a GraphQL schema as no longer supported.",
		DirectiveLocations: ast.DirectiveLocationKindFieldDefinition |
			ast.DirectiveLocationKindEnumValue,
		ArgumentsDefinition: (*ast.InputValueDefinitions)(nil).
			Add(ast.InputValueDefinition{
				Name: "reason",
				Description: "Explains why this element was deprecated, usually also including a " +
					"suggestion for how to access supported similar data. Formatted using " +
					"the Markdown syntax (as specified by [CommonMark](https://commonmark.org/).",
				Type: ast.Type{
					NamedType: "String",
					Kind:      ast.TypeKindNamed,
				},
				DefaultValue: &ast.Value{
					StringValue: "No longer supported",
					Kind:        ast.ValueKindString,
				},
			}),
	}
}
