package validation

import "github.com/bucketd/go-graphqlparser/ast"

// Rules is a slice of ValidationRuleFunc that contains all validation rules defined by the GraphQL
// specification.
//
// The order of this list is important, and is intended to produce the most clear output when
// encountering multiple validation errors.
var Rules = []ast.ValidationRuleFunc{
	executableDefinitions,
}

// RulesSDL is a slice of ValidationRuleFunc that contains validation rules defined by the GraphQL
// specification for validating schema definition language documents exclusively. This set of rules
// is useful for servers that are parsing schemas, and other tools.
//
// The order of this list is important, and is intended to produce the most clear output when
// encountering multiple validation errors.
var RulesSDL []ast.ValidationRuleFunc
