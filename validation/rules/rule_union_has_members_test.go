package rules_test

//func TestUnionHasMembers(t *testing.T) {
//	tt := []ruleTestCase{
//		{
//			msg: "one member",
//			query: `
//				union Foo = Bar
//			`,
//		},
//		{
//			msg: "many members",
//			query: `
//				union Foo = Bar | Baz
//			`,
//		},
//		{
//			msg: "no members",
//			query: `
//				union Foo
//			`,
//			errs: (*graphql.Errors)(nil).
//				Add(validation.UnionHasNoMembersError("Foo", 0, 0)),
//		},
//	}
//
//	sdlRuleTester(t, tt, rules.UnionHasMembers)
//}
