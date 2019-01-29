package rules

import (
	"testing"
)

func TestKnownDirectives(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, knownDirectives)
}
