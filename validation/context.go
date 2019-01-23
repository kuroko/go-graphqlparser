package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var contextDecoratorWalker = NewWalker([]VisitFunc{
	setDefinitionName,
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
	name                        string
	nameToSelections            map[string]*ast.Selections
}

func setDefinitionName(w *Walker) {
	w.AddOperationDefinitionEnterEventHandler(func(ctx *Context, od *ast.OperationDefinition) {
		ctx.name = od.Name
	})

	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		ctx.name = fd.Name
	})
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
	w.AddOperationDefinitionEnterEventHandler(func(ctx *Context, od *ast.OperationDefinition) {
		ctx.fragmentSpreadsSelectionSet = od.SelectionSet
	})

	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, fd *ast.FragmentDefinition) {
		ctx.fragmentSpreadsSelectionSet = fd.SelectionSet
	})

	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		if ctx.fragmentSpreads == nil {
			// TODO: Could this be keyed to operation / fragment name instead. We can extract all
			// fragments used within a definition, rather than being as granular as selections.
			// TODO: Is this going to be used elsewhere? If so, we might need to keep it like this.
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
	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		if ctx.nameToSelections == nil {
			ctx.nameToSelections = make(map[string]*ast.Selections)
		}

		// TODO: maybe this should happen later?
		if ctx.recursivelyReferencedFragments == nil {
			ctx.recursivelyReferencedFragments = make(map[string]map[string]bool)
		}

		// TODO: ctx.name here can be overwritten by a query and a fragment called the same thing.
		ctx.nameToSelections[ctx.name] = ctx.fragmentSpreadsSelectionSet
	})

	w.AddDocumentLeaveEventHandler(func(ctx *Context, d ast.Document) {
		for exDefName := range ctx.nameToSelections {
			if _, ok := ctx.recursivelyReferencedFragments[exDefName]; ok {
				continue
			}

			_ = recurseFrags(ctx, exDefName, []string{exDefName})
		}
	})
}

func recurseFrags(ctx *Context, name string, parents []string) []string {
	var children []string
	for frag := range ctx.FragmentSpreads(ctx.nameToSelections[name]) {
		if in(frag, parents) {
			continue
		}
		c := recurseFrags(
			ctx,
			frag,
			append(parents, name),
		)

		children = append(children, c...)
	}

	ctx.recursivelyReferencedFragments[name] = ctx.FragmentSpreads(ctx.nameToSelections[name])

	for _, child := range children {
		for frag := range ctx.FragmentSpreads(ctx.nameToSelections[child]) {
			ctx.recursivelyReferencedFragments[name][frag] = true
		}
	}

	return append(children, name)
}

func in(target string, candidates []string) bool {
	for _, c := range candidates {
		if c == target {
			return true
		}
	}
	return false
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(exDefName string) map[string]bool {
	return ctx.variableUsages[exDefName]
}

func setVariableUsages(w *Walker) {
	w.AddVariableValueEnterEventHandler(func(ctx *Context, v ast.Value) {
		if ctx.variableUsages == nil {
			ctx.variableUsages = make(map[string]map[string]bool)
		}

		if ctx.variableUsages[ctx.name] == nil {
			ctx.variableUsages[ctx.name] = make(map[string]bool)
		}

		ctx.variableUsages[ctx.name][v.StringValue] = true
	})
}

// RecursiveVariableUsages returns all recursively referenced variable usages for an operation.
func (ctx *Context) RecursiveVariableUsages(opName string) map[string]bool {
	//return ctx.recursiveVariableUsages[opName]
	return ctx.variableUsages[opName]
}

func setRecursiveVariableUsages(w *Walker) {

}
