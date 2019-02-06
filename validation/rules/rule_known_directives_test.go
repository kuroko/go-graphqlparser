package rules

import (
	"testing"
)

func TestKnownDirectives(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, knownDirectives)
}
