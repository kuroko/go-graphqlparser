package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bucketd/go-graphqlparser/lexer"
	"github.com/bucketd/go-graphqlparser/token"
)

var query = `
# Mutation for liking a story.
# Foo bar baz.
mutation {
  likeStory(storyID: 123.53e-10) {
    story {
      likeCount
    }
  }
}`

func main() {
	runtime.GOMAXPROCS(4)

	input := []byte("query \"\u4e16\" 0.001 foo { name model foo bar baz qux }")

	//lxr := lexer.New(input)
	//
	//for {
	//	tok, _ := lxr.Scan()
	//	if tok.Type == token.EOF {
	//		break
	//	}
	//
	//	spew.Dump(tok)
	//}
	//
	//os.Exit(0)

	fmt.Println("==> Single-threaded, 5,000,000 iterations:")

	start := time.Now()

	for j := 0; j < 5000000; j++ {
		lxr := lexer.New(input)

		for {
			tok, err := lxr.Scan()
			if err != nil {
				panic(err)
			}

			if tok.Type == token.EOF {
				break
			}

			_ = tok
		}
	}

	fmt.Printf("    %s\n", time.Since(start))
	fmt.Println("==> Multi-threaded, 5,000,000 iterations per core:")

	start = time.Now()

	var wg sync.WaitGroup

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func(w int) {
			wstart := time.Now()

			for j := 0; j < 5000000; j++ {
				lxr := lexer.New(input)

				for {
					tok, err := lxr.Scan()
					if err != nil {
						panic(err)
					}

					if tok.Type == token.EOF {
						break
					}

					_ = tok
				}
			}

			fmt.Printf("    Worker %d finished. Took: %s\n", w, time.Since(wstart))
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Printf("    %s\n", time.Since(start))

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
