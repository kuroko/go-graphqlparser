package ast

import "testing"

func TestType_String(t *testing.T) {
	tt := []struct {
		msg  string
		Type Type
		want string
	}{
		{
			msg: "Simple type",
			Type: Type{
				NamedType: "Foo",
			},
			want: "Foo",
		},
		{
			msg: "List",
			Type: Type{
				ListType: &Type{
					NamedType: "Bar",
				},
			},
			want: "[Bar]",
		},
		{
			msg: "Non nullable",
			Type: Type{
				NamedType:   "Baz",
				NonNullable: true,
			},
			want: "Baz!",
		},
		{
			msg: "Non nullable list",
			Type: Type{
				ListType: &Type{
					NamedType: "Qux",
				},
				NonNullable: true,
			},
			want: "[Qux]!",
		},
		{
			msg: "Non nullable list non nullable",
			Type: Type{
				ListType: &Type{
					NamedType:   "Quux",
					NonNullable: true,
				},
				NonNullable: true,
			},
			want: "[Quux!]!",
		},
		{
			msg: "Non nullable list non nullable list non nullable, you get the idea",
			Type: Type{
				ListType: &Type{
					ListType: &Type{
						ListType: &Type{
							NamedType:   "Quuz",
							NonNullable: true,
						},
						NonNullable: true,
					},
					NonNullable: true,
				},
				NonNullable: true,
			},
			want: "[[[Quuz!]!]!]!",
		},
	}
	for _, tc := range tt {
		t.Run(tc.msg, func(t *testing.T) {
			if got := tc.Type.String(); got != tc.want {
				t.Errorf("Type.String() = %v, want %v", got, tc.want)
			}
		})
	}
}
