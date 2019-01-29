package rules

import (
	"testing"
)

func TestPossibleFragmentSpreads(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, possibleFragmentSpreads)
}
