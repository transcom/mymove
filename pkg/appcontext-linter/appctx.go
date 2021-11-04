package appcontextlinter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var AppContextAnalyzer = &analysis.Analyzer{
	Name:     "appcontextlint",
	Doc:      "Make sure appContext is properly used throughout codebase",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	// Analyze code and make an AST from the file:
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.File)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		file := node.(*ast.File)

		allowList := map[string]bool{
			"appcontext":   true,
			"db":           true,
			"migrate":      true,
			"models":       true,
			"testdatagen":  true,
			"testingsuite": true,
		}

		if allowList[file.Name.Name] {
			return
		}

		for _, declaration := range file.Decls {
			t, ok := declaration.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range t.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if typeSpec.Name.Name == "handlerContext" && file.Name.Name == "handlers" {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				// Checking the fields of the structs
				for _, structField := range structType.Fields.List {
					if checkForPopConnection(structField) {
						pass.Reportf(typeSpec.Pos(), "Please remove pop.Connection from the struct if not in appcontext")
						continue
					}
				}

			}

		}
	})
	return nil, nil
}

// TODO: Add logic to get it to run in circleCI and when run locally

func checkForPopConnection(field *ast.Field) bool {
	// Look for a type called StarExpr where pop Connection might be
	if identifier, ok := field.Type.(*ast.StarExpr); ok {
		// Look for a Struct that may contain "pop" and "connection"
		if findPop, ok := identifier.X.(*ast.SelectorExpr); ok {
			foundPop := false
			// Once inside the struct, look for "pop"
			if popIdentifier, ok := findPop.X.(*ast.Ident); ok {
				if popIdentifier.Name == "pop" {
					foundPop = true
				}
			}
			// If "pop" not found, move on
			if !foundPop {
				return false
			}
			// After pop is found, look for "connection" and report if it's found
			if findPop.Sel.Name == "Connection" {
				return true
			}
		}
	}
	return false
}
