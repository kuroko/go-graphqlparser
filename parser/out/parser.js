const graphql = require("graphql/language");
const parser = graphql.parse;

const query = `
query first($foo: Boolean = true) {
    # How about some comments too?
    user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
        ...userFields
    }
    # Do comments even slow us down?
    post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
        ...postFields
    }
    # Directives
    fooa: foo @include(if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo) @include(if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo, if: $foo)
    bara: bar @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar)
    baza: baz @gt(val: $baz)
    # Inline fragments
    ... @include(if: $expandedInfo) {
        firstName
        lastName
        birthday
    }
}

mutation second($variable: String = "test") {
    sendEmail(message: """
        Hello,
            World!

        Yours,
            GraphQL
    """)
    sendEmail2(message: "Hello\\n,  World!\\n\\nYours,\\n  GraphQL.")
    intVal(foo: 12345)
    floatVal(bar: 123.456)
    floatVal2(bar: 123.456e10)
    boolVal(bool: false)
    listVal(list: [1, 2, 3])
    variableVal(var: $variable)
}

subscription third {
    ...postFields
}

fragment userFields on User {
    firstName
    lastName
    title
    company {
        name
        slug
    }
    email
    mobile
}

fragment postFields on Subscription {
    title
    subtitle
    slug
    author {
        ...userFields
    }
    category {
        name
        slug
    }
    content
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
