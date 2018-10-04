package validator

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
)

// Validator ...
type Validator struct {
	DocumentAST   *ast.Document
	GraphQLSchema *ast.Document
}

// New ...
// TODO: Error on nil document/schema?
func New(documentAST *ast.Document, gqlSchema *ast.Document) *Validator {
	return &Validator{
		DocumentAST:   documentAST,
		GraphQLSchema: gqlSchema,
	}
}

// ApplyExecutableRules ...
func (v *Validator) ApplyExecutableRules() *ast.Errors {
	return v.applyRules(executableRules)
}

// ApplySDLRules ...
func (v *Validator) ApplySDLRules() *ast.Errors {
	return v.applyRules(sdlRules)
}

// applyRules ...
func (v *Validator) applyRules(rules []rule) *ast.Errors {
	var errs *ast.Errors

	for _, erFn := range rules {
		e := erFn(v)
		if e == nil {
			continue
		}

		e.Join(errs)
		errs = e
	}

	return errs.Reverse()
}

type rule func(*Validator) *ast.Errors

var executableRules = []rule{
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

var sdlRules = []rule{
	(*Validator).loneSchemaDefinition,
	(*Validator).knownDirectives,
	(*Validator).uniqueDirectivesPerLocation,
	(*Validator).knownArgumentNamesOnDirectives,
	(*Validator).uniqueArgumentNames,
	(*Validator).uniqueInputFieldNames,
	(*Validator).providedRequiredArgumentsOnDirectives,
}

// https://facebook.github.io/graphql/June2018/#sec-Executable-Definitions
func (v *Validator) executableDefinitions() *ast.Errors {
	var errs *ast.Errors

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			err := fmt.Errorf("definition %d is of type %v, must be executable", i, def.Kind.String())
			errs = errs.Add(err)
		}
	})

	return errs
}

// https://facebook.github.io/graphql/June2018/#sec-Operation-Name-Uniqueness
func (v *Validator) uniqueOperationNames() *ast.Errors {
	var errs *ast.Errors

	seenOpName := make(map[string]bool, v.DocumentAST.Definitions.Len())

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			return
		}

		if seenOpName[def.ExecutableDefinition.Name] {
			err := fmt.Errorf("operation definitions should be unique, seen %s more than once", def.ExecutableDefinition.Name)
			errs = errs.Add(err)
		}
		seenOpName[def.ExecutableDefinition.Name] = true
	})

	return errs
}

// https://facebook.github.io/graphql/June2018/#sec-Lone-Anonymous-Operation
func (v *Validator) loneAnonymousOperation() *ast.Errors {
	var errs *ast.Errors

	var seenAnonOpCount int

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			return
		}

		if def.ExecutableDefinition.ShorthandQuery {
			seenAnonOpCount++
		}
	})

	if seenAnonOpCount > 1 {
		err := fmt.Errorf("seen %d shorthand queries", seenAnonOpCount)
		errs = errs.Add(err)
	}

	return errs
}

// https://facebook.github.io/graphql/June2018/#sec-Single-root-field
func (v *Validator) singleFieldSubscriptions() *ast.Errors {
	var errs *ast.Errors

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			return
		}

		if def.ExecutableDefinition.OperationType != ast.OperationTypeSubscription {
			return
		}

		if def.ExecutableDefinition.SelectionSet.Len() != 1 {
			err := fmt.Errorf("subscription %s must have exactly one root field", def.ExecutableDefinition.Name)
			errs = errs.Add(err)
		}
	})

	return errs
}

// https://facebook.github.io/graphql/June2018/#sec-Field-Selections-on-Objects-Interfaces-and-Unions-Types
func (v *Validator) fieldsOnCorrectType() *ast.Errors {
	var errs *ast.Errors

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			return
		}

		def.ExecutableDefinition.SelectionSet.ForEach(func(selection ast.Selection, _ int) {

		})

		if def.ExecutableDefinition.SelectionSet.Len() != 1 {
			err := fmt.Errorf("subscription %s must have exactly one root field", def.ExecutableDefinition.Name)
			errs = errs.Add(err)
		}
	})

	return errs
}

// recursive check
// func checkFieldsForType()

func (v *Validator) fragmentsOnCompositeTypes() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) knownArgumentNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) knownArgumentNamesOnDirectives() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) knownDirectives() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) knownFragmentNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) knownTypeNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) loneSchemaDefinition() *ast.Errors {
	var errs *ast.Errors

	var seenSchemaDefCount int

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind == ast.DefinitionKindTypeSystem {
			return
		}

		if def.TypeSystemDefinition.Kind == ast.TypeSystemDefinitionKindSchema {
			seenSchemaDefCount++
		}
	})

	if seenSchemaDefCount > 1 {
		err := fmt.Errorf("seen %d schema definitions", seenSchemaDefCount)
		errs = errs.Add(err)
	}

	return errs
}

func (v *Validator) noFragmentCycles() *ast.Errors {
	var errs *ast.Errors
	return errs
}
func (v *Validator) noUndefinedVariables() *ast.Errors {
	var errs *ast.Errors
	return errs
}
func (v *Validator) noUnusedFragments() *ast.Errors {
	var errs *ast.Errors
	return errs
}
func (v *Validator) noUnusedVariables() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) overlappingFieldsCanBeMerged() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) possibleFragmentSpreads() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) providedRequiredArguments() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) providedRequiredArgumentsOnDirectives() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) scalarLeafs() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) uniqueArgumentNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) uniqueDirectivesPerLocation() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) uniqueFragmentNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) uniqueInputFieldNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) uniqueVariableNames() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) valuesOfCorrectType() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) variablesAreInputTypes() *ast.Errors {
	var errs *ast.Errors
	return errs
}

func (v *Validator) variablesInAllowedPosition() *ast.Errors {
	var errs *ast.Errors
	return errs
}
