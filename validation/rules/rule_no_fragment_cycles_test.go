package rules

import (
	"testing"
)

func TestNoFragmentCycles(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, noFragmentCycles)
}
