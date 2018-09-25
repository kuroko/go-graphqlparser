# go-graphqlparser

A **work in progress** native Go GraphQL parser. Our aim is to produce an idiomatic, and extremely 
fast GraphQL parser that adheres to the [June 2018][1] GraphQL specification.

## Progress

* [x] Lexer
* [ ] Parser (in progress)
    * [x] Query parsing
    * [x] Type system parsing
    * [ ] Consistent, helpful errors
* [ ] Validation

### Benchmarks

Performance is one of this project's main goals, as such, we've kept a keen eye on benchmarks and 
tried to ensure that our benchmarks are fair, and reasonably comprehensive. Our results so far are
shown below.

Lexer benchmarks:

```
$ cd lexer
$ go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out
goos: darwin
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/lexer
BenchmarkLexer/github.com/bucketd/go-graphqlparser-8         	  300000	      4969 ns/op	     960 B/op	       5 allocs/op
BenchmarkLexer/github.com/graphql-go/graphql-8               	  100000	     14369 ns/op	    2176 B/op	      32 allocs/op
BenchmarkLexer/github.com/vektah/gqlparser-8                 	  300000	      4158 ns/op	    2000 B/op	       9 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/lexer	4.629s
```

Parser benchmarks:

```
$ cd parser
$ go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out
goos: darwin
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/parser
BenchmarkTypeSystemParser/tsQuery/github.com/bucketd/go-graphqlparser-8         	 2000000	       668 ns/op	     320 B/op	       8 allocs/op
BenchmarkTypeSystemParser/tsQuery/github.com/vektah/gqlparser-8                 	 1000000	      1224 ns/op	    1016 B/op	      16 allocs/op
BenchmarkParser/normalQuery/github.com/bucketd/go-graphqlparser-8               	  100000	     15289 ns/op	    6704 B/op	      81 allocs/op
BenchmarkParser/normalQuery/github.com/graphql-go/graphql-8                     	   30000	     41141 ns/op	   21505 B/op	     565 allocs/op
BenchmarkParser/normalQuery/github.com/vektah/gqlparser-8                       	   50000	     24931 ns/op	   15936 B/op	     244 allocs/op
BenchmarkParser/tinyQuery/github.com/bucketd/go-graphqlparser-8                 	 3000000	       543 ns/op	     384 B/op	       6 allocs/op
BenchmarkParser/tinyQuery/github.com/graphql-go/graphql-8                       	 1000000	      1525 ns/op	    1064 B/op	      27 allocs/op
BenchmarkParser/tinyQuery/github.com/vektah/gqlparser-8                         	 2000000	       944 ns/op	     968 B/op	      13 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/parser	15.885s
```

Test machine info:

* CPU: Intel Core i7-7700K @ 8x 4.2GHz
* RAM: 16GiB 3200MHz DDR4
* OS: macOS High Sierra (Version 10.13.6)
* Go: version go1.11 darwin/amd64

The benchmark code is included in this repository, please feel free to take a look at it yourself,
if you spot a mistake in our benchmark code that would give us an unfair advantage (or 
disadvantage!) then please let us know.

## License

MIT

[1]: http://facebook.github.io/graphql/June2018/
