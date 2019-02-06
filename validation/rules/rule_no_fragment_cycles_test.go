package rules

import (
	"testing"
)

func TestNoFragmentCycles(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, noFragmentCycles)
}
