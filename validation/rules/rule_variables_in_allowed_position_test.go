package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestVariablesInAllowedPosition(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, rules.VariablesInAllowedPosition)
}
