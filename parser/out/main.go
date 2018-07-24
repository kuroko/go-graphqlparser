package main

import (
	"fmt"
	"time"

	"github.com/bucketd/go-graphqlparser/parser"
)

func main() {
	query := []byte(`
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
	`)

	start := time.Now()

	for i := 0; i < 1000000; i++ {
		psr := parser.New(query)

		doc, err := psr.Parse()
		if err != nil {
			fmt.Println(err)
		}

		_ = doc
	}

	fmt.Println(time.Since(start))
}
