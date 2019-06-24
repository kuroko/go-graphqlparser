package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// ValidateSchema ...
func ValidateSchema(ctx *Context, schema *graphql.Schema) *graphql.Errors {
	validateRootTypes(ctx, schema)
	validateDirectives(ctx, schema)
	validateTypes(ctx, schema)

	if ctx.Errors.Len() > 0 {
		return ctx.Errors
	}

	return nil
}

// validateRootTypes ...
func validateRootTypes(ctx *Context, schema *graphql.Schema) {
	if schema.QueryType == nil {
		ctx.AddError(graphql.NewError(
			"Query root type must be provided.",
			// TODO: Location.
		))
	} else {
		// Query has a special case compared to Mutation and Subscription types in that the Query
		// type may automatically be assigned if not SchemaDefinition is provided. In that case,
		// it's possible for the Query type itself to not exist, despite QueryType being set on the
		// schema itself, so we must also check that.

		queryType, ok := schema.Types[schema.QueryType.NamedType]
		if ok && !ast.IsObjectTypeDefinition(queryType) {
			ctx.AddError(graphql.NewError(
				"Query root type must be Object type.",
				// TODO: Location.
			))
		} else if !ok {
			ctx.AddError(graphql.NewError(
				"Query root type must be provided.",
				// TODO: Location.
			))
		}
	}

	if schema.MutationType != nil {
		mutationType, ok := schema.Types[schema.MutationType.NamedType]
		if ok && !ast.IsObjectTypeDefinition(mutationType) {
			ctx.AddError(graphql.NewError(
				"Mutation root type must be Object type if provided.",
				// TODO: Location.
			))
		}
	}

	if schema.SubscriptionType != nil {
		subscriptionType, ok := schema.Types[schema.SubscriptionType.NamedType]
		if ok && !ast.IsObjectTypeDefinition(subscriptionType) {
			ctx.AddError(graphql.NewError(
				"Subscription root type must be Object type if provided.",
				// TODO: Location.
			))
		}
	}
}

// validateDirectives ...
func validateDirectives(ctx *Context, schema *graphql.Schema) {
	for directiveName, directiveDef := range schema.Directives {
		validateName(ctx, directiveName, 0, 0) // TODO: Location

		argNames := make(map[string]struct{})
		directiveDef.ArgumentsDefinition.ForEach(func(ivd ast.InputValueDefinition, i int) {
			if _, ok := argNames[ivd.Name]; ok {
				ctx.AddError(graphql.NewError(
					"Argument @" + directiveName + "(" + ivd.Name + ":) can only be defined once.",
					// TODO: Location.
				))
				return
			}

			validateName(ctx, ivd.Name, 0, 0) // TODO: Location

			argNames[ivd.Name] = struct{}{}
		})
	}
}

// validateTypes ...
func validateTypes(ctx *Context, schema *graphql.Schema) {
	for typeName, typeDef := range schema.Types {
		// If the name exactly matches one of the built-in introspection type names, don't bother
		// validating any further.
		if !isIntrospectionTypeName(typeName) {
			validateName(ctx, typeName, 0, 0) // TODO: Location
		}

		switch {
		case ast.IsObjectTypeDefinition(typeDef):
			validateFields(ctx, schema, typeDef)
			// TODO: ...
		case ast.IsInterfaceTypeDefinition(typeDef):
			validateFields(ctx, schema, typeDef)
			// TODO: ...
		case ast.IsUnionTypeDefinition(typeDef):
			// TODO: ...
		case ast.IsEnumTypeDefinition(typeDef):
			// TODO: ...
		case ast.IsInputObjectTypeDefinition(typeDef):
			// TODO: ...
		}
	}
}

// validateFields ...
func validateFields(ctx *Context, schema *graphql.Schema, typeDef *ast.TypeDefinition) {
	if typeDef.FieldsDefinition.Len() == 0 {
		ctx.AddError(graphql.NewError(
			"Type " + typeDef.Name + " must define one or more fields.",
			// TODO: Location.
		))
	}

	typeDef.FieldsDefinition.ForEach(func(field ast.FieldDefinition, i int) {
		validateName(ctx, field.Name, 0, 0) // TODO: Location.

		if !IsOutputType(schema, field.Type) {
			ctx.AddError(graphql.NewError(
				"The type of " + typeDef.Name + "." + field.Name + " must be Output Type.",
				// TODO: Location.
			))
		}

		argNames := make(map[string]struct{}, field.ArgumentsDefinition.Len())

		field.ArgumentsDefinition.ForEach(func(ivd ast.InputValueDefinition, i int) {
			validateName(ctx, ivd.Name, 0, 0) // TODO: Location.

			if _, ok := argNames[ivd.Name]; ok {
				ctx.AddError(graphql.NewError(
					"Field argument " + typeDef.Name + "." + field.Name + "(" + ivd.Name + ":) can only be defined " +
						"once.",
					// TODO: Location.
				))
			}

			argNames[ivd.Name] = struct{}{}

			if !IsInputType(schema, ivd.Type) {
				ctx.AddError(graphql.NewError(
					"The type of " + typeDef.Name + "." + field.Name + "(" + ivd.Name + ":) must be Input Type.",
					// TODO: Location.
				))
			}
		})
	})
}

// validateName ...
func validateName(ctx *Context, name string, line, col int) {
	nameLen := len(name)

	if nameLen == 0 {
		ctx.AddError(InvalidNameError(name, line, col))
		return
	}

	if len(name) > 1 && name[0] == '_' && name[1] == '_' {
		ctx.AddError(ReservedNameError(name, line, col))
		return
	}

	for i, r := range name {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_':
		case i > 0 && (r >= '0' && r <= '9'):
		default:
			ctx.AddError(InvalidNameError(name, line, col))
			return
		}
	}
}

// isIntrospectionTypeName ...
func isIntrospectionTypeName(typeName string) bool {
	switch typeName {
	case "__Schema":
		return true
	case "__Directive":
		return true
	case "__DirectiveLocation":
		return true
	case "__Type":
		return true
	case "__Field":
		return true
	case "__InputValue":
		return true
	case "__EnumValue":
		return true
	case "__TypeKind":
		return true
	default:
		return false
	}
}
