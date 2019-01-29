package rules

import (
	"testing"
)

func TestScalarLeafs(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, scalarLeafs)
}
