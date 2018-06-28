package main

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	query := `
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

	psr := parser.New([]byte(query))

	doc, err := psr.Parse()
	spew.Dump(doc)
	if err != nil {
		fmt.Println(err)
	}
}
