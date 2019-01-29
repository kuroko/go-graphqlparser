package rules

import (
	"testing"
)

func TestSingleFieldSubscriptions(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, singleFieldSubscriptions)
}
