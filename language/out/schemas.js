const {ApolloServer, makeExecutableSchema} = require("apollo-server");

const rootSchema = `
    schema {
        query: Query
        mutation: Mutation
    }
    
    input Foo {
        bar: String
        baz: String
    }
    
    directive @foo on OBJECT

    type Query {
        post(id: Int, code: String, foo: Foo): Post
    }
    
    type Mutation {
        _empty: String
    }
    
    interface SomeInterface {
        fieldOnInterface: String
    }

    type Post implements SomeInterface {
        title: String!
    }
    
    extend type Post @foo {
        author: String!
    }
`;

const postSchema = `
    type Mutation2 {
        doSomething: String!
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
