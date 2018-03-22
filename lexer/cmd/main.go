package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
)

func main() {
	runtime.GOMAXPROCS(1)

	input := "query foo { name model }"

	start := time.Now()

	for i := 0; i < 5000000; i++ {
		lxr := lexer.New(input)

		for {
			tok := lxr.Scan()
			if tok.Type == token.EOF {
				break
			}

			_ = tok
		}
	}

	fmt.Println(time.Since(start))

	//start = time.Now()
	//
	//for i := 0; i < 1000000; i++ {
	//	for _, r := range []rune(input) {
	//		_ = r
	//	}
	//}
	//
	//fmt.Println(time.Since(start))
}
