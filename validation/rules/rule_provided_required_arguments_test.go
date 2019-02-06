package rules

import (
	"testing"
)

func TestProvidedRequiredArguments(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, providedRequiredArguments)
}
