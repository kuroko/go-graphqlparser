# Questions

* How is type extension expected to work? 
    * For example, `extend schema` throws an error if used with `makeExecutableSchema` from 
    `graphql-tools`.
    * Should the spec define how schema stitching should work with regards to type system 
    extensions? Currently it's very ambiguous, and it seems like there are many different approaches
    out in the wild (e.g. merging schemas without even using the `extend` keyword, loading schema 
    files in a specific order, or loading them all and concatenating schemas into one document to be 
    parsed) - each of these approaches probably triggers different validation errors if something 
    has gone wrong, as validation rules behave differently when extending an existing schema. 
   
* Should all validation rules be run, and all errors be returned? Or should the first error be 
returned by the server? For example, given the below "query" sent to the server, what is expected?
Currently we get back an error saying Bar is not executable, and nothing about the undefined `Asdf`
type being used.

```graphql
query Foo {
    bar
}

type Baz {
    # The type `Asdf` doesn't exist in the server's schema.
    qux: Asdf
}
```
