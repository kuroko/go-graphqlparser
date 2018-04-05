package main

import (
	"fmt"
	"strings"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/davecgh/go-spew/spew"
)

var query = `
# Mutation for liking a story.
# Foo bar baz.
mutation {
  likeStory(storyID: 123.53e-10, foo: "Hello, \u4e16") {
    story {
      likeCount
    }
  }
}`

func main() {
	qry := strings.TrimSpace(query)
	lxr := lexer.New([]byte(qry))

	fmt.Println(qry)
	fmt.Println()

	for {
		tok, err := lxr.Scan()
		if err != nil {
			panic(err)
		}

		if tok.Type == token.EOF {
			break
		}

		spew.Dump(tok)
	}
}
