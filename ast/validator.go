package ast

import (
	"errors"
)

type Validator struct {
	DocumentAST   *Document
	GraphQLSchema *Document
}

func NewValidator(documentAST *Document, gqlSchema *Document) *Validator {
	return &Validator{
		DocumentAST:   documentAST,
		GraphQLSchema: gqlSchema,
	}
}

func (v *Validator) ApplyExecutableRules() *Errors {
	var errs *Errors

	for _, erFn := range executableRules {
		e := erFn(v).Reverse()

		if e == nil {
			continue
		}

		e.ForEach(func(err Error, _ int) {
			errs = errs.Add(err)
		})
	}

	return errs.Reverse()
}

func (v *Validator) ApplySDLRules() *Errors { return nil }

var executableRules = []func(*Validator) *Errors{
	(*Validator).executableDefinitions,
	(*Validator).uniqueOperationNames,
	(*Validator).loneAnonymousOperation,
	(*Validator).singleFieldSubscriptions,
	(*Validator).knownTypeNames,
	(*Validator).fragmentsOnCompositeTypes,
	(*Validator).variablesAreInputTypes,
	(*Validator).scalarLeafs,
	(*Validator).fieldsOnCorrectType,
	(*Validator).uniqueFragmentNames,
	(*Validator).knownFragmentNames,
	(*Validator).noUnusedFragments,
	(*Validator).possibleFragmentSpreads,
	(*Validator).noFragmentCycles,
	(*Validator).uniqueVariableNames,
	(*Validator).noUndefinedVariables,
	(*Validator).noUnusedVariables,
	(*Validator).knownDirectives,
	(*Validator).uniqueDirectivesPerLocation,
	(*Validator).knownArgumentNames,
	(*Validator).uniqueArgumentNames,
	(*Validator).valuesOfCorrectType,
	(*Validator).providedRequiredArguments,
	(*Validator).variablesInAllowedPosition,
	(*Validator).overlappingFieldsCanBeMerged,
	(*Validator).uniqueInputFieldNames,
}

var sdlRules = []func(*Validator) *Errors{
	(*Validator).loneSchemaDefinition,
	(*Validator).knownDirectives,
	(*Validator).uniqueDirectivesPerLocation,
	(*Validator).knownArgumentNamesOnDirectives,
	(*Validator).uniqueArgumentNames,
	(*Validator).uniqueInputFieldNames,
	(*Validator).providedRequiredArgumentsOnDirectives,
}

func (v *Validator) executableDefinitions() *Errors {
	var errs *Errors

	errs = errs.Add(Error(errors.New("first error")))
	errs = errs.Add(Error(errors.New("second error")))

	return errs
}
func (v *Validator) fieldsOnCorrectType() *Errors {
	var errs *Errors

	errs = errs.Add(Error(errors.New("third error")))
	errs = errs.Add(Error(errors.New("fourth error")))

	return errs
}
func (v *Validator) fragmentsOnCompositeTypes() *Errors             { return nil }
func (v *Validator) knownArgumentNames() *Errors                    { return nil }
func (v *Validator) knownArgumentNamesOnDirectives() *Errors        { return nil }
func (v *Validator) knownDirectives() *Errors                       { return nil }
func (v *Validator) knownFragmentNames() *Errors                    { return nil }
func (v *Validator) knownTypeNames() *Errors                        { return nil }
func (v *Validator) loneAnonymousOperation() *Errors                { return nil }
func (v *Validator) loneSchemaDefinition() *Errors                  { return nil }
func (v *Validator) noFragmentCycles() *Errors                      { return nil }
func (v *Validator) noUndefinedVariables() *Errors                  { return nil }
func (v *Validator) noUnusedFragments() *Errors                     { return nil }
func (v *Validator) noUnusedVariables() *Errors                     { return nil }
func (v *Validator) overlappingFieldsCanBeMerged() *Errors          { return nil }
func (v *Validator) possibleFragmentSpreads() *Errors               { return nil }
func (v *Validator) providedRequiredArguments() *Errors             { return nil }
func (v *Validator) providedRequiredArgumentsOnDirectives() *Errors { return nil }
func (v *Validator) scalarLeafs() *Errors                           { return nil }
func (v *Validator) singleFieldSubscriptions() *Errors              { return nil }
func (v *Validator) uniqueArgumentNames() *Errors                   { return nil }
func (v *Validator) uniqueDirectivesPerLocation() *Errors           { return nil }
func (v *Validator) uniqueFragmentNames() *Errors                   { return nil }
func (v *Validator) uniqueInputFieldNames() *Errors                 { return nil }
func (v *Validator) uniqueOperationNames() *Errors                  { return nil }
func (v *Validator) uniqueVariableNames() *Errors                   { return nil }
func (v *Validator) valuesOfCorrectType() *Errors                   { return nil }
func (v *Validator) variablesAreInputTypes() *Errors                { return nil }
func (v *Validator) variablesInAllowedPosition() *Errors            { return nil }
