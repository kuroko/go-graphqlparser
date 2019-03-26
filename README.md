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
$ go test -bench=. -benchmem -benchtime=10s
goos: linux
goarch: amd64
pkg: github.com/bucketd/go-graphqlparser/language
BenchmarkLexer/bucketd-8                       10000000     1866 ns/op      320 B/op      2 allocs/op
BenchmarkLexer/graphql-go-8                     1000000    10913 ns/op     1828 B/op     30 allocs/op
BenchmarkLexer/vektah-8                         5000000     3080 ns/op     1760 B/op      8 allocs/op
BenchmarkTypeSystemParser/tsQuery/bucketd-8     3000000     4205 ns/op     2480 B/op     38 allocs/op
BenchmarkTypeSystemParser/tsQuery/vektah-8      2000000     7819 ns/op     6000 B/op     91 allocs/op
BenchmarkParser/normalQuery/bucketd-8           1000000    10934 ns/op     7552 B/op     82 allocs/op
BenchmarkParser/normalQuery/graphql-go-8         500000    40478 ns/op    26983 B/op    736 allocs/op
BenchmarkParser/normalQuery/vektah-8            1000000    20071 ns/op    15792 B/op    243 allocs/op
BenchmarkParser/tinyQuery/bucketd-8            30000000      460 ns/op      448 B/op      7 allocs/op
BenchmarkParser/tinyQuery/graphql-go-8         10000000     1504 ns/op     1320 B/op     35 allocs/op
BenchmarkParser/tinyQuery/vektah-8             20000000      821 ns/op      984 B/op     13 allocs/op
PASS
ok  	github.com/bucketd/go-graphqlparser/language	190.444s
```

Test machine info:

* CPU: Intel Core i7-7700K @ 8x 5.0GHz
* RAM: 16GiB 3200MHz DDR4
* OS: Arch Linux 5.0.4-arch1-1-ARCH
* Go: version go1.12.1 linux/amd64

The benchmark code is included in this repository, please feel free to take a look at it yourself,
if you spot a mistake in our benchmark code that would give us an unfair advantage (or 
disadvantage!) then please let us know.

## License

MIT

[1]: http://facebook.github.io/graphql/June2018/
