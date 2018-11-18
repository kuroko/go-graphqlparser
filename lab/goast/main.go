package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"
)

const (
	astFileName   = "ast.go"
	listsFileName = "lists.go"
)

type SymbolTable struct {
	Package string
	Consts  map[string][]Const
	Structs map[string]Struct
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		Consts:  make(map[string][]Const),
		Structs: make(map[string]Struct),
	}
}

type Const struct {
	Name  string
	Field string
}

type Struct struct {
	Fields map[string]Type
}

func NewStruct() Struct {
	return Struct{
		Fields: make(map[string]Type),
	}
}

type Type struct {
	TypeName  string
	IsArray   bool
	IsPointer bool
}

type Symbols struct {
	ast, list SymbolTable
}

func main() {
	if len(os.Args) < 2 {
		log.Println("you must specify a path to the ast package")
		os.Exit(1)
	}

	astFile, err := readFile(astFileName)
	if err != nil {
		log.Fatal(err)
	}

	astSymbols, err := createSymbolTable(astFile)
	if err != nil {
		log.Fatal(err)
	}

	listFile, err := readFile(listsFileName)
	if err != nil {
		log.Fatal(err)
	}

	listSymbols, err := createSymbolTable(listFile)
	if err != nil {
		log.Fatal(err)
	}

	symbols := Symbols{
		ast:  astSymbols,
		list: listSymbols,
	}

	gen(&symbols)
}

// readFile ...
func readFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(os.Args[1], fileName), nil, parser.ParseComments)
}

// createSymbolTable ...
func createSymbolTable(file *ast.File) (SymbolTable, error) {
	symbols := NewSymbolTable()
	var err error

	for _, decl := range file.Decls {
		gdecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch gdecl.Tok {
		case token.CONST:
			err = processConstDeclaration(&symbols, gdecl)
		case token.TYPE:
			err = processTypeDeclaration(&symbols, gdecl)
		default:
			continue
		}
		if err != nil {
			return symbols, err
		}
	}

	return symbols, nil
}

// processConstDeclaration ...
func processConstDeclaration(symbols *SymbolTable, decl *ast.GenDecl) error {
	if len(decl.Specs) < 1 {
		return nil
	}

	var t Type
	var ok bool

	for i, spec := range decl.Specs {
		switch v := spec.(type) {
		case *ast.ValueSpec:
			if i == 0 {
				t, ok = processExpr(v.Type)
				if !ok {
					continue
				}
			}
			err := processValueSpec(symbols, t, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// processTypeDeclaration ...
func processTypeDeclaration(symbols *SymbolTable, decl *ast.GenDecl) error {
	// TODO(elliot): Move into some kind of function that extracts this stuff into a struct or
	// something? Then just make a simple if statement here.
	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if strings.Contains(comment.Text, "@wg:ignore") {
				return nil
			}
		}
	}

	for _, spec := range decl.Specs {
		switch v := spec.(type) {
		case *ast.TypeSpec:
			err := processTypeSpec(symbols, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// processTypeSpec ...
func processTypeSpec(symbols *SymbolTable, tspec *ast.TypeSpec) error {
	// Get the name of the type.
	name := tspec.Name.Name

	// If the type doesn't already exist in the symbol table, then we need to add it.
	if _, ok := symbols.Structs[name]; !ok {
		symbols.Structs[name] = NewStruct()
	}

	switch v := tspec.Type.(type) {
	case *ast.StructType:
		err := processStructType(symbols, name, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func processValueSpec(symbols *SymbolTable, t Type, vspec *ast.ValueSpec) error {
	var consts []Const

	field, err := processFieldName(vspec)
	if err != nil {
		return err
	}

	for _, name := range vspec.Names {
		consts = append(consts, Const{
			Name:  name.Name,
			Field: field,
		})
	}

	symbols.Consts[t.TypeName] = append(symbols.Consts[t.TypeName], consts...)
	return nil
}

func processFieldName(vspec *ast.ValueSpec) (string, error) {
	if vspec.Doc != nil {
		for _, comment := range vspec.Doc.List {
			if strings.Contains(comment.Text, "@wg:field") {
				f := strings.Split(comment.Text, " ")
				if len(f) < 3 {
					return "", fmt.Errorf("wg metadata '%v' invalid", comment.Text)
				}
				return f[2], nil
			}
		}
	}

	f := strings.Split(vspec.Names[0].Name, "Kind")
	if len(f) < 2 {
		return "", fmt.Errorf("name %v not properly formatted, should have Kind flanked by a word either side", vspec.Names[0].Name)
	}
	return f[1] + f[0], nil
}

// processStructType ...
func processStructType(symbols *SymbolTable, name string, st *ast.StructType) error {
	// For a struct switch.
	for _, field := range st.Fields.List {
		for _, fieldIdent := range field.Names {
			if t, ok := processExpr(field.Type); ok {
				symbols.Structs[name].Fields[fieldIdent.Name] = t
			}
		}
	}
	return nil
}

// processExpr ...
func processExpr(expr ast.Expr) (Type, bool) {
	var t Type
	var ident *ast.Ident
	var isArray, isPointer, ok bool

	switch v := expr.(type) {
	case *ast.Ident:
		ident = v
	case *ast.ArrayType:
		isArray = true
		if ident, ok = v.Elt.(*ast.Ident); !ok {
			return t, false
		}
	case *ast.StarExpr:
		isPointer = true
		if ident, ok = v.X.(*ast.Ident); !ok {
			return t, false
		}
	default:
		return t, false
	}

	return Type{
		TypeName:  ident.Name,
		IsArray:   isArray,
		IsPointer: isPointer,
	}, true
}

func fgen(w io.Writer, s *Symbols) {
	generate(w, s)
}

func gen(s *Symbols) {
	fgen(os.Stdout, s)
}

func sgen(s *Symbols) string {
	buf := bytes.Buffer{}

	fgen(&buf, s)

	return buf.String()
}

// TODO: check output.
// TODO: clean code.
// TODO: kind types still coming through on struct data.
func generate(w io.Writer, s *Symbols) {
	var foo []string
	for tn := range s.ast.Structs {
		foo = append(foo, tn)
	}
	sort.Strings(foo)

	var bar []walkTemplateData
	for _, tn := range foo {
		litn := tn + "s"
		if _, hasList := s.list.Structs[litn]; hasList {
			bar = append(bar,
				walkTemplateData{
					Type: Type{
						TypeName:  litn,
						IsPointer: true,
					},
					ListType: true,
				},
			)
		}

		wtd := walkTemplateData{
			Type: Type{
				TypeName:  tn,
				IsArray:   isFieldArray(s, tn),
				IsPointer: isFieldPointer(s, tn),
			},
		}

		if f, ok := s.ast.Structs[tn].Fields["Kind"]; ok {
			wtd.IsSwitcher = true
			wtd.Consts = s.ast.Consts[f.TypeName]
		}

		bar = append(bar, wtd)
	}

	for _, baz := range bar {
		walkTemplate.Execute(w, baz)
	}
}

func isFieldPointer(s *Symbols, tn string) bool {
	for _, strc := range s.ast.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsPointer
		}
	}
	return false
}

func isFieldArray(s *Symbols, tn string) bool {
	for _, strc := range s.ast.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsArray
		}
	}
	return false
}

type walkTemplateData struct {
	Type
	ListType   bool
	IsSwitcher bool
	Consts     []Const
}

func (ttd walkTemplateData) ShortTN() string {
	stn := strings.Map(abridger, ttd.TypeName)
	if ttd.ListType {
		return stn + "s"
	}
	return stn
}

func allButLastLetter(s string) string {
	return s[:len(s)-1]
}

func (ttd walkTemplateData) TypeNameMinusS() string {
	return allButLastLetter(ttd.TypeName)
}

func (ttd walkTemplateData) ShortTNMinusS() string {
	return allButLastLetter(ttd.ShortTN())
}

func abridger(r rune) rune {
	if unicode.IsUpper(r) {
		return unicode.ToLower(r)
	}
	return -1
}

var walkTemplate = template.Must(template.New("walkTemplate").Parse(`
// walk{{.TypeName}} ...
func (w *Walker) walk{{.TypeName}}(ctx *Context, {{.ShortTN}} {{if .IsArray}}[]{{end}}{{if .IsPointer}}*{{end}}ast.{{.TypeName}}) {
	w.On{{.TypeName}}Enter(ctx, {{.ShortTN}})
	{{if .ListType}}{{.ShortTN}}.ForEach(func({{.ShortTNMinusS}} ast.{{.TypeNameMinusS}}, i int) {
		w.walk{{.TypeNameMinusS}}(ctx, {{.ShortTNMinusS}})
	})
	{{$parent := .}}{{else if .Consts}}switch {{.ShortTN}}.Kind {
	{{range .Consts}}case ast.{{.Name}}:
		w.walk{{.Name}}(ctx, {{$.ShortTN}}.{{.Field}})
	{{end}}}
	{{end}}w.On{{.TypeName}}Leave(ctx, {{.ShortTN}})
}
`))
