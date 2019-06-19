const {ApolloServer, makeExecutableSchema} = require("apollo-server");

const rootSchema = `
  type Query {
    bar: String!
  }
  
  directive @foo(foo: String, foo: String) on QUERY
`;

const schema = makeExecutableSchema({
    typeDefs: [rootSchema],
    resolvers: {}
});

const server = new ApolloServer({schema});

server.listen().then(({url}) => {
    console.log(`Server running at ${url}`)
});
