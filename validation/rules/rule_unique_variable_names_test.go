package rules

import (
	"testing"
)

func TestUniqueVariableNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, uniqueVariableNames)
}
