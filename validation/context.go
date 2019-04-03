package validation

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql/types"
)

var (
	// queryContextDecoratorWalker ...
	queryContextDecoratorWalker = NewWalker([]VisitFunc{
		setMapSizes,
		setExecutableDefinition,
		setFragments,
		setReferencedFragments,
		setVariableUsages,
		setTypeDefinitions,
	})
	// sdlContextDecoratorWalker ...
	sdlContextDecoratorWalker = NewWalker([]VisitFunc{
		setMapSizes,
		setDirectiveDefinitions,
		setTypeDefinitions,
		//setSchemaDefinitionTypes,
	})
)

// NewContext instantiates a validation context struct, this involves the walker doing a
// preliminary pass of a query document, gathering basic information for the more complicated
// validation walk to come.
func NewContext(doc ast.Document, schema *types.Schema) *Context {
	if schema == nil {
		schema = &types.Schema{}
	}

	ctx := &Context{
		Document: doc,
		Schema:   schema,
	}

	queryContextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// NewSDLContext ...
func NewSDLContext(doc ast.Document, schema *types.Schema) *Context {
	isExtending := schema != nil
	if !isExtending {
		schema = &types.Schema{}
	}

	ctx := &Context{
		Document: doc,
		Schema:   schema,
	}

	// Construct SDL specific structures.
	ctx.SDLContext = &SDLContext{
		IsExtending: isExtending,
	}

	sdlContextDecoratorWalker.Walk(ctx, doc)

	return ctx
}

// Context ...
type Context struct {
	Document ast.Document
	Errors   *types.Errors
	Schema   *types.Schema

	// Used if we're validating an SDL file.
	SDLContext *SDLContext

	// TypeDefinitions is a map of all TypeDefinition nodes in the current document.
	TypeDefinitions map[string]*ast.TypeDefinition
	TypeExtensions  map[string]*ast.TypeExtension

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
}

// AddError adds an error to the linked list of errors on this Context.
func (ctx *Context) AddError(err types.Error) {
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
	// TODO: Maybe we could make this a slice too?.. and maybe we should, given VariableUsages
	// returns a string slice too, this is pretty inconsistent.
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

// SDLContext ...
type SDLContext struct {
	DirectiveDefinitions map[string]*ast.DirectiveDefinition
	KnownTypeNames       map[string]struct{}
	KnownDirectiveNames  map[string]struct{}
	KnownEnumValueNames  map[string]map[string]struct{}
	KnownFieldNames      map[string]map[string]struct{}

	QueryTypeDefined        bool
	MutationTypeDefined     bool
	SubscriptionTypeDefined bool

	// HasSeenSchemaDefinition ...
	HasSeenSchemaDefinition bool

	// IsExtending is true if this context was created with an existing Schema, and it's being
	// extended by another SDL file.
	IsExtending bool
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

// setMapSizes ...
func setMapSizes(w *Walker) {
	w.AddDocumentEnterEventHandler(func(ctx *Context, doc ast.Document) {
		if ctx.TypeDefinitions == nil {
			ctx.TypeDefinitions = make(map[string]*ast.TypeDefinition, doc.TypeDefinitions)
			ctx.TypeExtensions = make(map[string]*ast.TypeExtension, doc.TypeExtensions)
		}

		// SDL documents only (for now).
		if ctx.SDLContext != nil {
			if ctx.SDLContext.DirectiveDefinitions == nil {
				ctx.SDLContext.DirectiveDefinitions = make(map[string]*ast.DirectiveDefinition, doc.DirectiveDefinitions)
			}
		}
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

// setDirectiveDefinitions ...
func setDirectiveDefinitions(w *Walker) {
	w.AddDirectiveDefinitionEnterEventHandler(func(ctx *Context, def *ast.DirectiveDefinition) {
		ctx.SDLContext.DirectiveDefinitions[def.Name] = def
	})
}

// setTypeDefinitions ...
func setTypeDefinitions(w *Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *Context, def *ast.TypeDefinition) {
		ctx.TypeDefinitions[def.Name] = def
	})

	w.AddTypeExtensionEnterEventHandler(func(ctx *Context, ext *ast.TypeExtension) {
		ctx.TypeExtensions[ext.Name] = ext
	})
}
