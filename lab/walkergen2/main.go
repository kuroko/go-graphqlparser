package main

import (
	"text/template"
	"fmt"
)

var tis = []TypeInfo{
	{
		TypeName: "Document",
		VariableName: "document",
	},
}

type TypeInfo struct {
	TypeName string
	VariableName string
	BuiltIn bool
	Pointer bool
}

func (i TypeInfo) ArgumentName() string {
	var pointer string
	if i.Pointer {
		pointer = "*"
	}

	var pkgName string
	if !i.BuiltIn {
		pkgName = "ast."
	}

	return fmt.Sprintf("%s%s%s", pointer, pkgName, i.TypeName)
}

func main() {

}

var walkMethodTpl = template.Must(template.New("walkType").Parse(`
func (w *Walker) walk{{.TypeName}}(ctx *Context, doc {{.ArgumentName}}) {

}
`))
