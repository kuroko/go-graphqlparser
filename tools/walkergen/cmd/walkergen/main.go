package main

import (
	"flag"
	"log"
	"os"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
	"github.com/bucketd/go-graphqlparser/tools/walkergen/walker"
)

const (
	astFileName   = "ast.go"
	listsFileName = "lists.go"
)

func main() {
	var astPath string
	var packageName string
	var noImports bool

	flag.StringVar(&astPath, "ast-path", "", "The path to the AST package on the filesystem.")
	flag.StringVar(&packageName, "package", "", "The package name to use in the generated code.")
	flag.BoolVar(&noImports, "no-imports", false, "Use this flag to exclude imports.")
	flag.Parse()

	if astPath == "" {
		log.Println("you must specify a path to the ast package")
		os.Exit(1)
	}

	if packageName == "" {
		log.Println("you must specify a package name for the generated code")
		os.Exit(1)
	}

	astFile, err := goast.ReadFile(astPath, astFileName)
	if err != nil {
		log.Fatal(err)
	}

	astSymbols, err := goast.CreateSymbolTable(astFile)
	if err != nil {
		log.Fatal(err)
	}

	listFile, err := goast.ReadFile(astPath, listsFileName)
	if err != nil {
		log.Fatal(err)
	}

	listSymbols, err := goast.CreateSymbolTable(listFile)
	if err != nil {
		log.Fatal(err)
	}

	symbols := goast.Symbols{
		AST:  astSymbols,
		List: listSymbols,
	}

	// Output walker.
	walker.Generate(os.Stdout, packageName, noImports, &symbols)
}
