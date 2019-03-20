# Questions

**Q:** How is type extension expected to work? 

Should the spec define how schema stitching should work with regards to type system extensions? 
Currently it's very ambiguous, and it seems like there are many different approaches out in the wild
(e.g. merging schemas without even using the `extend` keyword, loading schema files in a specific 
order, or loading them all and concatenating schemas into one document to be parsed) - each of these 
approaches probably triggers different validation errors if something has gone wrong, as validation 
rules behave differently when extending an existing schema.

**A:** This is still not clear, it's left to interpretation somewhat. Currently extending a schema
seems to behave inconsistently too, i.e. in the JS implementation if you extend a schema it does 
nothing, but if you extend a schema with an operation type that's already on the schema definition
it will throw an error.

In reality, I think it should be possible to extend anything, even in the same file, as long as it
doesn't conflict with existing fields. Type extensions would probably need to be evaluated after all 
type definitions have been processed, otherwise errors might not make sense (i.e. we'd show an error
for the main definition if it collided with an extension that was processed first).

--- 
   
**Q:** Should all validation rules be run, and all errors be returned? Or should the first error be 
returned by the server? For example, given the below "query" sent to the server, what is expected?


```graphql
query Foo {
    bar
}

type Baz {
    # The type `Asdf` doesn't exist in the server's schema.
    qux: Asdf
}
```

**A:** Validation rules are all run, and each error is returned.
