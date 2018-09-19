package ast

import (
	"go/scanner"
	"go/token"
)

type validationRule func() error

func (d Document) ValidateExecutable(document Document) error {

	var err scanner.ErrorList

	// Two lists of rules, one for TypeSystem one for Executable?
	var validationRules = []validationRule{
		d.executableDefinitions,
		d.fieldsOnCorrectType,
		d.fragmentsOnCompositeTypes,
		d.knownArgumentNames,
		d.knownDirectives,
		d.knownFragmentNames,
		d.knownTypeNames,
		d.loneAnonymousOperation,
		d.loneSchemaDefinition,
		d.noFragmentCycles,
		d.noUndefinedVariables,
		d.noUnusedFragments,
		d.noUnusedVariables,
		d.overlappingFieldsCanBeMerged,
		d.possibleFragmentSpreads,
		d.providedRequiredArguments,
		d.scalarLeafs,
		d.singleFieldSubscriptions,
		d.uniqueArgumentNames,
		d.uniqueDirectivesPerLocation,
		d.uniqueFragmentNames,
		d.uniqueInputFieldNames,
		d.uniqueOperationNames,
		d.uniqueVariableNames,
		d.valuesOfCorrectType,
		d.variablesAreInputTypes,
		d.variablesInAllowedPosition,
	}

	for _, vr := range validationRules {
		errString := vr().Error()
		if len(errString) != 0 {
			err.Add(token.Position{}, errString)
		}
	}

	return err.Err()
}

func (d Document) ValidateTypeSystem(document Document) error {
	return nil
}

// Validation rules
// From: https://github.com/graphql/graphql-js/tree/8682f57095b38742ef2aae24915d5f26be8ab97d/src/validation/rules

func (d Document) executableDefinitions() error {

	var err scanner.ErrorList

	d.Definitions.ForEach(func(definition Definition, _ int) {
		if definition.Kind != DefinitionKindExecutable {
			pos := token.Position{
				// Filename: ??,
				// Offset: ??,
				// Line: ??,
				// Column: ??,
			}

			// TODO: better error message, with token.Position context.
			err.Add(pos, "found non-executable definiton - this needs context")
		}
	})

	return err.Err()
}

func (d Document) fieldsOnCorrectType() error          { return nil }
func (d Document) fragmentsOnCompositeTypes() error    { return nil }
func (d Document) knownArgumentNames() error           { return nil }
func (d Document) knownDirectives() error              { return nil }
func (d Document) knownFragmentNames() error           { return nil }
func (d Document) knownTypeNames() error               { return nil }
func (d Document) loneAnonymousOperation() error       { return nil }
func (d Document) loneSchemaDefinition() error         { return nil }
func (d Document) noFragmentCycles() error             { return nil }
func (d Document) noUndefinedVariables() error         { return nil }
func (d Document) noUnusedFragments() error            { return nil }
func (d Document) noUnusedVariables() error            { return nil }
func (d Document) overlappingFieldsCanBeMerged() error { return nil }
func (d Document) possibleFragmentSpreads() error      { return nil }
func (d Document) providedRequiredArguments() error    { return nil }
func (d Document) scalarLeafs() error                  { return nil }
func (d Document) singleFieldSubscriptions() error     { return nil }
func (d Document) uniqueArgumentNames() error          { return nil }
func (d Document) uniqueDirectivesPerLocation() error  { return nil }
func (d Document) uniqueFragmentNames() error          { return nil }
func (d Document) uniqueInputFieldNames() error        { return nil }
func (d Document) uniqueOperationNames() error         { return nil }
func (d Document) uniqueVariableNames() error          { return nil }
func (d Document) valuesOfCorrectType() error          { return nil }
func (d Document) variablesAreInputTypes() error       { return nil }
func (d Document) variablesInAllowedPosition() error   { return nil }
