package rules

import (
	"testing"
)

func TestNoUndefinedVariables(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, noUndefinedVariables)
}
