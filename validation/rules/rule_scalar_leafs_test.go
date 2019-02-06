package rules

import (
	"testing"
)

func TestScalarLeafs(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, scalarLeafs)
}
