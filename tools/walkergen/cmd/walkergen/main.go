package main

import (
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
	if len(os.Args) < 2 {
		log.Println("you must specify a path to the ast package")
		os.Exit(1)
	}

	astFile, err := goast.ReadFile(astFileName)
	if err != nil {
		log.Fatal(err)
	}

	astSymbols, err := goast.CreateSymbolTable(astFile)
	if err != nil {
		log.Fatal(err)
	}

	listFile, err := goast.ReadFile(listsFileName)
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
	walker.Generate(os.Stdout, &symbols)
}
