package validation

import (
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
)

const (
	scalar = ast.TypeDefinitionKindScalar
	object = ast.TypeDefinitionKindObject
	iface  = ast.TypeDefinitionKindInterface
	union  = ast.TypeDefinitionKindUnion
	enum   = ast.TypeDefinitionKindEnum
	inobj  = ast.TypeDefinitionKindInputObject
)

// IsInputType ...
func IsInputType(schema *graphql.Schema, t ast.Type) bool {
	if t.Kind == ast.TypeKindList {
		return IsInputType(schema, *t.ListType)
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		switch def.Kind {
		case scalar, enum, inobj:
			return true
		}
	}

	return false
}

// IsOutputType ...
func IsOutputType(schema *graphql.Schema, t ast.Type) bool {
	if t.Kind == ast.TypeKindList {
		return IsOutputType(schema, *t.ListType)
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		switch def.Kind {
		case scalar, object, iface, union, enum:
			return true
		}
	}

	return false
}

// IsUnionMemberType ...
func IsUnionMemberType(schema *graphql.Schema, t ast.Type) bool {
	if t.NonNullable {
		return false
	}

	if t.Kind == ast.TypeKindList {
		return false
	}

	if def, ok := schema.Types[t.NamedType]; ok {
		if def.Kind == object {
			return true
		}
	}

	return false
}

// OrList takes a string slice, and returns each item comma separated, up to the last item which
// will be separated with ' or '.
func OrList(items []string) string {
	itemCount := len(items)
	if itemCount == 0 {
		panic("validation: expected at least one item")
	}

	if itemCount == 1 {
		return items[0]
	}

	if itemCount == 2 {
		return items[0] + " or " + items[1]
	}

	if itemCount > 5 {
		items = items[:5]
		itemCount = 5
	}

	items, last := items[:itemCount-1], items[itemCount-1]

	return strings.Join(items, ", ") + ", or " + last
}

// QuotedOrList is the same as OrList, but each item will also be wrapped in double-quotes.
func QuotedOrList(items []string) string {
	for i := 0; i < len(items); i++ {
		items[i] = `"` + items[i] + `"`
	}

	return OrList(items)
}

// SuggestionList accepts an invalid input string and a list of valid options, and returns a
// filtered list of the valid options sorted based on their similarity with the input.
//
// Source: https://github.com/graphql/graphql-js/blob/v14.2.0/src/jsutils/suggestionList.js
func SuggestionList(input string, options []string) []string {
	//ol := len(options)
	it := len(input) / 2

	optionsByDistance := make(map[string]int)

	for _, option := range options {
		ot := len(option) / 2

		threshold := it
		if ot > it {
			threshold = ot
		}

		if threshold < 1 {
			threshold = 1
		}

		distance := lexicalDistance(input, option)

		if distance <= threshold {
			optionsByDistance[option] = distance
		}
	}

	results := make([]string, 0, len(optionsByDistance))
	for _, o := range options {
		if _, ok := optionsByDistance[o]; ok {
			results = append(results, o)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		// Sort by length if the distance is the same, so we get the most specific suggestion first.
		if optionsByDistance[results[i]] == optionsByDistance[results[j]] {
			return len(results[i]) > len(results[j])
		}

		return optionsByDistance[results[i]] < optionsByDistance[results[j]]
	})

	return results
}

// lexicalDistance computes the lexical distance between strings A and B.
//
// The "distance" between two strings is given by counting the minimum number of edits needed to
// transform string A into string B. An edit can be an insertion, deletion, or substitution of a
// single character, or a swap of two adjacent characters.
//
// This distance can be useful for detecting typos in input or sorting
//
// Source: https://github.com/graphql/graphql-js/blob/v14.2.0/src/jsutils/suggestionList.js
func lexicalDistance(a, b string) int {
	if a == b {
		return 0
	}

	a = strings.ToLower(a)
	b = strings.ToLower(b)

	// Any case change counts as a single edit
	if a == b {
		return 1
	}

	return levenshtein.ComputeDistance(a, b)
}
