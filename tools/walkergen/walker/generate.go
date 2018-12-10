package walker

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/bucketd/go-graphqlparser/tools/walkergen/goast"
	"github.com/davecgh/go-spew/spew"
)

// Generate ...
func Generate(w io.Writer, packageName string, noImports bool, st goast.SymbolTable) {
	// First pass, get walker types for types that are actually defined in the symbol table. The
	// actual list of walker types will grow once we add the "kind" types later.
	wts := buildBaseTypes(st)

	// Second pass, now we have all base types, with as much information populated as possible, we
	// need to attach those types to fields, and kinds.
	wts = hydrateAllTypes(wts)

	// Third pass, this one adds in the types that aren't actually in the AST, as in, if we have a
	// "self" kind type, we add all of those types too, so that we do actually generate walker
	// functions for them.
	wts = injectSelfKindTypes(wts)

	// Then sort them all into order so that we end up with a consistent result.
	sort.Slice(wts, func(i, j int) bool {
		return wts[i].FuncName < wts[j].FuncName
	})

	// Avoid printing this into generated file.
	spew.Fdump(os.Stderr, wts)

	// Header and package name
	fmt.Fprintf(os.Stdout, strings.TrimSpace(header))
	fmt.Fprintf(os.Stdout, "\npackage %s\n", packageName)

	if !noImports {
		fmt.Fprintf(os.Stdout, "%s", imports)
	}

	err := walkerTypeTmpl.Execute(w, wts)
	if err != nil {
		log.Fatal(err)
	}

	for _, wt := range wts {
		err := eventHandlersTmpl.Execute(w, wt)
		if err != nil {
			log.Fatal(err)
		}

		err = walkerFnTmpl(w, wt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// buildBaseTypes ...
func buildBaseTypes(st goast.SymbolTable) []walkerType {
	var wts []walkerType

	var tns []string
	for tn := range st.Structs {
		tns = append(tns, tn)
	}

	sort.Strings(tns)

	for _, tn := range tns {
		wt := walkerType{}
		wt.TypeName = tn
		wt.FuncName = tn
		wt.Fields = buildBaseTypeFields(st.Structs[tn].Fields)
		wt.Kinds = buildBaseTypeKinds(st.Consts[fmt.Sprintf("%sKind", tn)])
		wt.IsAlwaysPointer = isTypeAlwaysPointer(st, tn)
		wt.IsLinkedList = isTypeLinkedList(tn, st.Structs[tn])
		wt.ShortTypeName = buildTypeShortName(wt)
		wt.FullTypeName = buildTypeFullName(wt)

		wts = append(wts, wt)
	}

	return wts
}

// hydrateAllTypes ...
func hydrateAllTypes(wts []walkerType) []walkerType {
	// Firstly, hydrate all fields, now that we can.
	for i, wt := range wts {
		wts[i] = hydrateAllFields(wts, wt)
	}

	// Then hydrate the kinds, with type and field information.
	for i, wt := range wts {
		wts[i] = hydrateAllKinds(wts, wt)
	}

	// We also want all linked list types to be attached to their types.
	for i, wt := range wts {
		wts[i] = hydrateLinkedListTypes(wts, wt)
	}

	return wts
}

// buildBaseTypeFields ...
func buildBaseTypeFields(fields map[string]goast.Type) []walkerTypeField {
	var wtfs []walkerTypeField

	var fns []string
	for fn := range fields {
		fns = append(fns, fn)
	}

	sort.Strings(fns)

	for _, fn := range fns {
		wtf := walkerTypeField{}
		wtf.Name = fn
		wtf.IsPointerType = fields[fn].IsPointer
		wtf.IsSliceType = fields[fn].IsArray
		wtf.typeName = fields[fn].TypeName
		// wtf.Type is set once we've gathered all types.

		wtfs = append(wtfs, wtf)
	}

	return wtfs
}

// hydrateAllFields ...
func hydrateAllFields(wts []walkerType, wt walkerType) walkerType {
	for i, fld := range wt.Fields {
		for _, wt := range wts {
			if wt.TypeName == fld.typeName {
				// We must first set the referenced type's fields and kinds to nil, otherwise we can
				// end up with huge sub-structures. At this level, we're simply not interested in
				// this information too.
				wt.Fields = nil
				wt.Kinds = nil

				fld.Type = wt
				fld.isASTType = true
				break
			}
		}

		wt.Fields[i] = fld
	}

	// Second pass, to clear out non-AST types.
	for i := len(wt.Fields) - 1; i >= 0; i-- {
		fld := wt.Fields[i]
		if fld.isASTType {
			continue
		}

		wt.Fields = append(wt.Fields[:i], wt.Fields[i+1:]...)
	}

	return wt
}

// buildBaseTypeKinds ...
func buildBaseTypeKinds(consts []goast.Const) []walkerTypeKind {
	var wtks []walkerTypeKind

	sort.Slice(consts, func(i, j int) bool {
		return consts[i].Name < consts[j].Name
	})

	for _, c := range consts {
		wtk := walkerTypeKind{}
		wtk.ConstName = c.Name
		wtk.IsSelf = c.Field == "self"
		wtk.fieldName = c.Field
		// wtk.Type is set once we've gathered all types.
		// wtk.Field is set once we've gathered all types and fields.

		wtks = append(wtks, wtk)
	}

	return wtks
}

// hydrateAllKinds ...
func hydrateAllKinds(wts []walkerType, wt walkerType) walkerType {
	for i, knd := range wt.Kinds {
		kwt := wt

		for _, fld := range wt.Fields {
			if fld.Name == knd.fieldName {
				knd.Field = &fld
				break
			}
		}

		if knd.IsSelf {
			kwt.FuncName = buildTypeKindFuncName(knd.ConstName)
		} else {
			kwt = knd.Field.Type
		}

		// We must first set the referenced type's fields and kinds to nil, otherwise we can
		// end up with huge sub-structures. At this level, we're simply not interested in
		// this information too.
		kwt.Fields = nil
		kwt.Kinds = nil

		knd.Type = kwt

		wt.Kinds[i] = knd
	}

	return wt
}

// hydrateLinkedListTypes ...
func hydrateLinkedListTypes(wts []walkerType, wt walkerType) walkerType {
	if !wt.IsLinkedList {
		return wt
	}

	nodeTypeName := strings.TrimRight(wt.TypeName, "s")
	for _, nwt := range wts {
		rnwt := nwt

		if nwt.FuncName != nodeTypeName {
			continue
		}

		// We must first set the referenced type's fields and kinds to nil, otherwise we can
		// end up with huge sub-structures. At this level, we're simply not interested in
		// this information too.
		nwt.Fields = nil
		nwt.Kinds = nil

		wt.LinkedListType = &rnwt
	}

	return wt
}

// injectSelfKindTypes ...
func injectSelfKindTypes(wts []walkerType) []walkerType {
	for _, wt := range wts {
		if len(wt.Kinds) == 0 {
			continue
		}

		for _, knd := range wt.Kinds {
			if !knd.IsSelf {
				continue
			}

			// Find the real type info, but update the FuncName, then remove kinds, as we don't want
			// to recurse!
			kwt := findWalkerTypeByName(wts, knd.Type.TypeName)
			kwt.FuncName = buildTypeKindFuncName(knd.ConstName)
			kwt.Kinds = nil

			wts = append(wts, kwt)
		}
	}

	return wts
}

// findWalkerTypeByName ...
func findWalkerTypeByName(wts []walkerType, name string) walkerType {
	var wt walkerType
	for _, wt := range wts {
		if wt.TypeName == name {
			return wt
		}
	}

	return wt
}

// buildTypeShortName ...
func buildTypeShortName(wt walkerType) string {
	stn := strings.Map(abridger, wt.TypeName)
	if wt.IsLinkedList {
		return stn + "s"
	}

	return stn
}

// buildTypeFullName ...
func buildTypeFullName(wt walkerType) string {
	var tn string
	if wt.IsAlwaysPointer {
		tn = "*"
	}

	return tn + "ast." + wt.TypeName
}

// buildTypeKindFuncName takes the name of a kind constant, and turns it into the expected format
// for a walk function's name for that type.
func buildTypeKindFuncName(constName string) string {
	f := strings.Split(constName, "Kind")
	return f[1] + f[0]
}

// isTypeAlwaysPointer ...
// TODO(seeruk): Verify that this behaves as expected...
func isTypeAlwaysPointer(st goast.SymbolTable, tn string) bool {
	var referenced bool

	for _, str := range st.Structs {
		for _, fld := range str.Fields {
			if fld.TypeName != tn {
				continue
			}

			referenced = true

			if !fld.IsPointer {
				return false
			}
		}
	}

	return referenced
}

// isTypeLinkedList ...
func isTypeLinkedList(tn string, str goast.Struct) bool {
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

// abridger is a strings.Map function that is used to return a variable name from a type name that
// makes sense by taking each capital letter from the given string and converting them to lowercase.
// The input string should start with a capital letter.
func abridger(r rune) rune {
	if unicode.IsUpper(r) {
		return unicode.ToLower(r)
	}
	return -1
}
