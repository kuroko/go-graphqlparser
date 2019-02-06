package rules

import (
	"testing"
)

func TestUniqueVariableNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueVariableNames)
}
