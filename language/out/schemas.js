const {ApolloServer, makeExecutableSchema} = require("apollo-server");

const rootSchema = `
    schema {
        query: Query
        mutation: Mutation
    }

    type Query {
        post(id: Int, code: String): Post
    }
    
    type Mutation {
        _empty: String
    }   
    
    extend type Post {
        title: Int!
        author: String!
    }
`;

const postSchema = `    
    type Mutation2 {
        doSomething: String!
    }

    type Post {
        title: String!
    }
`;

const schema = makeExecutableSchema({
    typeDefs: [rootSchema, postSchema],
    resolvers: {
        Query: {
            post: () => ({
                title: "Hello",
                author: "World",
            })
        },
        Mutation2: {
            doSomething: () => "Something happened"
        }
    }
});

const server = new ApolloServer({schema});

server.listen().then(({url}) => {
    console.log(`Server running at ${url}`)
});
