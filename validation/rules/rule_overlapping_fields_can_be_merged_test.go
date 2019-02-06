package rules

import (
	"testing"
)

func TestOverlappingFieldsCanBeMerged(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, overlappingFieldsCanBeMerged)
}
