package rules

import (
	"github.com/bucketd/go-graphqlparser/validation"
)

// Specified is a slice of ValidationRuleFunc that contains all validation rules defined by the
// GraphQL specification.
//
// The order of this list is important, and is intended to produce the most clear output when
// encountering multiple validation errors.
var Specified = []validation.VisitFunc{
	ExecutableDefinitions,
	// UniqueOperationNames,
	LoneAnonymousOperation,
	// SingleFieldSubscriptions,
	KnownTypeNames,
	// FragmentsOnCompositeTypes,
	// VariablesAreInputTypes,
	// ScalarLeafs,Uni
	// FieldsOnCorrectType,
	// UniqueFragmentNames,
	// KnownFragmentNames,
	// NoUnusedFragments,
	// PossibleFragmentSpreads,
	// NoFragmentCycles,
	// UniqueVariableNames,
	// NoUndefinedVariables,
	NoUnusedVariables,
	KnownDirectives,
	UniqueDirectivesPerLocation,
	// KnownArgumentNames,
	UniqueArgumentNames,
	// ValuesOfCorrectType,
	// ProvidedRequiredArguments,
	// VariablesInAllowedPosition,
	// OverlappingFieldsCanBeMerged,
	UniqueInputFieldNames,
}

// SpecifiedSDL is a slice of ValidationRuleFunc that contains validation rules defined by the
// GraphQL specification for validating schema definition language documents exclusively. This set
// of rules is useful for servers that are parsing schemas, and other tools.
//
// The order of this list is important, and is intended to produce the most clear output when
// encountering multiple validation errors.
var SpecifiedSDL = []validation.VisitFunc{
	// These rules are handled before walking:
	//LoneSchemaDefinition,
	//UniqueDirectiveNames,
	//UniqueTypeNames,

	UniqueOperationTypes,
	UniqueDirectivesPerLocation,
	UniqueArgumentNames,
	UniqueInputFieldNames,
	KnownArgumentNamesOnDirectives,
	KnownDirectives,
	KnownTypeNames,
	PossibleNames,
	PossibleTypeExtensions,
	ProvidedRequiredArgumentsOnDirectives,

	// These rules are handled after walking:
	//UniqueEnumValueNames,
	//UniqueFieldDefinitionNames,
}
