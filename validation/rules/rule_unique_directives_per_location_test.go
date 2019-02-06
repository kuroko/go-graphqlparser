package rules

import (
	"testing"
)

func TestUniqueDirectivesPerLocation(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, uniqueDirectivesPerLocation)
}
