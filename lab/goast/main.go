package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	Name string
	// TODO(elliot): Populate this.
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

	_ = astSymbols
	_ = listSymbols

	spew.Dump(astSymbols)
	//spew.Dump(listSymbols)
	//spew.Dump(astFile)
}

// readFile ...
func readFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(os.Args[1], fileName), nil, parser.ParseComments)
}

// createSymbolTable ...
func createSymbolTable(file *ast.File) (SymbolTable, error) {
	symbols := NewSymbolTable()

	for _, decl := range file.Decls {
		gdecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch gdecl.Tok {
		case token.CONST:
			processConstDeclaration(&symbols, gdecl)
		case token.TYPE:
			processTypeDeclaration(&symbols, gdecl)
		default:
			continue
		}
	}

	return symbols, nil
}

// processConstDeclaration ...
func processConstDeclaration(symbols *SymbolTable, decl *ast.GenDecl) {
	if len(decl.Specs) < 1 {
		return
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

			processValueSpec(symbols, t, v)
		}
	}
}

// processTypeDeclaration ...
func processTypeDeclaration(symbols *SymbolTable, decl *ast.GenDecl) {
	// TODO(elliot): Move into some kind of function that extracts this stuff into a struct or
	// something? Then just make a simple if statement here.
	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if strings.Contains(comment.Text, "@wg:ignore") {
				return
			}
		}
	}

	for _, spec := range decl.Specs {
		switch v := spec.(type) {
		case *ast.TypeSpec:
			processTypeSpec(symbols, v)
		}
	}
}

// processTypeSpec ...
func processTypeSpec(symbols *SymbolTable, tspec *ast.TypeSpec) {
	// Get the name of the type.
	name := tspec.Name.Name

	// If the type doesn't already exist in the symbol table, then we need to add it.
	if _, ok := symbols.Structs[name]; !ok {
		symbols.Structs[name] = NewStruct()
	}

	switch v := tspec.Type.(type) {
	case *ast.StructType:
		processStructType(symbols, name, v)
	}
}

func processValueSpec(symbols *SymbolTable, t Type, vspec *ast.ValueSpec) {
	var consts []Const
	for _, name := range vspec.Names {
		consts = append(consts, Const{
			Name: name.Name,
		})
	}

	symbols.Consts[t.TypeName] = append(symbols.Consts[t.TypeName], consts...)
}

// processStructType ...
func processStructType(symbols *SymbolTable, name string, st *ast.StructType) {
	// For a struct switch.
	for _, field := range st.Fields.List {
		for _, fieldIdent := range field.Names {
			if t, ok := processExpr(field.Type); ok {
				symbols.Structs[name].Fields[fieldIdent.Name] = t
			}
		}
	}
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
