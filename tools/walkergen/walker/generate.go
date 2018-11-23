package walker

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
)

// TODO: check output.
func Generate(w io.Writer, packageName string, noImports bool, s *goast.Symbols) {
	// Header and package name
	fmt.Fprintf(os.Stdout, strings.TrimSpace(header))
	fmt.Fprintf(os.Stdout, "\npackage %s\n", packageName)

	if !noImports {
		fmt.Fprintf(os.Stdout, "%s", imports)
	}

	var typeNames []string
	for tn := range s.AST.Structs {
		typeNames = append(typeNames, tn)
	}

	sort.Strings(typeNames)

	// Construct our template data, based on the symbol table.
	// For each struct type name...
	var tds []templateData
	for _, tn := range typeNames {
		td := templateData{
			Type: goast.Type{
				TypeName:  tn,
				IsArray:   isFieldArray(s, tn),
				IsPointer: isFieldPointer(s, tn),
			},
		}

		listName := tn + "s"

		// Does a list type exist for this type?
		if _, ok := s.List.Structs[listName]; ok {
			tds = append(tds,
				templateData{
					NodeType: &td,
					Type: goast.Type{
						TypeName:  listName,
						IsPointer: true,
					},
					IsListType: true,
				},
			)
		}

		// If we have a field called "Kind", then we need to generate a switch statement too.
		if f, ok := s.AST.Structs[tn].Fields["Kind"]; ok {
			td.IsSwitcher = true
			td.Consts = s.AST.Consts[f.TypeName]
		}

		tds = append(tds, td)
	}

	err := walkerTypeTmpl.Execute(w, tds)
	if err != nil {
		log.Fatal(err)
	}

	templates := []*template.Template{
		eventHandlersTmpl,
		walkFnTmpl,
	}

	for _, td := range tds {
		for _, tmpl := range templates {
			err := tmpl.Execute(w, td)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// isFieldPointer checks the symbol table for references of the given type name on fields of other
// types, returning true if the given type name is ever used as a pointer.
func isFieldPointer(s *goast.Symbols, tn string) bool {
	for _, strc := range s.AST.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsPointer
		}
	}

	return false
}

// isFieldArray checks the symbol table for references of the given type name on fields of other
// types, returning true if the given type name is ever used as an array.
func isFieldArray(s *goast.Symbols, tn string) bool {
	for _, strc := range s.AST.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsArray
		}
	}

	return false
}
