package validation

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
)

// DuplicateTypeNameError ...
func DuplicateTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one type named " + typeName + ".",
		// TODO: Location.
	)
}

// ExistedTypeNameError ...
func ExistedTypeNameError(typeName string, line, col int) types.Error {
	return types.NewError(
		"Type " + typeName + " already exists in the schema. It cannot also be defined in this type definition.",
		// TODO: Location.
	)
}

// DuplicateDirectiveNameError ...
func DuplicateDirectiveNameError(directiveName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one directive named \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// ExistedDirectiveNameError ...
func ExistedDirectiveNameError(directiveName string, line, col int) types.Error {
	return types.NewError(
		"Directive \"" + directiveName + "\" already exists in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}

// SchemaDefinitionNotAloneError ...
func SchemaDefinitionNotAloneError(line, col int) types.Error {
	return types.NewError(
		"Must provide only one schema definition.",
		// TODO: Location.
	)
}

// CanNotDefineSchemaWithinExtensionError ...
func CanNotDefineSchemaWithinExtensionError(line, col int) types.Error {
	return types.NewError(
		"Cannot define a new schema within a schema extension.",
		// TODO: Location.
	)
}

// DuplicateFieldDefinitionNameError ...
func DuplicateFieldDefinitionNameError(typeName, fieldName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + typeName + "." + fieldName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// DuplicateEnumValueNameError ...
func DuplicateEnumValueNameError(typeName, valueName string, line, col int) types.Error {
	return types.NewError(
		"Enum value \"" + typeName + "." + valueName + "\" can only be defined once.",
		// TODO: Location.
	)
}

// DuplicateOperationTypeError ...
func DuplicateOperationTypeError(operation string, line, col int) types.Error {
	return types.NewError(
		"There can be only one " + operation + " type in schema.",
		// TODO: Location.
	)
}

// ExistedOperationTypeError ...
func ExistedOperationTypeError(operation string, line, col int) types.Error {
	return types.NewError(
		"Type for " + operation + " already defined in the schema. It cannot be redefined.",
		// TODO: Location.
	)
}

// DuplicateOperationNameError ...
func DuplicateOperationNameError(operationName string, line, col int) types.Error {
	return types.NewError("There can be only one operation named " + operationName + ".")
}

// DuplicateInputFieldError ...
func DuplicateInputFieldError(fieldName string, line, col int) types.Error {
	return types.NewError(
		"There can be only one input field named \"" + fieldName + "\".",
		// TODO: Location.
	)
}

// DuplicateDirectiveError ...
func DuplicateDirectiveError(directiveName string, line, col int) types.Error {
	return types.NewError(
		"The directive \"" + directiveName + "\" can only be used once at this location.",
		// TODO: Location.
	)
}

// DuplicateArgError ...
func DuplicateArgError(argName string, line, col int) types.Error {
	return types.NewError("There can be only one argument named \"" + argName + "\".")
}

// MissingFieldArgError ...
func MissingFieldArgError(fieldName, argName, typeName string, line, col int) types.Error {
	return types.NewError(
		"Field \"" + fieldName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// MissingDirectiveArgError ...
func MissingDirectiveArgError(directiveName, argName, typeName string, line, col int) types.Error {
	return types.NewError(
		"Directive \"" + directiveName + "\" argument \"" + argName + "\" of type \"" + typeName + "\" is required, but it was not provided",
		// TODO: Location.
	)
}

// ExtendingUnknownTypeError ...
func ExtendingUnknownTypeError(typeName string, line, col int) types.Error {
	return types.NewError("Cannot extend type \"" + typeName + "\" because it is not defined.")
}

// ExtendingDifferentTypeKindError ...
func ExtendingDifferentTypeKindError(typeName, kind string, line, col int) types.Error {
	return types.NewError("Cannot extend non-" + kind + "type \"" + typeName + "\".")
}

// UnusedVariableError ...
func UnusedVariableError(varName, opName string, line, col int) types.Error {
	return types.NewError(
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

// AnonOperationNotAloneError ...
func AnonOperationNotAloneError(line, col int) types.Error {
	return types.NewError(
		"This anonymous operation must be the only defined operation.",
		// TODO: Location.
	)
}

// UnknownTypeError ...
func UnknownTypeError(typeName string, line, col int) types.Error {
	return types.NewError(
		"Unknown type \"" + typeName + "\".",
		// TODO: Location.
	)
}

// UnknownDirectiveError ...
func UnknownDirectiveError(directiveName string, line, col int) types.Error {
	return types.NewError(
		"Unknown directive \"" + directiveName + "\".",
		// TODO: Location.
	)
}

// MisplacedDirectiveError ...
func MisplacedDirectiveError(directiveName string, location ast.DirectiveLocation, line, col int) types.Error {
	return types.NewError(
		"Directive \"" + directiveName + "\" may not be used on " + ast.NamesByDirectiveLocations[location] + ".",
		// TODO: Location.
	)
}

// NonExecutableDefinitionError ...
func NonExecutableDefinitionError(name string, line, col int) types.Error {
	return types.NewError(
		"The \"" + name + "\" definition is not executable.",
		// TODO: Location.
	)
}
