package goast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

// Symbols is a collection of symbol tables for walker generation from multiple files, just here for
// convenience when passing the symbol tables around.
type Symbols struct {
	AST  SymbolTable
	List SymbolTable
}

// SymbolTable takes the information we need from the Go AST, and represents it in an easy to access
// structure for the walker generator.
type SymbolTable struct {
	Package string
	Consts  map[string]Consts
	Structs map[string]Struct
}

// NewSymbolTable returns a new SymbolTable value with the maps on it initialised.
func NewSymbolTable() SymbolTable {
	return SymbolTable{
		Consts:  make(map[string]Consts),
		Structs: make(map[string]Struct),
	}
}

// Annotations contains any annotations that are relevant to a section of Go code being processed.
type Annotations struct {
	Field  string
	Ignore bool
}

// Consts is a slice of Const, with some additional methods.
type Consts []Const

// HasNonSelfConsts returns true if any of these Consts have a Field that isn't set to self.
func (cs Consts) HasNonSelfConsts() bool {
	for _, c := range cs {
		if c.Field != "self" {
			return true
		}
	}

	return false
}

// Const represents the information we need from the Go AST for constants.
type Const struct {
	Name  string
	Field string
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
	IsArray   bool
	IsPointer bool
}

// ReadFile ...
func ReadFile(filePath, fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(filePath, fileName), nil, parser.ParseComments)
}

// CreateSymbolTable ...
func CreateSymbolTable(file *ast.File) (SymbolTable, error) {
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

// readAnnotations ...
func readAnnotations(cg *ast.CommentGroup, annotations Annotations) (Annotations, error) {
	if cg != nil {
		for _, comment := range cg.List {
			switch {
			case strings.Contains(comment.Text, "@wg:ignore"):
				annotations.Ignore = true
			case strings.Contains(comment.Text, "@wg:field"):
				f := strings.Split(comment.Text, " ")
				if len(f) < 3 {
					return annotations, fmt.Errorf("wg metadata '%v' invalid", comment.Text)
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
	annotations, err := readAnnotations(vspec.Doc, annotations)
	if err != nil {
		return "", err
	}

	if annotations.Field != "" {
		return annotations.Field, nil
	}

	// Otherwise, construct a sane field name out of the constant's name.
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
