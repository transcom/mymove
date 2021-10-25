package appcontextlinter

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"go/ast"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/ast/inspector"
)

var AppContextAnalyzer = &analysis.Analyzer{
	Name:     "appcontext-lint",
	Doc:      "Make sure appContext is properly used throughout codebase",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.File)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		file := node.(*ast.File)
		fmt.Print("🐲🐲🐲🐲🐲🐲🐲")
		spew.Dump(file.Decls)
	})

	// NEXT Steps: Find out how we import pop.Connection?, What exactly in file.Decls do we want to look at to find the connection we're looking for, look at AST package to see what tools are available to look for different types in a file.

	return nil, fmt.Errorf("BAHHHHHHHHH")
}
