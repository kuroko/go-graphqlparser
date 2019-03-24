package graphql

import (
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/language"
)

// Parse ...
func Parse(doc []byte) (ast.Document, error) {
	return language.NewParser(doc).Parse()
}
