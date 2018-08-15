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
    foo(message: """
    
    
        Hello,
            World

        From,
            GraphQL
            
            
    """)
}
`;

let ast = parser(query2);

console.log(util.inspect(ast, {showHidden: false, depth: null}))

process.exit(1);
//
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
