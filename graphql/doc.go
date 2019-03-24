// Package graphql presents an easy to use public interface for all of the functionality provided in
// this library. The major stages of dealing with a GraphQL query before execution can all be found
// as functions in this package. Each function is fairly lightweight, and it'd be easy to build your
// own variations of these functions if necessary.
//
// From a development standpoint, it also makes it easier to tie together parts of the library that
// would otherwise cause issues (e.g. with circular imports).
package graphql
