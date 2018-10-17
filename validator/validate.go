package validator

import "github.com/bucketd/go-graphqlparser/ast"

// rules ...
var rules = []validationRuleFunc{
	executableDefinitions,
}

// sdlRules ...
var sdlRules []validationRuleFunc

// validationContext ...
type validationContext struct {
	errors *ast.Errors
	schema *Schema
}

// validationRuleFunc ...
type validationRuleFunc func(vctx *validationContext, walker *ast.Walker)

// Validate ...
func Validate(doc ast.Document) *ast.Errors {
	return validate(doc, rules)
}

// ValidateSDL ...
func ValidateSDL(doc ast.Document) *ast.Errors {
	return validate(doc, sdlRules)
}

// validate ....
func validate(doc ast.Document, ruleFns []validationRuleFunc) *ast.Errors {
	vctx := &validationContext{}
	walker := ast.NewWalker()

	for _, rule := range ruleFns {
		rule(vctx, walker)
	}

	walker.Walk(doc)

	return vctx.errors
}
