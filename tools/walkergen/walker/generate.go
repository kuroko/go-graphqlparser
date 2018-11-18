package walker

import (
	"io"
	"sort"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
)

// TODO: check output.
// TODO: clean code.
func Generate(w io.Writer, s *goast.Symbols) {
	var typeNames []string
	for tn := range s.AST.Structs {
		typeNames = append(typeNames, tn)
	}

	sort.Strings(typeNames)

	// For each struct type name...
	var wtds []walkTemplateData
	for _, tn := range typeNames {
		wtd := walkTemplateData{
			Type: goast.Type{
				TypeName:  tn,
				IsArray:   isFieldArray(s, tn),
				IsPointer: isFieldPointer(s, tn),
			},
		}

		listName := tn + "s"

		// Does a list type exist for this type?
		if _, ok := s.List.Structs[listName]; ok {
			wtds = append(wtds,
				walkTemplateData{
					NodeType: &wtd,
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
			wtd.IsSwitcher = true
			wtd.Consts = s.AST.Consts[f.TypeName]
		}

		wtds = append(wtds, wtd)
	}

	for _, baz := range wtds {
		walkFnTmpl.Execute(w, baz)
	}
}

func isFieldPointer(s *goast.Symbols, tn string) bool {
	for _, strc := range s.AST.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsPointer
		}
	}
	return false
}

func isFieldArray(s *goast.Symbols, tn string) bool {
	for _, strc := range s.AST.Structs {
		if fld, ok := strc.Fields[tn]; ok {
			return fld.IsArray
		}
	}
	return false
}
