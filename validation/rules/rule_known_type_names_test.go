package rules

import (
	"testing"
)

func TestKnownTypeNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, knownTypeNames)
}
