# TODO

3 steps to validating SDL document:

* `graphql.ValidateSDL`: Runs the rules in `rules.go`.
* `buildSchema`: Builds up a `*graphql.Schema`, possibly extending an existing one.
* `graphql.ValidateSchema`: Validates that built schema.

BuildSchema is used in tests, which wil run ValidateSDL, and is more lenient, allowing us to write
shorter schema documents in tests.
