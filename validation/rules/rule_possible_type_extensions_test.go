package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestPossibleTypeExtensions(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, rules.PossibleTypeExtensions)
}
