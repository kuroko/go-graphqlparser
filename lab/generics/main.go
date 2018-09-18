package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"unicode/utf8"
)

var (
	packageName string
	typeNames   string
)

func main() {
	flag.StringVar(&packageName, "package", "", "The package name to use in the generated code.")
	flag.StringVar(&typeNames, "types", "", "Comma separated names of types to generate.")
	flag.Parse()

	fmt.Fprintf(os.Stdout, "// Package %s - THIS CODE IS GENERATED, DO NOT EDIT MANUALLY\n", packageName)
	fmt.Fprintf(os.Stdout, "package %s\n", packageName)

	tns := strings.Split(typeNames, ",")

	sort.Strings(tns)

	for _, tn := range tns {
		typeNameLCF := lcfirst(tn)
		if typeNameLCF == "type" {
			typeNameLCF = "t"
		}

		linkedList.Execute(os.Stdout, map[string]string{
			"TypeNameLCF": typeNameLCF,
			"TypeName":    tn,
		})
	}
}

var linkedList = template.Must(template.New("linkedList").Parse(`
// {{.TypeName}}s is a linked list that contains {{.TypeName}} values.
type {{.TypeName}}s struct {
	Data {{.TypeName}}
	Next *{{.TypeName}}s
}

// ForEach applies the given map function to each item in this linked list.
func (d *{{.TypeName}}s) ForEach(fn func({{.TypeNameLCF}} {{.TypeName}}, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

	for {
		fn(current.Data, iter)

		if current.Next == nil {
			break
		}

		iter++
		current = current.Next
	}
}

// Len returns the length of this linked list. 
func (d *{{.TypeName}}s) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of {{.TypeName}}. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *{{.TypeName}}s) Reverse() *{{.TypeName}}s {
	current := d

	var prev *{{.TypeName}}s
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}
`))

func lcfirst(in string) string {
	if len(in) == 0 {
		return in
	}

	if len(in) == 1 {
		return strings.ToLower(in)
	}

	fr, w := utf8.DecodeRuneInString(in)

	return strings.ToLower(string(fr)) + in[w:]
}
