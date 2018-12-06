package walker

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
	"github.com/davecgh/go-spew/spew"
)

// TODO: check output.
func Generate(w io.Writer, packageName string, noImports bool, s *goast.SymbolTable) {
	// Header and package name
	fmt.Fprintf(os.Stdout, strings.TrimSpace(header))
	fmt.Fprintf(os.Stdout, "\npackage %s\n", packageName)

	if !noImports {
		fmt.Fprintf(os.Stdout, "%s", imports)
	}

	var typeNames []string
	for tn := range s.Structs {
		typeNames = append(typeNames, tn)
	}

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
			IsListType: isStructListType(tn, s.Structs[tn]),
		}

		// Add field information.
		for k, f := range s.Structs[tn].Fields {
			if _, ok := s.Structs[f.TypeName]; ok {
				td.Fields = append(td.Fields, field{
					Name: k,
					Type: f,
				})
			}
		}

		// If we have a field called "Kind", then we need to generate a switch statement too.
		if f, ok := s.Structs[tn].Fields["Kind"]; ok {
			td.IsSwitcher = true
			td.Consts = s.Consts[f.TypeName]

			//for _, c := range td.Consts {
			//	if c.Field == "self" {
			//		tds = append(tds, templateData{})
			//	}
			//}
		}

		tds = append(tds, td)
	}

	for i, ltd := range tds {
		if !ltd.IsListType {
			continue
		}

		tn := strings.TrimRight(ltd.TypeName, "s")

		// Find the non-list type for this list type.
		for _, td := range tds {
			if td.TypeName != tn {
				continue
			}

			ltd.NodeType = &td
			break
		}

		tds[i] = ltd
	}

	err := walkerTypeTmpl.Execute(w, tds)
	if err != nil {
		log.Fatal(err)
	}

	templates := []*template.Template{
		eventHandlersTmpl,
		walkFnTmpl,
	}

	spew.Fdump(os.Stderr, tds)

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
func isFieldPointer(s *goast.SymbolTable, tn string) bool {
	for _, strc := range s.Structs {
		for _, f := range strc.Fields {
			if f.TypeName == tn {
				return f.IsPointer
			}
		}
	}

	return false
}

// isFieldArray checks the symbol table for references of the given type name on fields of other
// types, returning true if the given type name is ever used as an array.
func isFieldArray(s *goast.SymbolTable, tn string) bool {
	for _, strc := range s.Structs {
		for _, f := range strc.Fields {
			if f.TypeName == tn {
				return f.IsArray
			}
		}
	}

	return false
}

// isStructListType ...
func isStructListType(tn string, str goast.Struct) bool {
	if len(str.Fields) != 3 {
		return false
	}

	// A list type only has these 3 fields.
	_, hasDataField := str.Fields["Data"]
	next, hasNextField := str.Fields["next"]
	pos, hasPosField := str.Fields["pos"]

	// The type of the "next" field should match the name of this list type.
	hasCorrectNextType := next.TypeName == tn
	hasCorrectPosType := pos.TypeName == "int"

	return hasDataField && hasNextField && hasPosField && hasCorrectNextType && hasCorrectPosType
}
