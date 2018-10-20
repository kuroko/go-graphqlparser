package graphql

import "testing"
import "github.com/stretchr/testify/require"

func TestPathInfo_MarshalJSON(t *testing.T) {
	t.Run("should return valid JSON", func(t *testing.T) {
		info := PathInfo{
			Kind:   PathInfoKindString,
			String: "Hello World",
		}

		bs, err := info.MarshalJSON()
		require.NoError(err)
	})
}
