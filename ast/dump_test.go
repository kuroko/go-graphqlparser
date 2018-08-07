package ast_test

import (
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
)

var (
	simpleSelectionField = strings.TrimSpace(`
{
	selection
}	
`)
)

func TestSdump(t *testing.T) {
	tt := []struct {
		descr string
		query string
	}{
		{descr: "simple selection field", query: simpleSelectionField},
	}

	for _, tc := range tt {
		psr := parser.New([]byte(tc.query))

		doc, err := psr.Parse()
		if err != nil {
			t.Fatal(err)
		}

		if sdump := ast.Sdump(doc); sdump != tc.query {
			t.Errorf("issue with %s:\n%s\n", tc.descr, sdump)
		}
	}
}
