package ast_test

import (
	"strings"
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/stretchr/testify/assert"
)

var (
	simpleSelectionField = strings.TrimSpace(`
{
  selection
}
`)

	wipTest = strings.TrimSpace(`
query Var($v: Int! = $var) {
  selection
}

query Vars($v: Int! = $var, $i: Int! = 123, $f: Float! = 1.23e+10, $s: String! = "string") {
  selection
}

query Vars2($b: Boolean! = true, $b2: Boolean! = false, $n: Int = null, $e: Enum = ENUM_VALUE) {
  selection
}

query Vars3($l: [Int!]! = [1, 2, 3], $o: Point2D = { x: 1.2, y: 3.4 }) {
  selection
}

query Directive @foo {
  selection
}

query Directives @foo(bar: $baz, baz: "qux") @bar @baz(foo: 123) {
  selection
}

query Selection {
  selection
}

query Selections {
  selection1
  selection2
  selection3 @foo
  selection4 @bar(baz: "qux")
  selection5 @baz(qux: 123) @foo @bar {
    nested {
      aliased: selections
    }
  }
}
`)
)

func TestSdump(t *testing.T) {
	tt := []struct {
		descr string
		query string
	}{
		{descr: "simple selection field", query: simpleSelectionField},
		{descr: "wip test", query: wipTest},
	}

	for _, tc := range tt {
		psr := parser.New([]byte(tc.query))

		doc, err := psr.Parse()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.query, ast.Sdump(doc))
	}
}
