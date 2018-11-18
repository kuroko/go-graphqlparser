package goast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
func ReadFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(os.Args[1], fileName), nil, parser.ParseComments)
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
