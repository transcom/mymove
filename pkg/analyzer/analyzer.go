package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "goprintffuncname",
	Doc:  "Checks that printf-like functions are named with `f` at the end.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		params := funcDecl.Type.Params.List
		if len(params) != 2 { // [0] must be format (string), [1] must be args (...interface{})
			return true
		}

		firstParamType, ok := params[0].Type.(*ast.Ident)
		if !ok { // first param type isn't identificator so it can't be of type "string"
			return true
		}

		if firstParamType.Name != "string" { // first param (format) type is not string
			return true
		}

		secondParamType, ok := params[1].Type.(*ast.Ellipsis)
		if !ok { // args are not ellipsis (...args)
			return true
		}

		elementType, ok := secondParamType.Elt.(*ast.InterfaceType)
		if !ok { // args are not of interface type, but we need interface{}
			return true
		}

		if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
			return true // has >= 1 method in interface, but we need an empty interface "interface{}"
		}

		if strings.HasSuffix(funcDecl.Name.Name, "f") {
			return true
		}

		pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
			funcDecl.Name.Name, funcDecl.Name.Name)
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
