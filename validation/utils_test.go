package validation_test

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/validation"
	"github.com/stretchr/testify/assert"
)

func TestOrList(t *testing.T) {
	t.Run("should panic if no items are given", func(t *testing.T) {
		assert.Panics(t, func() {
			validation.OrList([]string{})
		})
	})

	t.Run("should return a lone item", func(t *testing.T) {
		assert.Equal(t, "A", validation.OrList([]string{"A"}))
	})

	t.Run("should return 2 items separate by ' or '", func(t *testing.T) {
		assert.Equal(t, "A or B", validation.OrList([]string{"A", "B"}))
	})

	t.Run("should return many items separate by commas, and the last by ' or ', up to 5 items", func(t *testing.T) {
		assert.Equal(t, "A, B, or C", validation.OrList([]string{"A", "B", "C"}))
		assert.Equal(t, "A, B, C, or D", validation.OrList([]string{"A", "B", "C", "D"}))
		assert.Equal(t, "A, B, C, D, or E", validation.OrList([]string{"A", "B", "C", "D", "E"}))
		assert.Equal(t, "A, B, C, D, or E", validation.OrList([]string{"A", "B", "C", "D", "E", "F"}))
		assert.Equal(t, "A, B, C, D, or E", validation.OrList([]string{"A", "B", "C", "D", "E", "F", "G"}))
	})
}

func TestQuotedOrList(t *testing.T) {
	t.Run("should panic if no items are given", func(t *testing.T) {
		assert.Panics(t, func() {
			validation.QuotedOrList([]string{})
		})
	})

	t.Run("should return a lone item", func(t *testing.T) {
		assert.Equal(t, `"A"`, validation.QuotedOrList([]string{"A"}))
	})

	t.Run("should return 2 items separate by ' or '", func(t *testing.T) {
		assert.Equal(t, `"A" or "B"`, validation.QuotedOrList([]string{"A", "B"}))
	})

	t.Run("should return many items separate by commas, and the last by ' or ', up to 5 items", func(t *testing.T) {
		assert.Equal(t, `"A", "B", or "C"`, validation.QuotedOrList([]string{"A", "B", "C"}))
		assert.Equal(t, `"A", "B", "C", or "D"`, validation.QuotedOrList([]string{"A", "B", "C", "D"}))
		assert.Equal(t, `"A", "B", "C", "D", or "E"`, validation.QuotedOrList([]string{"A", "B", "C", "D", "E"}))
		assert.Equal(t, `"A", "B", "C", "D", or "E"`, validation.QuotedOrList([]string{"A", "B", "C", "D", "E", "F"}))
		assert.Equal(t, `"A", "B", "C", "D", or "E"`, validation.QuotedOrList([]string{"A", "B", "C", "D", "E", "F", "G"}))
	})
}

func BenchmarkSuggestionList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		suggestions := validation.SuggestionList("Peettt", []string{
			"Pet",
			"Pet Shop Boys",
			"Badger",
			"Goat",
			"Bucketd",
			"Pucketd",
			"String",
			"Int",
			"ID",
			"Float",
			"Boolean",
			"Query",
			"Mutation",
			"Subscription",
		})
		_ = suggestions
	}
}

func TestSuggestionList(t *testing.T) {
	t.Run("should return results when input is empty", func(t *testing.T) {
		assert.NotEmpty(t, validation.SuggestionList("", []string{"a"}))
	})

	t.Run("should return an empty array when there are no options", func(t *testing.T) {
		assert.Empty(t, validation.SuggestionList("input", []string{}))
	})

	t.Run("should return options sorted based on similarity", func(t *testing.T) {
		assert.Equal(t, []string{"abc", "ab"}, validation.SuggestionList("abc", []string{"a", "ab", "abc"}))
	})
}
