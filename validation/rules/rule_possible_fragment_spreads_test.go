package rules

import (
	"testing"
)

func TestPossibleFragmentSpreads(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, possibleFragmentSpreads)
}
