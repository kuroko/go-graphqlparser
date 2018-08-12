package parser

var (
	bigQuery = []byte(`
query first($foo: Boolean = true, $foo: Boolean = true) {
    # How about some comments too?
    user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376", id: "3931a3fc-d4f9-4faa-bcf5-882022617376", id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
        ...userFields
    }
    # Do comments even slow us down?
    post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e", id: "489c9250-50b9-4612-b930-56dc4e1ae44e", id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
        ...postFields
    }
    # Directives
    fooa: foo @include(if: $foo) @include(if: $foo) @include(if: $foo) @include(if: $foo) @include(if: $foo) @include(if: $foo)
    bara: bar @skip(if: $bar) @skip(if: $bar) @skip(if: $bar) @skip(if: $bar)
    baza: baz @gt(val: $baz)
    # Inline fragments
    ... @include(if: $expandedInfo) {
        firstName
        lastName
        birthday
    }
}
mutation second($variable: String = "test") {
    sendEmail(message: """
        Hello,
            World!

        Yours,
            GraphQL
    """)
    sendEmail2(message: "Hello\n,  World!\n\nYours,\n  GraphQL.")
    intVal(foo: 12345, foo: 12345, foo: 12345)
    floatVal(bar: 123.456, bar: 123.456, bar: 123.456, bar: 123.456, bar: 123.456, bar: 123.456, bar: 123.456)
    floatVal2(bar: 123.456e10)
    boolVal(bool: false)
    listVal(list: [1, 2, 3], list: [1, 2, 3], list: [1, 2, 3], list: [1, 2, 3], list: [1, 2, 3], list: [1, 2, 3], list: [1, 2, 3])
    variableVal(var: $variable, var: $variable, var: $variable, var: $variable, var: $variable, var: $variable)
}
`)

	normalQuery = []byte(`
query first($foo: Boolean = true, $foo: Boolean = true) {
    # How about some comments too?
    user(id: "3931a3fc-d4f9-4faa-bcf5-882022617376") {
        ...userFields
    }
    # Do comments even slow us down?
    post(id: "489c9250-50b9-4612-b930-56dc4e1ae44e") {
        ...postFields
    }
    # Directives
    fooa: foo @include(if: $foo)
    bara: bar @skip(if: $bar)
    baza: baz @gt(val: $baz)
    # Inline fragments
    ... @include(if: $expandedInfo) {
        firstName
        lastName
        birthday
    }
}

mutation second($variable: String = "test") {
    sendEmail(message: """
        Hello,
            World!

        Yours,
            GraphQL
    """)
    sendEmail2(message: "Hello\n,  World!\n\nYours,\n  GraphQL.")
    intVal(foo: 12345)
    floatVal(bar: 123.456)
    floatVal2(bar: 123.456e10)
    boolVal(bool: false)
    listVal(list: [1, 2, 3])
    variableVal(var: $variable)
}

fragment userFields on User {
    firstName
    lastName
    title
    company {
        name
        slug
    }
    email
    mobile
}

fragment postFields on Subscription {
    title
    subtitle
    slug
    author {
        ...userFields
    }
    category {
        name
        slug
    }
    content
}
	`)

	tinyQuery = []byte(`
		{
			person {
				name
			}
		}
	`)
)
