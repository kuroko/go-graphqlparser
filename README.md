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
goos: darwin
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/language
BenchmarkLexer/bucketd-8                      200000   7588 ns/op    960 B/op    5 allocs/op
BenchmarkLexer/graphql-go-8                   100000  21537 ns/op   2176 B/op   32 allocs/op
BenchmarkLexer/vektah-8                       200000   6301 ns/op   2000 B/op    9 allocs/op
BenchmarkTypeSystemParser/tsQuery/bucketd-8  1000000   2158 ns/op    624 B/op   16 allocs/op
BenchmarkTypeSystemParser/tsQuery/vektah-8    500000   3356 ns/op   1392 B/op   24 allocs/op
BenchmarkParser/normalQuery/bucketd-8         100000  22362 ns/op   6768 B/op   81 allocs/op
BenchmarkParser/normalQuery/graphql-go-8       20000  60579 ns/op  21504 B/op  565 allocs/op
BenchmarkParser/normalQuery/vektah-8           50000  36795 ns/op  15936 B/op  244 allocs/op
BenchmarkParser/tinyQuery/bucketd-8          2000000    807 ns/op    400 B/op    6 allocs/op
BenchmarkParser/tinyQuery/graphql-go-8       1000000   2276 ns/op   1064 B/op   27 allocs/op
BenchmarkParser/tinyQuery/vektah-8           1000000   1377 ns/op    968 B/op   13 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/language	21.860s
```

Test machine info:

* CPU: Intel Core i7-4770HQ @ 8x 2.2GHz (boost to 3.4GHz)
* RAM: 16GiB 1600MHz DDR3
* OS: macOS Mojave (Version 10.14)
* Go: version go1.11.1 darwin/amd64

The benchmark code is included in this repository, please feel free to take a look at it yourself,
if you spot a mistake in our benchmark code that would give us an unfair advantage (or 
disadvantage!) then please let us know.

## License

MIT

[1]: http://facebook.github.io/graphql/June2018/
