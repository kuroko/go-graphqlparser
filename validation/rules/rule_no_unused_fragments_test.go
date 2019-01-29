package rules

import (
	"testing"
)

func TestNoUnusedFragments(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, noUnusedFragments)
}
