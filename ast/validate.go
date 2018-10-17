package ast

// ValidationContext ...
type ValidationContext struct {
	Errors *Errors
	Schema *Schema
}

// ValidationRuleFunc ...
type ValidationRuleFunc func(vctx *ValidationContext, walker *Walker)

// Validate ...
func Validate(doc Document, rules []ValidationRuleFunc) *Errors {
	vctx := &ValidationContext{}
	walker := NewWalker()

	// Apply all specified validation rules to the walker so that the AST can be validated as it is
	// traversed by the walker, populating the ValidationContext along the way.
	for _, rule := range rules {
		rule(vctx, walker)
	}

	walker.Walk(doc)

	return vctx.Errors
}
