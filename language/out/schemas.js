const {ApolloServer, makeExecutableSchema} = require("apollo-server");

const rootSchema = `
  schema {
    query: Foo
  }

  schema {
    query: Foo
  }

  type Foo {
    bar: String!
  }

  extend type Foo {
    baz: Int!
  }
`;

const schema = makeExecutableSchema({
    typeDefs: [rootSchema],
    resolvers: {}
});

const server = new ApolloServer({schema});

server.listen().then(({url}) => {
    console.log(`Server running at ${url}`)
});
