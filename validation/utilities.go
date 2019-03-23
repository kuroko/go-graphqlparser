package validation

import (
	"errors"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
)

// BuildASTSchema ...
// TODO: Do we need a previous schema that we can pass in here?
// TODO: Do we really need to pass in the walker here...
func BuildASTSchema(doc ast.Document, walker *Walker) (*Schema, error) {
	errs := ValidateSDL(doc, nil, walker)
	if errs.Len() > 0 {
		// TODO: This is kinda useless right now...
		return nil, errors.New("validation: found errors validating schema")
	}

	schema := &Schema{}
	schemaVisitFns := []VisitFunc{
		setSchemaTypeDefinitions,
	}

	// Traverse the schema AST, populating the schema object with relevant information.
	NewWalker(schemaVisitFns).Walk(&Context{Schema: schema}, doc)

	return schema, nil
}

// BuildSchema ...
// TODO: Do we need a previous schema that we can pass in here?
// TODO: Should the input be bytes here?
func BuildSchema(doc string, walker *Walker) (*Schema, error) {
	schemaParser := language.NewParser([]byte(doc))

	schemaDoc, err := schemaParser.Parse()
	if err != nil {
		// TODO: Error wrapping. Maybe some kind of context?
		return nil, err
	}

	return BuildASTSchema(schemaDoc, walker)
}

// setSchemaTypeDefinitions ...
func setSchemaTypeDefinitions(w *Walker) {
	w.AddTypeDefinitionEnterEventHandler(func(ctx *Context, def *ast.TypeDefinition) {
		if ctx.Schema.Types == nil {
			ctx.Schema.Types = make(map[string]*ast.TypeDefinition)
		}

		ctx.Schema.Types[def.Name] = def
	})
}
