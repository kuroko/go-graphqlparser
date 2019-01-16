package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var contextDecoratorWalker = NewWalker([]VisitFunc{
	setFragment,
	setVariableUsages,
	setRecursiveVariableUsages,
	setRecursivelyReferencedFragments,
	setFragmentSpreads,
})

// NewContext instantiates a validation context struct, this involves the walker
// doing a preliminary pass of the document, gathering basic information for the
// more complicated validation walk to come.
func NewContext(doc ast.Document) *Context {
	ctx := &Context{}

	contextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Errors *graphql.Errors
	Schema *graphql.Schema

	// Used by validation rules.
	VariableDefs *ast.VariableDefinitions

	// Internal pre-cached with methods to access.
	fragments                      map[string]*ast.FragmentDefinition
	fragmentSpreads                map[*ast.Selections]map[string]bool
	recursivelyReferencedFragments map[string]map[string]bool
	variableUsages                 map[string]map[string]bool
	recursiveVariableUsages        map[string]map[string]bool

	fragmentSpreadsSelectionSet *ast.Selections
}

// Fragment returns a FragmentDefinition by name.
func (ctx *Context) Fragment(fragName string) *ast.FragmentDefinition {
	return ctx.fragments[fragName]
}

func setFragment(w *Walker) {
	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		if ctx.fragments == nil {
			ctx.fragments = make(map[string]*ast.FragmentDefinition)
		}

		ctx.fragments[fd.Name] = fd
	})
}

// FragmentSpreads returns all nested usages of fragment spreads in this Selections.
func (ctx *Context) FragmentSpreads(ss *ast.Selections) map[string]bool {
	return ctx.fragmentSpreads[ss]
}

func setFragmentSpreads(w *Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *Context, def *ast.OperationDefinition) {
		ctx.fragmentSpreadsSelectionSet = def.SelectionSet
	})

	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, def *ast.FragmentDefinition) {
		ctx.fragmentSpreadsSelectionSet = def.SelectionSet
	})

	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		if ctx.fragmentSpreads == nil {
			ctx.fragmentSpreads = make(map[*ast.Selections]map[string]bool)
		}

		if ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet] == nil {
			ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet] = make(map[string]bool)
		}

		ctx.fragmentSpreads[ctx.fragmentSpreadsSelectionSet][s.Name] = true
	})
}

// RecursivelyReferencedFragments returns all the recursively referenced
// fragments used by an operation or fragment definition.
func (ctx *Context) RecursivelyReferencedFragments(exDefName string) map[string]bool {
	return ctx.recursivelyReferencedFragments[exDefName]
}

func setRecursivelyReferencedFragments(w *Walker) {

}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(exDefName string) map[string]bool {
	return ctx.variableUsages[exDefName]
}

func setVariableUsages(w *Walker) {

}

// RecursiveVariableUsages returns all recursively referenced variable usages for an operation.
func (ctx *Context) RecursiveVariableUsages(opName string) map[string]bool {
	return ctx.recursiveVariableUsages[opName]
}

func setRecursiveVariableUsages(w *Walker) {

}
