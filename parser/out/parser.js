const graphql = require("graphql/language");
const parser = graphql.parse;

const query = `
query foo($foo: Boolean = 2) {
	hello @foo(bar: "baz") {
		foo
		bar
	}
	world
}
    `

var start = process.hrtime();

let ast;

for (let i = 0; i < 1000000; i++) {
    ast = parser(query)
}

const NS_PER_SEC = 1e9;
const diff = process.hrtime(start);

console.log(`Benchmark took ${(diff[0] * NS_PER_SEC + diff[1]) / 1000000} ms`);
console.log(ast)
