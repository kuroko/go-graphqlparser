package rules

import (
	"fmt"

	"github.com/bucketd/go-graphqlparser/ast"
	"github.com/bucketd/go-graphqlparser/graphql"
	"github.com/bucketd/go-graphqlparser/validation"
)

type operationInfo struct {
	definedVars, seenVars, frags []string
}

type fragmentInfo struct {
	seenVars, frags []string
}

func noUnusedVariables(walker *validation.Walker) {
	fragName := ""
	fragInfos := make(map[string]fragmentInfo)

	opName := ""
	opInfos := make(map[string]operationInfo)

	walker.AddFragmentDefinitionEnterEventHandler(func(context *validation.Context, fd *ast.FragmentDefinition) { fragName = fd.Name })
	walker.AddFragmentDefinitionLeaveEventHandler(func(context *validation.Context, fd *ast.FragmentDefinition) { fragName = "" })
	walker.AddOperationDefinitionEnterEventHandler(func(context *validation.Context, od *ast.OperationDefinition) { opName = od.Name })
	walker.AddOperationDefinitionLeaveEventHandler(func(context *validation.Context, od *ast.OperationDefinition) { opName = "" })

	walker.AddVariableDefinitionEnterEventHandler(func(context *validation.Context, vd ast.VariableDefinition) {
		if len(opName) > 0 {
			oi := opInfos[opName]
			oi.definedVars = append(oi.definedVars, vd.Name)
			opInfos[opName] = oi
		}
	})

	walker.AddVariableValueEnterEventHandler(func(context *validation.Context, vv ast.Value) {
		if len(fragName) > 0 {
			fi := fragInfos[fragName]
			fi.seenVars = append(fi.seenVars, vv.StringValue)
			fragInfos[fragName] = fi
		}

		if len(opName) > 0 {
			oi := opInfos[opName]
			oi.seenVars = append(oi.seenVars, vv.StringValue)
			opInfos[opName] = oi
		}
	})

	walker.AddFragmentSpreadSelectionEnterEventHandler(func(context *validation.Context, fs ast.Selection) {
		if len(opName) > 0 {
			oi := opInfos[opName]
			oi.frags = append(oi.frags, fs.Name)
			opInfos[opName] = oi
		}
	})

	walker.AddDocumentLeaveEventHandler(func(context *validation.Context, d ast.Document) {
		for opName, oi := range opInfos {

			oi.frags = whatTheFrag(oi.frags, fragInfos, []string{})

			for _, frag := range oi.frags {
				oi.seenVars = append(oi.seenVars, fragInfos[frag].seenVars...)
			}

			for _, def := range oi.definedVars {
				var used bool
				for _, seen := range oi.seenVars {
					if def == seen {
						used = true
					}
				}
				if !used {
					context.Errors = context.Errors.Add(unusedVariableMessage(def, opName, 0, 0))
				}
			}
		}
	})
}

func whatTheFrag(frags []string, fragInfos map[string]fragmentInfo, seenFrags []string) []string {
	for _, frag := range frags {
		var seen bool
		for _, seenFrag := range seenFrags {
			if frag == seenFrag {
				seen = true
			}
		}
		if !seen {
			seenFrags = append(seenFrags, whatTheFrag([]string{frag}, fragInfos, append(seenFrags, frag))...)
		}
	}
	return seenFrags
}

func unusedVariableMessage(varName, opName string, line, col int) graphql.Error {
	msg := fmt.Sprintf("Variable %s is never used", varName)

	if len(opName) > 0 {
		msg += fmt.Sprintf(" in operation %s", opName)
	}

	return graphql.NewError(
		msg + ".",
		// TODO(seeruk): Location.
	)
}
