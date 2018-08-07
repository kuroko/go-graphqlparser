# go-graphqlparser

A **work in progress** native Go GraphQL parser. Our aim is to produce an idiomatic, and extremely 
fast GraphQL parser that aheres to the [June 2018][1] GraphQL specification.

## Progress

* [x] Lexer
* [ ] Parser (in progress)
* [ ] Validation

### Benchmarks

Performance is one of this project's main goals, as such, we've kept a keen eye on benchmarks and 
tried to ensure that our benchmarks are fair, and reasonably comprehensive. Our results so far are
shown below.

Lexer benchmarks:

```
$ cd lexer
$ go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out
goos: linux
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/lexer
BenchmarkLexer/github.com/bucketd/go-graphqlparser-8         	 1000000	      1322 ns/op	       0 B/op	       0 allocs/op
BenchmarkLexer/github.com/graphql-go/graphql-8               	  300000	      4661 ns/op	     688 B/op	      17 allocs/op
BenchmarkLexer/github.com/graphql-gophers/graphql-go-8       	  300000	      5440 ns/op	    3184 B/op	     102 allocs/op
BenchmarkLexer/github.com/vektah/gqlparser-8                 	 1000000	      2000 ns/op	     208 B/op	       4 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/lexer	6.642s
```

Parser benchmarks:

```
$ cd parser
$ go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out
goos: linux
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/parser
BenchmarkParser/ultraMegaQuery/github.com/bucketd/go-graphqlparser-8         	      30	  36668885 ns/op	14071811 B/op	  129602 allocs/op
BenchmarkParser/ultraMegaQuery/github.com/graphql-go/graphql-8               	      10	 111869839 ns/op	37819180 B/op	 1059722 allocs/op
BenchmarkParser/ultraMegaQuery/github.com/graphql-gophers/graphql-go-8       	      30	  46932711 ns/op	19306851 B/op	  588318 allocs/op
BenchmarkParser/ultraMegaQuery/github.com/vektah/gqlparser-8                 	      20	  65345045 ns/op	29188164 B/op	  507217 allocs/op
BenchmarkParser/monsterQuery/github.com/bucketd/go-graphqlparser-8           	    3000	    366611 ns/op	  140854 B/op	    1298 allocs/op
BenchmarkParser/monsterQuery/github.com/graphql-go/graphql-8                 	    2000	    963880 ns/op	  378202 B/op	   10609 allocs/op
BenchmarkParser/monsterQuery/github.com/graphql-gophers/graphql-go-8         	    3000	    475450 ns/op	  194483 B/op	    5893 allocs/op
BenchmarkParser/monsterQuery/github.com/vektah/gqlparser-8                   	    3000	    597224 ns/op	  292099 B/op	    5081 allocs/op
BenchmarkParser/normalQuery/github.com/bucketd/go-graphqlparser-8            	  100000	     19851 ns/op	    6304 B/op	      76 allocs/op
BenchmarkParser/normalQuery/github.com/graphql-go/graphql-8                  	   30000	     55771 ns/op	   21524 B/op	     568 allocs/op
BenchmarkParser/normalQuery/github.com/graphql-gophers/graphql-go-8          	   50000	     31546 ns/op	   14064 B/op	     388 allocs/op
BenchmarkParser/normalQuery/github.com/vektah/gqlparser-8                    	   50000	     36563 ns/op	   15936 B/op	     244 allocs/op
BenchmarkParser/tinyQuery/github.com/bucketd/go-graphqlparser-8              	 2000000	       795 ns/op	     384 B/op	       6 allocs/op
BenchmarkParser/tinyQuery/github.com/graphql-go/graphql-8                    	 1000000	      2188 ns/op	    1064 B/op	      27 allocs/op
BenchmarkParser/tinyQuery/github.com/graphql-gophers/graphql-go-8            	 1000000	      1388 ns/op	    2040 B/op	      12 allocs/op
BenchmarkParser/tinyQuery/github.com/vektah/gqlparser-8                      	 1000000	      1365 ns/op	     968 B/op	      13 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/parser	28.466s
```

Test machine info:

* Model: Dell XPS 9560
* CPU: Intel Core i7-7700HQ @ 8x 3.8GHz
* RAM: 32GiB
* OS: Arch Linux
* Kernel: x86_64 Linux 4.17.12-arch1-1-ARCH
* Go: version go1.10.3 linux/amd64

The benchmark code is included in this repository, please feel free to take a look at it yourself,
if you spot a mistake in our benchmark code that would give us an unfair advantage (or 
disadvantage!) then please let us know.

## License

MIT

[1]: http://facebook.github.io/graphql/June2018/
