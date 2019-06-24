const { buildSchema, graphql, validateSchema } = require("graphql");

let schema = buildSchema(`
\t\ttype Query {
\t\t\tfoo: Foo
\t\t}

\t\tinterface Bar {
\t\t\tsup: String
\t\t\tbar: ID!
\t\t}

\t\ttype Foo implements Bar {
\t\t\tsup: String
\t\t\tbar: String!
\t\t}
`);

console.log(validateSchema(schema));

// graphql(schema, `{
//     bar
// }`).then(function (res) {
//     console.log(res);
// });
