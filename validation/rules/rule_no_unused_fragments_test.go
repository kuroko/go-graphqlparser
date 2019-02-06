package rules

import (
	"testing"
)

func TestNoUnusedFragments(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, noUnusedFragments)
}
