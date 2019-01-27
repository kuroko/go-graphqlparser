# Notes

## Context

What information is needed on Context? We should think about prioritising how we optimise storage of
this information, because it's this information that we may reference multiple times. The rest of
the validation process happens in the rules.

* The full AST document.
* Any errors that have been found.
* All fragments?
    * But can this just be kept on the AST document?
* Fragment spreads (i.e. non-recursive).
* Recursive fragment spreads.
* Variable usages.
* Recursive variable usages.

Is one of our main problems the fact that we haven't got a more strictly defined order of nodes for
our walker to follow when it's traversing the AST? For example, if we traversed selections on an
operation before it's variable definitions, could we walk only the nodes in the tree that we cared 
about walking, and in doing so be able to solve the recursive lookups by performing the traversal 
there and then?
