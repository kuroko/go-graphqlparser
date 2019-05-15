const { buildSchema, graphql } = require("graphql");

let schema = buildSchema(`
    type Query {
        foo: String
    }
    
    directive @foo on SCHEMA | UNION
    
    type Mut {
        sup: String
    }
    
    union Foo
`);

graphql(schema, `{
    foo
}`).then(function (res) {
    console.log(res);
});
