const graphql = require("graphql/language");
const parser = graphql.parse;

const query = `
        query withFragments {
            user(id: 4) {
                friends(first: 10) {
                    ...friendFields
                }
                mutualFriends(first: 10) {
                    ...friendFields
                }
            }
        }

        fragment friendFields on User {
            id
            name
            profilePic(size: 50)
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
