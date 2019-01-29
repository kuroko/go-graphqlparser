package rules

import (
	"testing"
)

func TestPossibleTypeExtensions(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, possibleTypeExtensions)
}
