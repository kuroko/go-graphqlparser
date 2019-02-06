package rules

import (
	"testing"
)

func TestLoneAnonymousOperation(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, loneAnonymousOperation)
}
