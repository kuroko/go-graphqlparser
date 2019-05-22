package validation

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// AnonOperationNotAloneError ...
func AnonOperationNotAloneError(line, col int) graphql.Error {
	return graphql.NewError(
		"This anonymous operation must be the only defined operation.",
		// TODO: Location.
	)
}

// CanNotDefineSchemaWithinExtensionError ...
func CanNotDefineSchemaWithinExtensionError(line, col int) graphql.Error {
	return graphql.NewError(
		"Cannot define a new schema within a schema extension.",
		// TODO: Location.
	)
}

// DuplicateArgError ...
func DuplicateArgError(argName string, line, col int) graphql.Error {
	return graphql.NewError("There can be only one argument named \"" + argName + "\".")
}

// DuplicateDirectiveError ...
func DuplicateDirectiveError(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"The directive \"" + directiveName + "\" can only be used once at this location.",
		// TODO: Location.
	)
}

// DuplicateDirectiveNameError ...
func DuplicateDirectiveNameError(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one directive named \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// DuplicateEnumValueNameError ...
func DuplicateEnumValueNameError(typeName, valueName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Enum value \"" + typeName + "." + valueName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// DuplicateFieldDefinitionNameError ...
func DuplicateFieldDefinitionNameError(typeName, fieldName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Field \"" + typeName + "." + fieldName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// DuplicateInputFieldError ...
func DuplicateInputFieldError(fieldName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one input field named \"" + fieldName + "\".",
		// TODO: Location.
	)
}

// DuplicateOperationNameError ...
func DuplicateOperationNameError(operationName string, line, col int) graphql.Error {
	return graphql.NewError("There can be only one operation named " + operationName + ".")
}

// DuplicateOperationTypeError ...
func DuplicateOperationTypeError(operation string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one " + operation + " type in schema.",
		// TODO: Location.
	)
}

// DuplicateTypeNameError ...
func DuplicateTypeNameError(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// ExistedDirectiveNameError ...
func ExistedDirectiveNameError(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Directive \"" + directiveName + "\" already exists in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}

// ExistedOperationTypeError ...
func ExistedOperationTypeError(operation string, line, col int) graphql.Error {
	return graphql.NewError(
		"Type for " + operation + " already defined in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}

// ExistedTypeNameError ...
func ExistedTypeNameError(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}

// ExtendingDifferentTypeKindError ...
func ExtendingDifferentTypeKindError(typeName, kind string, line, col int) graphql.Error {
	return graphql.NewError("Cannot extend non-" + kind + "type \"" + typeName + "\".")
}

// ExtendingUnknownTypeError ...
func ExtendingUnknownTypeError(typeName string, line, col int) graphql.Error {
	return graphql.NewError("Cannot extend type \"" + typeName + "\" because it is not defined.")
}

// MisplacedDirectiveError ...
func MisplacedDirectiveError(directiveName string, location ast.DirectiveLocation, line, col int) graphql.Error {
	return graphql.NewError(
		"Directive \"" + directiveName + "\" may not be used on " + ast.NamesByDirectiveLocations[location] + ".",
		// TODO: Location.
	)
}

// MissingDirectiveArgError ...
func MissingDirectiveArgError(directiveName, argName, typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Directive \"" + directiveName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// MissingFieldArgError ...
func MissingFieldArgError(fieldName, argName, typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Field \"" + fieldName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// NameStartsWithTwoUnderscoresError ...
func NameStartsWithTwoUnderscoresError(name string, line, col int) graphql.Error {
	return graphql.NewError(
		"Name \"" + name + "\" must not begin with \"__\" (two underscores)",
		// TODO: Location.
	)
}

// NonExecutableDefinitionError ...
func NonExecutableDefinitionError(name string, line, col int) graphql.Error {
	return graphql.NewError(
		"The \"" + name + "\" definition is not executable.",
		// TODO: Location.
	)
}

// SchemaDefinitionNotAloneError ...
func SchemaDefinitionNotAloneError(line, col int) graphql.Error {
	return graphql.NewError(
		"Must provide only one schema definition.",
		// TODO: Location.
	)
}

// UnionHasNoMembersError ...
func UnionHasNoMembersError(unionName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Union \"" + unionName + "\" must have at least one member.",
		// TODO: Location.
	)
}

// UnknownDirectiveArgError ...
func UnknownDirectiveArgError(argName, directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Unknown argument \"" + argName + "\" on directive \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// UnknownDirectiveError ...
func UnknownDirectiveError(directiveName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Unknown directive \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// UnknownTypeError ...
func UnknownTypeError(typeName string, line, col int) graphql.Error {
	return graphql.NewError(
		"Unknown type \"" + typeName + "\".",
		// TODO: Location.
	)
}

// UnusedVariableError ...
func UnusedVariableError(varName, opName string, line, col int) graphql.Error {
	return graphql.NewError(
		unusedVariableMessage(varName, opName),
		// TODO: Location.
	)
}

// unusedVariableMessage ...
func unusedVariableMessage(varName, opName string) string {
	if len(opName) > 0 {
		return fmt.Sprintf("Variable %s is never used in operation %s", varName, opName)
	}

	return fmt.Sprintf("Variable %s is never used", varName)
}
