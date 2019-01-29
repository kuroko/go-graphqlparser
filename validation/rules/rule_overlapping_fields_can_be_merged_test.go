package rules

import (
	"testing"
)

func TestOverlappingFieldsCanBeMerged(t *testing.T) {
	tt := []ruleTestCase{}

	ruleTester(t, tt, overlappingFieldsCanBeMerged)
}
