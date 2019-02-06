package rules

import (
	"testing"
)

func TestFragmentsOnCompositeTypes(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, fragmentsOnCompositeTypes)
}
