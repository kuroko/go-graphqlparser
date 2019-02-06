package rules

import (
	"testing"
)

func TestSingleFieldSubscriptions(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, singleFieldSubscriptions)
}
