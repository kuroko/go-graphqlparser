package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

// NewContext instantiates a validation context struct, this involves the walker
// doing a preliminary pass of the document, gathing basic information for the
// more complicated validation walk to come.
func NewContext(doc ast.Document) *Context {
	ctx := &Context{}

	visitFns := []ast.VisitFunc{
		setFragment(ctx),
		setVariableUsages(ctx),
		setRecursiveVariableUsages(ctx),
		setRecursivelyReferencedFragments(ctx),
		setFragmentSpreads(ctx),
	}

	walker := ast.NewWalker(visitFns)
	walker.Walk(doc)

	return ctx
}

// Context ...
type Context struct {
	Errors *graphql.Errors
	Schema *graphql.Schema

	fragment                       map[string]*ast.FragmentDefinition
	fragmentSpreads                map[*ast.Selections]map[string]bool
	recursivelyReferencedFragments map[*ast.ExecutableDefinition]map[string]bool
	variableUsages                 map[*ast.ExecutableDefinition]map[string]bool
	recursiveVariableUsages        map[string]map[string]bool
}

// Fragment returns a FragmentDefinition by name.
func (ctx *Context) Fragment(fragName string) *ast.FragmentDefinition {
	return ctx.fragment[fragName]
}

func setFragment(ctx *Context) ast.VisitFunc {
	return func(w *ast.Walker) {
		w.AddFragmentDefinitionEnterEventHandler(func(fragDef *ast.FragmentDefinition) {
			ctx.fragment[fragDef.Name] = fragDef
		})
	}
}

// FragmentSpreads returns all nested usages of fragment spreads in this Selections.
func (ctx *Context) FragmentSpreads(selections *ast.Selections) map[string]bool {
	return ctx.fragmentSpreads[selections]
}

func setFragmentSpreads(ctx *Context) ast.VisitFunc {
	return func(w *ast.Walker) {}
}

// RecursivelyReferencedFragments returns all the recursively referenced
// fragments used by an operation or fragment definition.
func (ctx *Context) RecursivelyReferencedFragments(exDef *ast.ExecutableDefinition) map[string]bool {
	return ctx.recursivelyReferencedFragments[exDef]
}

func setRecursivelyReferencedFragments(ctx *Context) ast.VisitFunc {
	return func(w *ast.Walker) {}
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(exDef *ast.ExecutableDefinition) map[string]bool {
	return ctx.variableUsages[exDef]
}

func setVariableUsages(ctx *Context) ast.VisitFunc {
	return func(w *ast.Walker) {}
}

// RecursiveVariableUsages returns all recursively referenced variable usages for an operation.
func (ctx *Context) RecursiveVariableUsages(opName string) map[string]bool {
	return ctx.recursiveVariableUsages[opName]
}

func setRecursiveVariableUsages(ctx *Context) ast.VisitFunc {
	return func(w *ast.Walker) {}
}
