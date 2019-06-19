const { buildSchema, graphql, validateSchema } = require("graphql");

let schema = buildSchema(`
  type Query {
    bar: String @foo(foo: "hello")
  }
  
  directive @foo(foo: String, foo: Int) on QUERY | FIELD_DEFINITION
`);

console.log(validateSchema(schema));

// graphql(schema, `{
//     bar
// }`).then(function (res) {
//     console.log(res);
// });
