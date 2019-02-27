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

Benchmarks:

```
$ go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/language
BenchmarkLexer/bucketd-8                    	  300000	      4378 ns/op	     656 B/op	       3 allocs/op
BenchmarkLexer/graphql-go-8                 	  200000	     12037 ns/op	    1844 B/op	      30 allocs/op
BenchmarkLexer/vektah-8                     	  500000	      3435 ns/op	    1760 B/op	       8 allocs/op
BenchmarkTypeSystemParser/tsQuery/bucketd-8 	 1000000	      1225 ns/op	     688 B/op	      14 allocs/op
BenchmarkTypeSystemParser/tsQuery/vektah-8  	 1000000	      1784 ns/op	    1392 B/op	      24 allocs/op
BenchmarkParser/normalQuery/bucketd-8       	  100000	     13601 ns/op	    7648 B/op	      83 allocs/op
BenchmarkParser/normalQuery/graphql-go-8    	   30000	     40802 ns/op	   26975 B/op	     736 allocs/op
BenchmarkParser/normalQuery/vektah-8        	  100000	     20363 ns/op	   15776 B/op	     243 allocs/op
BenchmarkParser/tinyQuery/bucketd-8         	 3000000	       533 ns/op	     464 B/op	       7 allocs/op
BenchmarkParser/tinyQuery/graphql-go-8      	 1000000	      1558 ns/op	    1320 B/op	      35 allocs/op
BenchmarkParser/tinyQuery/vektah-8          	 2000000	       795 ns/op	     968 B/op	      13 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/language	20.184s
```

Test machine info:

* CPU: Intel Core i7-7700K @ 8x 5.0GHz
* RAM: 16GiB 3200MHz DDR4
* OS: Arch Linux 4.20.10-arch1-1-ARCH
* Go: version go1.12 linux/amd64

The benchmark code is included in this repository, please feel free to take a look at it yourself,
if you spot a mistake in our benchmark code that would give us an unfair advantage (or 
disadvantage!) then please let us know.

## License

MIT

[1]: http://facebook.github.io/graphql/June2018/
