package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
)

const (
	astFileName   = "ast.go"
	listsFileName = "lists.go"
)

type File struct {
	Package string
	Nodes   map[string]Node
}

type Node struct {
	IsPointer bool
	Fields    map[string]Field
}

type Field struct {
	IsListType bool
	IsKindType bool
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

	astTypes, err := readTypeNames(astFile)
	if err != nil {
		log.Fatal(err)
	}

	listFile, err := readFile(listsFileName)
	if err != nil {
		log.Fatal(err)
	}

	listTypes, err := readTypeNames(listFile)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(astTypes)
	spew.Dump(listTypes)

	spew.Dump(astFile)
}

func readFile(fileName string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), path.Join(os.Args[1], fileName), nil, parser.ParseComments)
}

func readTypeNames(file *ast.File) ([]string, error) {
	var types []string

	for _, decl := range file.Decls {
		gdec, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gdec.Tok != token.TYPE {
			continue
		}

		for _, spec := range gdec.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			types = append(types, tspec.Name.Name)
		}
	}

	return types, nil
}
