package rules_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation/rules"
)

func TestUniqueOperationNames(t *testing.T) {
	tt := []ruleTestCase{}

	queryRuleTester(t, tt, rules.UniqueOperationNames)
}
