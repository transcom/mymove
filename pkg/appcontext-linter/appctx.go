// Package appcontextlinter This linter makes sure that we only use *pop.Connections in a few allowed places. We want
// to do this so that we can ensure our database connections are set up properly and remain consistent for each request
// we handle. See checkIfPackageCanBeSkipped below for list of packages allowed to use *pop.Connection directly.
// Another allowed use is handlers.handlerConfig because it sets up appcontext with a DB connection for our handlers.
package appcontextlinter

import (
	"go/ast"
	"go/token"

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
	inspectorInstance := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.File)(nil),
	}

	inspectorInstance.Preorder(nodeFilter, func(node ast.Node) {
		file := node.(*ast.File)
		packageName := file.Name.Name

		if checkIfPackageCanBeSkipped(packageName) {
			return
		}

		for _, declaration := range file.Decls {
			switch decl := declaration.(type) {
			case *ast.FuncDecl:
				paramsIncludePopConnection := checkIfFuncParamsIncludePopConnection(decl, packageName)

				if paramsIncludePopConnection {
					pass.Reportf(decl.Pos(), "Please use appcontext instead of pop.Connection.")
				}

				continue

			case *ast.GenDecl:
				positionToFlag := checkForPopConnectionUsesInDeclaration(decl, packageName)

				if positionToFlag.IsValid() {
					pass.Reportf(positionToFlag, "Please remove pop.Connection from the struct if not in allowed places. See pkg/appcontext-linter/appctx.go for valid placements.")
				}
			default:
				continue
			}
		}
	})
	return nil, nil
}

func checkIfPackageCanBeSkipped(packageName string) bool {
	// These are the packages that are allowed to have pop.Connection in them. This is strictly
	// at the package level and does not include subpackages.
	allowedPackages := map[string]bool{
		"appcontext":   true,
		"db":           true,
		"migrate":      true,
		"models":       true,
		"roles":        true,
		"testdatagen":  true,
		"testingsuite": true,
		"utilities":    true,
	}

	return allowedPackages[packageName]
}

func checkIfFuncParamsIncludePopConnection(funcToCheck *ast.FuncDecl, packageName string) bool {
	if funcToCheck.Name.Name == "NewHandlerConfig" && packageName == "handlers" {
		return false
	}

	if funcToCheck.Name.Name == "NewHandlerConfigForTest" && packageName == "handlers" {
		return false
	}
	for _, param := range funcToCheck.Type.Params.List {
		if checkForPopConnection(param) {
			return true
		}
	}

	return false
}

func checkForPopConnectionUsesInDeclaration(declarationToCheck *ast.GenDecl, packageName string) token.Pos {
	for _, spec := range declarationToCheck.Specs {
		// Only want types, not imports, variables, or constants
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		// Special case because this is the one that sets up all handlers so it's allowed to use *pop.Connection
		if typeSpec.Name.Name == "Config" && packageName == "handlers" {
			continue
		}

		// Specifically care about struct types
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// Checking the fields of the struct
		for _, structField := range structType.Fields.List {
			if checkForPopConnection(structField) {
				return typeSpec.Pos()
			}
		}
	}

	return token.NoPos
}

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
