package types

import (
	"testing"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_MarshalJSON(t *testing.T) {
	t.Run("should return valid JSON", func(t *testing.T) {
		locations := []ast.Location{
			{Line: 6, Column: 7},
		}

		nodes := []ast.PathNode{
			{Kind: ast.PathNodeKindString, String: "hero"},
			{Kind: ast.PathNodeKindString, String: "heroFriends"},
			{Kind: ast.PathNodeKindInt, Int: 1},
			{Kind: ast.PathNodeKindString, String: "name"},
		}

		gqlErr := Error{
			Message:   "Name for character with ID 1002 could not be fetched.",
			Locations: ast.LocationsFromSlice(locations),
			Path:      ast.PathNodesFromSlice(nodes),
		}

		actual, err := gqlErr.MarshalJSON()
		require.NoError(t, err)

		expected := `{"message":"Name for character with ID 1002 could not be fetched.","locations":[{"line":6,"column":7}],"path":["hero","heroFriends",1,"name"]}`

		assert.Equal(t, expected, string(actual))
	})
}
