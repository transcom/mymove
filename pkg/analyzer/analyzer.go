package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer describes an analysis function and its options.
var Analyzer = &analysis.Analyzer{
	Name: "goprintffuncname",
	Doc:  "Checks that printf-like functions are named with `f` at the end.",
	Run:  run,
}

// render returns the pretty-print of the given node
/*func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
*/

func containsGosecDisableNoReason(comments []*ast.Comment) bool {
	const disableNoSec = "#nosec"
	for _, comment := range comments {
		if strings.Contains(comment.Text, disableNoSec) {
			individualCommentArr := strings.Split(comment.Text, " ")
			for index, str := range individualCommentArr {
				if str == disableNoSec && index == len(individualCommentArr)-1 {
					return true
				}
			}
		}
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		comments, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		containsDisablingGosecWithNoReason := containsGosecDisableNoReason(comments.List)

		if containsDisablingGosecWithNoReason {
			pass.Reportf(node.Pos(), "Please provide gosec rule that is being disabled")
			return true
		}
		/*fmt.Println(funcDecl)
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
			funcDecl.Name.Name, funcDecl.Name.Name)*/
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
