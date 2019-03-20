const { ApolloServer,  makeExecutableSchema } = require("apollo-server");

const rootSchema = `
    schema {
        query: Query
    }

    type Query {
        post(id: Int, code: String): Post
    }
    
    extend type Post {
        author: String!
    }
`;

const postSchema = `
    type Post {
        title: String!
    }
`;

const schema = makeExecutableSchema({
    typeDefs: [ rootSchema, postSchema ],
    resolvers: {
        Query: {
            post: () => ({ title: null })
        }
    }
});

const server = new ApolloServer({ schema });

server.listen().then(({ url }) => {
    console.log(`Server running at ${url}`)
});
