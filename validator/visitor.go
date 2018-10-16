package validator

import (
	"github.com/bucketd/go-graphqlparser/ast"
)

type Visitor struct {
	Schema *Schema
	Errors *ast.Errors

	documentFuncs            []func(*Visitor, *ast.Document)
	operationDefinitionFuncs []func(*Visitor, *ast.ExecutableDefinition)
	queryDefinitionFuncs     []func(*Visitor, *ast.ExecutableDefinition)
}

func (v *Visitor) visitDocument(document *ast.Document) {
	document.Definitions.ForEach(func(definition *ast.Definition, _ int) {
		v.visitDefinition(definition)
	})

	for _, f := range v.documentFuncs {
		f(v, document)
	}
}

func (v *Visitor) visitDefinition(definition *ast.Definition) {
	switch definition.Kind {
	case ast.DefinitionKindExecutable:
		v.visitExecutableDefinition(definition.ExecutableDefinition)
	case ast.DefinitionKindTypeSystem:
		v.visitTypeSystemDefinition(definition.TypeSystemDefinition)
	case ast.DefinitionKindTypeSystemExtension:
		v.visitTypeSystemExtension(definition.TypeSystemExtension)
	}
}

func (v *Visitor) visitExecutableDefinition(executableDefinition *ast.ExecutableDefinition) {
	switch executableDefinition.Kind {
	case ast.ExecutableDefinitionKindOperation:
		v.visitOperationDefinition(executableDefinition)
	case ast.ExecutableDefinitionKindFragment:
		v.visitFragmentDefinition(executableDefinition)
	}
}

func (v *Visitor) visitTypeSystemDefinition(definition *ast.TypeSystemDefinition) {}
func (v *Visitor) visitTypeSystemExtension(definition *ast.TypeSystemExtension)   {}

func (v *Visitor) visitOperationDefinition(operationDefinition *ast.ExecutableDefinition) {
	switch operationDefinition.Kind {
	case ast.OperationTypeQuery:
		v.visitQuery(operationDefinition)
	case ast.OperationTypeMutation:
		v.visitMutation(operationDefinition)
	case ast.OperationTypeSubscription:
		v.visitSubscription(operationDefinition)
	}

	for _, f := range v.operationDefinitionFuncs {
		f(v, operationDefinition)
	}
}

func (v *Visitor) visitQuery(query *ast.ExecutableDefinition) {
	for _, f := range v.queryDefinitionFuncs {
		f(v, query)
	}
}

func (v *Visitor) visitMutation(query *ast.ExecutableDefinition)     {}
func (v *Visitor) visitSubscription(query *ast.ExecutableDefinition) {}
