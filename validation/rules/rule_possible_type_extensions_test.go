package rules

import (
	"testing"
)

func TestPossibleTypeExtensions(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, possibleTypeExtensions)
}
