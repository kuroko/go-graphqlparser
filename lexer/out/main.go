package main

import (
	"fmt"
	"strings"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
	"github.com/davecgh/go-spew/spew"
)

const query = `query foo {
	name
	model
	1.423e-10
}`

func main() {
	ipt := strings.NewReader(query)
	lxr := lexer.New(ipt)

	fmt.Println(query)
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
