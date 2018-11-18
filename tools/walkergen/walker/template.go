package walker

import (
	"strings"
	"text/template"
	"unicode"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
)

// walkFnTmpl is the template for generating a walker function.
var walkFnTmpl = template.Must(template.New("walkFnTmpl").Parse(`
// walk{{.TypeName}} ...
func (w *Walker) walk{{.TypeName}}(ctx *Context, {{.ShortTypeName}} {{if .IsArray}}[]{{end}}{{if .IsPointer}}*{{end}}ast.{{.TypeName}}) {
	w.On{{.TypeName}}Enter(ctx, {{.ShortTypeName}})
	{{if .IsListType}}{{.ShortTypeName}}.ForEach(func({{.NodeType.ShortTypeName}} ast.{{.NodeType.TypeName}}, i int) {
		w.walk{{.NodeType.TypeName}}(ctx, {{.NodeType.ShortTypeName}})
	})
	{{$parent := .}}{{else if .Consts}}switch {{.ShortTypeName}}.Kind {
	{{range .Consts}}case ast.{{.Name}}:
		w.walk{{.Name}}(ctx, {{$.ShortTypeName}}.{{.Field}})
	{{end}}}
	{{end}}w.On{{.TypeName}}Leave(ctx, {{.ShortTypeName}})
}
`))

// walkTemplateData contains data and methods necessary to generate a walk function.
type walkTemplateData struct {
	goast.Type
	// Consts contains the constants related to this type.
	Consts []goast.Const
	// NodeType is populated if this type is a list type (i.e. IsListType = true).
	NodeType *walkTemplateData
	// IsListType is true if this type is a linked list type.
	IsListType bool
	// IsSwitcher is a type that we need to generate a switch statement for.
	IsSwitcher bool
}

// ShortTypeName returns an abridged version of the embedded type's name.
func (ttd walkTemplateData) ShortTypeName() string {
	stn := strings.Map(abridger, ttd.TypeName)
	if ttd.IsListType {
		return stn + "s"
	}

	return stn
}

// abridger is a strings.Map function that is used to return a variable name from a type name that
// makes sense by taking each capital letter from the given string and converting them to lowercase.
// The input string should start with a capital letter.
func abridger(r rune) rune {
	if unicode.IsUpper(r) {
		return unicode.ToLower(r)
	}
	return -1
}
