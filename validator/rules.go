package validator

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
)

func executableDefinitions(v *Visitor)     {} // Document
func fieldsOnCorrectType(v *Visitor)       {} // Field
func fragmentsOnCompositeTypes(v *Visitor) {} // InlineFragment | FragmentDefinition
func knownArgumentNames(v *Visitor)        {} // Argument
func knownDirectives(v *Visitor)           {} // Directive
func knownFragmentNames(v *Visitor)        {} // FragmentSpread
func knownTypeNames(v *Visitor)            {} // NamedType

func loneAnonymousOperation(v *Visitor) {
	var operationCount int

	v.documentFuncs = append(v.documentFuncs, func(v *Visitor, document *ast.Document) {
		operationCount = document.Definitions.Len()
	})

	v.queryDefinitionFuncs = append(v.queryDefinitionFuncs, func(v *Visitor, definition *ast.ExecutableDefinition) {
		if len(definition.Name) == 0 && operationCount > 1 {
			err := fmt.Errorf("anonymous operations should have no accompanying operations")
			v.Errors = v.Errors.Add(err)
		}
	})
}

func loneSchemaDefinition(v *Visitor)         {} // SchemaDefinition
func noFragmentCycles(v *Visitor)             {} //
func noUndefinedVariables(v *Visitor)         {} //
func noUnusedFragments(v *Visitor)            {} //
func noUnusedVariables(v *Visitor)            {} //
func overlappingFieldsCanBeMerged(v *Visitor) {} //
func possibleFragmentSpreads(v *Visitor)      {} //
func providedRequiredArguments(v *Visitor)    {} //
func scalarLeafs(v *Visitor)                  {} //
func singleFieldSubscriptions(v *Visitor)     {} //
func uniqueArgumentNames(v *Visitor)          {} //
func uniqueDirectivesPerLocation(v *Visitor)  {} //
func uniqueFragmentNames(v *Visitor)          {} //
func uniqueInputFieldNames(v *Visitor)        {} //

func uniqueOperationNames(v *Visitor) {
	names := map[string]struct{}{}

	v.operationDefinitionFuncs = append(v.operationDefinitionFuncs, func(v *Visitor, definition *ast.ExecutableDefinition) {
		if _, seen := names[definition.Name]; seen {
			err := fmt.Errorf("operation names must be unique, seen %s before", definition.Name)
			v.Errors = v.Errors.Add(err)
		}

		names[definition.Name] = struct{}{}
	})
}

func uniqueVariableNames(v *Visitor)        {} //
func valuesOfCorrectType(v *Visitor)        {} //
func variablesAreInputTypes(v *Visitor)     {} //
func variablesInAllowedPosition(v *Visitor) {} //
