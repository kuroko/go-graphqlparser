package rules

import (
	"testing"
)

func TestKnownFragmentNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, knownFragmentNames)
}
