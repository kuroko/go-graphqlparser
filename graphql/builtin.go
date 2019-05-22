package graphql

import "github.com/bucketd/go-graphqlparser/ast"

var (
	// SkipDirective is the definition for the built-in "@skip" directive.
	SkipDirective = &ast.DirectiveDefinition{
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
	// IncludeDirective is the definition for the built-in "@include" directive.
	IncludeDirective = &ast.DirectiveDefinition{
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
	// DeprecatedDirective is the definition for the built-in "@deprecated" directive.
	DeprecatedDirective = &ast.DirectiveDefinition{
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

	// IDType ...
	// TODO: Fully implement.
	IDType = &ast.TypeDefinition{
		Name: "ID",
		Kind: ast.TypeDefinitionKindScalar,
	}

	// TODO: Implement all built-in scalar types
)

// SpecifiedDirectives returns a map similar to the one found on the Schema type, containing all
// pre-defined GraphQL directives. It is returned from a function to avoid this map being mutated,
// as it is used in several places.
func SpecifiedDirectives() map[string]*ast.DirectiveDefinition {
	return map[string]*ast.DirectiveDefinition{
		"skip":       SkipDirective,
		"include":    IncludeDirective,
		"deprecated": DeprecatedDirective,
	}
}

// SpecifiedTypes ...
func SpecifiedTypes() map[string]*ast.TypeDefinition {
	return map[string]*ast.TypeDefinition{
		"ID": IDType,
	}
}
