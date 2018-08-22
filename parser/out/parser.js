const graphql = require("graphql/language");
const util = require('util');
const parser = graphql.parse;

const query = `
query foo($foo: Boolean = 2) {
	hello @foo(bar: "baz") {
		foo
		bar
	}
	world
}
`;

const query2 = `
query {
	foo(content: """
		Hello,
	
			Welcome to GraphQL. \\""" \t
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
			Welcome to GraphQL.
			Lets make this string a little bigger then. Because the larger this string
			becomes, the more efficient our lexer should look...
	
		From, Bucketd
	""")
}
`;

let ast = parser(query2);

console.log(util.inspect(ast, {showHidden: false, depth: null}));

process.exit(1);

// var start = process.hrtime();
//
// let ast;
//
// for (let i = 0; i < 1000000; i++) {
//     ast = parser(query2)
// }
//
// const NS_PER_SEC = 1e9;
// const diff = process.hrtime(start);
//
// console.log(`Benchmark took ${(diff[0] * NS_PER_SEC + diff[1]) / 1000000} ms`);
// console.log(ast)
