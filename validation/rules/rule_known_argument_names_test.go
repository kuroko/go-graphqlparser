package rules

import (
	"testing"
)

func TestKnownArgumentNames(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, knownArgumentNames)
}
