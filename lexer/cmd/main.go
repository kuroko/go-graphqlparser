package main

import (
    "bytes"
    "fmt"
    "time"

    "github.com/bucketd/go-graphqlparser/lexer"
)

func main() {
    start := time.Now()

    for i := 0; i < 1000000; i++ {
        l := lexer.New(bytes.NewReader([]byte("Hello, 世界")))

        for {
            r, w := l.Read()
            if r == rune(0) {
                break;
            }

            _ = r
            _ = w
        }
    }

    fmt.Println(time.Since(start))
}

