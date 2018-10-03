package validator

import (
	"errors"
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

func (v *Validator) executableDefinitions() *ast.Errors {
	var errs *ast.Errors

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if def.Kind != ast.DefinitionKindExecutable {
			err := fmt.Errorf("definition %d is of type %v, must be executable", i, def.Kind.String())
			errs.Add(err)
		}
	})

	return errs
}
func (v *Validator) fieldsOnCorrectType() *ast.Errors {
	var errs *ast.Errors

	errs = errs.Add(ast.Error(errors.New("third error")))
	errs = errs.Add(ast.Error(errors.New("fourth error")))

	return errs
}
func (v *Validator) fragmentsOnCompositeTypes() *ast.Errors      { return nil }
func (v *Validator) knownArgumentNames() *ast.Errors             { return nil }
func (v *Validator) knownArgumentNamesOnDirectives() *ast.Errors { return nil }
func (v *Validator) knownDirectives() *ast.Errors                { return nil }
func (v *Validator) knownFragmentNames() *ast.Errors             { return nil }
func (v *Validator) knownTypeNames() *ast.Errors                 { return nil }
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
		errs.Add(err)
	}

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
		errs.Add(err)
	}

	return errs
}
func (v *Validator) noFragmentCycles() *ast.Errors                      { return nil }
func (v *Validator) noUndefinedVariables() *ast.Errors                  { return nil }
func (v *Validator) noUnusedFragments() *ast.Errors                     { return nil }
func (v *Validator) noUnusedVariables() *ast.Errors                     { return nil }
func (v *Validator) overlappingFieldsCanBeMerged() *ast.Errors          { return nil }
func (v *Validator) possibleFragmentSpreads() *ast.Errors               { return nil }
func (v *Validator) providedRequiredArguments() *ast.Errors             { return nil }
func (v *Validator) providedRequiredArgumentsOnDirectives() *ast.Errors { return nil }
func (v *Validator) scalarLeafs() *ast.Errors                           { return nil }
func (v *Validator) singleFieldSubscriptions() *ast.Errors              { return nil }
func (v *Validator) uniqueArgumentNames() *ast.Errors                   { return nil }
func (v *Validator) uniqueDirectivesPerLocation() *ast.Errors           { return nil }
func (v *Validator) uniqueFragmentNames() *ast.Errors                   { return nil }
func (v *Validator) uniqueInputFieldNames() *ast.Errors                 { return nil }
func (v *Validator) uniqueOperationNames() *ast.Errors {
	var errs *ast.Errors

	seenOpName := make(map[string]bool, v.DocumentAST.Definitions.Len())

	v.DocumentAST.Definitions.ForEach(func(def ast.Definition, i int) {
		if seenOpName[def.ExecutableDefinition.Name] {
			err := fmt.Errorf("operation definitions should be unique, seen %s more than once", def.ExecutableDefinition.Name)
			errs.Add(err)
		}
		seenOpName[def.ExecutableDefinition.Name] = true
	})

	return errs
}
func (v *Validator) uniqueVariableNames() *ast.Errors        { return nil }
func (v *Validator) valuesOfCorrectType() *ast.Errors        { return nil }
func (v *Validator) variablesAreInputTypes() *ast.Errors     { return nil }
func (v *Validator) variablesInAllowedPosition() *ast.Errors { return nil }
