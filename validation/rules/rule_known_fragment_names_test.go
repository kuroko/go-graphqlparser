package rules

import (
	"testing"
)

func TestKnownFragmentNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, knownFragmentNames)
}
