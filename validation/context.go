package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

var (
	// queryContextDecoratorWalker ...
	queryContextDecoratorWalker = NewWalker([]VisitFunc{
		setExecutableDefinition,
		setFragments,
		setReferencedFragments,
		setVariableUsages,
	})
	// sdlContextDecoratorWalker ...
	sdlContextDecoratorWalker = NewWalker([]VisitFunc{
		setSchemaDefinitionTypes,
	})
)

// NewQueryContext instantiates a validation context struct, this involves the walker doing a
// preliminary pass of a query document, gathering basic information for the more complicated
// validation walk to come.
func NewQueryContext(doc ast.Document, schema *Schema) *Context {
	return newContext(doc, schema, queryContextDecoratorWalker)
}

// NewSDLContext ...
func NewSDLContext(doc ast.Document, schema *Schema) *Context {
	return newContext(doc, schema, sdlContextDecoratorWalker)
}

// newContext ...
func newContext(doc ast.Document, schema *Schema, walker *Walker) *Context {
	isExtending := schema != nil
	if !isExtending {
		schema = &Schema{}
	}

	ctx := &Context{
		Document:    doc,
		Schema:      schema,
		IsExtending: isExtending,
	}

	walker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Document ast.Document
	Errors   *graphql.Errors
	Schema   *Schema

	// fragments contains all fragment definitions found in the input query, accessible by name.
	fragments map[string]*ast.FragmentDefinition

	// referencedFragments stores the fragment definitions referenced directly by an executable
	// definition, i.e. this is not recursively referenced fragments.
	referencedFragments map[*ast.ExecutableDefinition][]ast.Definition

	// variableUsages stores the variable usages referenced directly by an executable definition,
	// i.e. this is not recursive variable usages.
	variableUsages map[*ast.ExecutableDefinition][]string

	// executableDefinition is the current executable definition being walked over.
	executableDefinition *ast.ExecutableDefinition

	// IsExtending is true if this context was created with an existing Schema, and it's being
	// extended by another SDL file.
	IsExtending bool

	// HasSeenSchemaDefinition ...
	HasSeenSchemaDefinition bool
}

// AddError adds an error to the linked list of errors on this Context.
func (ctx *Context) AddError(err graphql.Error) {
	ctx.Errors = ctx.Errors.Add(err)
}

// Fragment ...
func (ctx *Context) Fragment(name string) *ast.FragmentDefinition {
	return ctx.fragments[name]
}

// VariableUsages returns the variable usages in an operation or fragment definition.
func (ctx *Context) VariableUsages(def *ast.ExecutableDefinition) []string {
	return ctx.variableUsages[def]
}

// RecursiveVariableUsages ...
func (ctx *Context) RecursiveVariableUsages(def *ast.ExecutableDefinition) map[string]struct{} {
	// Maybe we could make this a slice too?
	result := make(map[string]struct{})

	ctx.recursiveVariableUsagesIter(def, result, make(map[*ast.ExecutableDefinition]struct{}))

	return result
}

// recursiveVariableUsagesIter ...
func (ctx *Context) recursiveVariableUsagesIter(def *ast.ExecutableDefinition, agg map[string]struct{}, seen map[*ast.ExecutableDefinition]struct{}) {
	for _, vu := range ctx.variableUsages[def] {
		agg[vu] = struct{}{}
	}

	// TODO: Can this be swapped to use a cached version of recursively referenced fragments.
	for _, rd := range ctx.referencedFragments[def] {
		// We only want to recurse deeper if we've never seen the fragment before.
		if _, ok := seen[rd.ExecutableDefinition]; !ok {
			seen[rd.ExecutableDefinition] = struct{}{}
			ctx.recursiveVariableUsagesIter(rd.ExecutableDefinition, agg, seen)
		}
	}
}

// ReferencedFragments returns the fragments directly referenced by the given executable definition.
func (ctx *Context) ReferencedFragments(def *ast.ExecutableDefinition) []ast.Definition {
	return ctx.referencedFragments[def]
}

// RecursivelyReferencedFragments ...
// NOTE: In practice, def would likely be an operation definition, but this isn't a requirement.
func (ctx *Context) RecursivelyReferencedFragments(def *ast.ExecutableDefinition) map[ast.Definition]struct{} {
	if ctx.referencedFragments[def] == nil || len(ctx.referencedFragments[def]) == 0 {
		return nil
	}

	// Maybe we could make this a slice too?
	result := make(map[ast.Definition]struct{})

	ctx.recursivelyReferencedFragmentsIter(def, result, make(map[*ast.ExecutableDefinition]struct{}))

	return result
}

// recursivelyReferencedFragmentsIter is the inner iteration method for finding recursively
// referenced fragments for a given executable definition. It modifies the given aggregate map of
// results.
func (ctx *Context) recursivelyReferencedFragmentsIter(def *ast.ExecutableDefinition, agg map[ast.Definition]struct{}, seen map[*ast.ExecutableDefinition]struct{}) {
	// For each referenced fragment in the current executable definition...
	for _, rd := range ctx.referencedFragments[def] {
		agg[rd] = struct{}{}

		// We only want to recurse deeper if we've never seen the fragment before.
		if _, ok := seen[rd.ExecutableDefinition]; !ok {
			seen[rd.ExecutableDefinition] = struct{}{}
			ctx.recursivelyReferencedFragmentsIter(rd.ExecutableDefinition, agg, seen)
		}
	}
}

// setExecutableDefinition ...
func setExecutableDefinition(w *Walker) {
	w.AddExecutableDefinitionEnterEventHandler(func(ctx *Context, def *ast.ExecutableDefinition) {
		ctx.executableDefinition = def
	})
}

// setFragments ...
func setFragments(w *Walker) {
	w.AddFragmentDefinitionEnterEventHandler(func(ctx *Context, def *ast.FragmentDefinition) {
		if ctx.fragments == nil {
			ctx.fragments = make(map[string]*ast.FragmentDefinition)
		}

		ctx.fragments[def.Name] = def
	})
}

// setReferencedFragments ...
func setReferencedFragments(w *Walker) {
	w.AddFragmentSpreadSelectionEnterEventHandler(func(ctx *Context, s ast.Selection) {
		ctx.Document.Definitions.ForEach(func(d ast.Definition, i int) {
			if d.Kind != ast.DefinitionKindExecutable {
				return
			}

			if d.ExecutableDefinition.Kind != ast.ExecutableDefinitionKindFragment {
				return
			}

			if d.ExecutableDefinition.FragmentDefinition.Name == s.Name {
				if ctx.referencedFragments == nil {
					ctx.referencedFragments = make(map[*ast.ExecutableDefinition][]ast.Definition)
				}

				for _, v := range ctx.referencedFragments[ctx.executableDefinition] {
					if v == d {
						return
					}
				}

				ctx.referencedFragments[ctx.executableDefinition] =
					append(ctx.referencedFragments[ctx.executableDefinition], d)
			}
		})
	})
}

// setSchemaDefinitionTypes ...
func setSchemaDefinitionTypes(w *Walker) {
	w.AddSchemaDefinitionEnterEventHandler(func(ctx *Context, def *ast.SchemaDefinition) {
		def.RootOperationTypeDefinitions.ForEach(func(otd ast.RootOperationTypeDefinition, i int) {
			switch otd.OperationType {
			case ast.OperationDefinitionKindQuery:
				ctx.Schema.QueryType = &otd.NamedType
			case ast.OperationDefinitionKindMutation:
				ctx.Schema.MutationType = &otd.NamedType
			case ast.OperationDefinitionKindSubscription:
				ctx.Schema.SubscriptionType = &otd.NamedType
			}
		})
	})
}

// setVariableUsages ...
func setVariableUsages(w *Walker) {
	w.AddVariableValueEnterEventHandler(func(ctx *Context, v ast.Value) {
		if ctx.variableUsages == nil {
			ctx.variableUsages = make(map[*ast.ExecutableDefinition][]string)
		}

		for _, u := range ctx.variableUsages[ctx.executableDefinition] {
			if u == v.StringValue {
				return
			}
		}

		ctx.variableUsages[ctx.executableDefinition] = append(ctx.variableUsages[ctx.executableDefinition], v.StringValue)
	})
}
