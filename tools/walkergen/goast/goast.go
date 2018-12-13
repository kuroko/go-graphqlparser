package goast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

// SymbolTable takes the information we need from the Go AST, and represents it in an easy to access
// structure for the walker generator.
type SymbolTable struct {
	Package string
	Consts  map[string][]Const
	Structs map[string]Struct
}

// NewSymbolTable returns a new SymbolTable value with the maps on it initialised.
func NewSymbolTable() SymbolTable {
	return SymbolTable{
		Consts:  make(map[string][]Const),
		Structs: make(map[string]Struct),
	}
}

// Annotations contains any annotations that are relevant to a section of Go code being processed.
type Annotations struct {
	Field   string
	Ignore  bool
	OnKinds []string
}

// Const represents the information we need from the Go AST for constants.
// TODO: This isn't really a "Const", it's more specifically a "Kind", isn't it? This looks more
// like some template data.
type Const struct {
	Name  string
	Field string
}

// SelfName ...
func (c Const) SelfName() string {
	n, _ := constructSaneFieldName(c.Name)
	return n
}

// Struct represents the information we need from the Go AST for structs.
type Struct struct {
	Fields map[string]Type
}

// NewStruct returns a new Struct value with the map on it initialised.
func NewStruct() Struct {
	return Struct{
		Fields: make(map[string]Type),
	}
}

// Type represents the information we need from the Go AST for types.
type Type struct {
	TypeName  string
	OnKinds   []string
	IsArray   bool
	IsPointer bool
}

// ReadFile ...
func ReadFile(filePath, fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(filePath, fileName), nil, parser.ParseComments)
}

// PopulateSymbolTable ...
func PopulateSymbolTable(file *ast.File, symbols *SymbolTable) error {
	var err error

	for _, decl := range file.Decls {
		gdecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch gdecl.Tok {
		case token.CONST:
			err = processConstDeclaration(symbols, gdecl)
		case token.TYPE:
			err = processTypeDeclaration(symbols, gdecl)
		default:
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// readAnnotations ...
func readAnnotations(cg *ast.CommentGroup, annotations Annotations) (Annotations, error) {
	if cg != nil {
		for _, comment := range cg.List {
			switch {
			case strings.Contains(comment.Text, "@wg:on_kinds"):
				f := strings.Split(comment.Text, " ")
				if len(f) < 3 {
					return annotations, fmt.Errorf("wg metadata %q invalid", comment.Text)
				}

				annotations.OnKinds = strings.Split(f[2], ",")

			case strings.Contains(comment.Text, "@wg:ignore"):
				annotations.Ignore = true

			case strings.Contains(comment.Text, "@wg:field"):
				f := strings.Split(comment.Text, " ")
				if len(f) < 3 {
					return annotations, fmt.Errorf("wg metadata %q invalid", comment.Text)
				}

				annotations.Field = f[2]
			}
		}
	}

	return annotations, nil
}

// processConstDeclaration ...
func processConstDeclaration(symbols *SymbolTable, decl *ast.GenDecl) error {
	if len(decl.Specs) < 1 {
		return nil
	}

	annotations, err := readAnnotations(decl.Doc, Annotations{})
	if err != nil {
		return err
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

			err := processValueSpec(symbols, t, v, annotations)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processTypeDeclaration ...
func processTypeDeclaration(symbols *SymbolTable, decl *ast.GenDecl) error {
	annotations, err := readAnnotations(decl.Doc, Annotations{})
	if err != nil {
		return err
	}

	if annotations.Ignore {
		return nil
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

	switch v := tspec.Type.(type) {
	case *ast.StructType:
		// If the type doesn't already exist in the symbol table, then we need to add it.
		if _, ok := symbols.Structs[name]; !ok {
			symbols.Structs[name] = NewStruct()
		}

		err := processStructType(symbols, name, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func processValueSpec(symbols *SymbolTable, t Type, vspec *ast.ValueSpec, annotations Annotations) error {
	var consts []Const

	annotations, err := readAnnotations(vspec.Doc, annotations)
	if err != nil {
		return err
	}

	field, err := processFieldName(vspec, annotations)
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

func processFieldName(vspec *ast.ValueSpec, annotations Annotations) (string, error) {
	if annotations.Field != "" {
		return annotations.Field, nil
	}

	return constructSaneFieldName(vspec.Names[0].Name)
}

// constructSaneFieldName creates a nicer name from the const name.
func constructSaneFieldName(nonSane string) (string, error) {
	f := strings.Split(nonSane, "Kind")
	if len(f) < 2 {
		return "", fmt.Errorf("name %v not properly formatted, should have Kind flanked by a word either side", nonSane)
	}
	return f[1] + f[0], nil
}

// processStructType ...
func processStructType(symbols *SymbolTable, name string, st *ast.StructType) error {
	// For a struct switch.
	for _, field := range st.Fields.List {
		annotations, err := readAnnotations(field.Doc, Annotations{})
		if err != nil {
			return err
		}

		for _, fieldIdent := range field.Names {
			if t, ok := processExpr(field.Type); ok {
				t.OnKinds = annotations.OnKinds
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
