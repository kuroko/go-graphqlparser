const {ApolloServer, makeExecutableSchema} = require("apollo-server");

const rootSchema = `
    type Query {
        foo: String
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
