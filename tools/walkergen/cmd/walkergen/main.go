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

	flag.StringVar(&astPath, "ast-path", "", "The path to the AST package on the filesystem.")
	flag.StringVar(&packageName, "package", "", "The package name to use in the generated code.")
	flag.Parse()

	if astPath == "" {
		log.Fatalln("walkergen: you must specify a path to the ast package")
	}

	if packageName == "" {
		log.Fatalln("walkergen: you must specify a package name for the generated code")
	}

	astFile, err := goast.ReadFile(astPath, astFileName)
	if err != nil {
		log.Fatal(err)
	}

	listFile, err := goast.ReadFile(astPath, listsFileName)
	if err != nil {
		log.Fatal(err)
	}

	symbols := goast.NewSymbolTable()

	err = goast.PopulateSymbolTable(astFile, &symbols)
	if err != nil {
		log.Fatal(err)
	}

	err = goast.PopulateSymbolTable(listFile, &symbols)
	if err != nil {
		log.Fatal(err)
	}

	// Output walker.
	walker.Generate(os.Stdout, packageName, symbols)
}
