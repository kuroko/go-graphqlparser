package rules

import (
	"testing"
)

func TestFragmentsOnCompositeTypes(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, fragmentsOnCompositeTypes)
}
