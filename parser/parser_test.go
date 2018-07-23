package parser

import (
	"testing"

	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

const (
	// query is a valid GraphQL query that contains at least one of each token type. It's re-used to
	// compare against other GraphQL libraries.
	query = `
		query withFragments {
			user(id: 4) {
				friends(first: 10) {
					...friendFields
				}
				mutualFriends(first: 10) {
					...friendFields
				}
			}
		}

		fragment friendFields on User {
			id
			name
			profilePic(size: 50)
		}
	`
)

func BenchmarkParser(b *testing.B) {
	qry := []byte(query)

	b.Run("github.com/bucketd/go-graphqlparser", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			psr := New(qry)

			ast, err := psr.Parse()
			if err != nil {
				b.Fatal(err)
			}

			_ = ast
		}
	})

	b.Run("github.com/graphql-go/graphql", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			params := parser.ParseParams{
				Source: source.NewSource(&source.Source{
					Body: qry,
				}),
			}

			ast, err := parser.Parse(params)
			if err != nil {
				b.Fatal(err)
			}

			_ = ast
		}
	})
}
