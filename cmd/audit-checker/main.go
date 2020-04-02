package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func initFlags(flag *pflag.FlagSet) {

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Handler represents information about the handler being scanned
type Handler struct {
	Name         string
	HandleMethod *ast.FuncDecl
}

func loadHandlersFromFile(path string) ([]*Handler, error) {
	var handlers []*Handler
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, decl := range parsedFile.Decls {
		if g, ok := decl.(*ast.GenDecl); ok {
			if g.Tok == token.TYPE {
				for _, spec := range g.Specs {
					typeSpec := spec.(*ast.TypeSpec)

					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						h := Handler{
							Name: typeSpec.Name.Name,
						}

						handlers = append(handlers, &h)
					}
				}
			}
		}

		if g, ok := decl.(*ast.FuncDecl); ok {
			if g.Name.Name == "Handle" {
				receivingHandler, ok := g.Recv.List[0].Type.(*ast.Ident)

				if !ok {
					continue
				}

				for _, h := range handlers {
					if h.Name == receivingHandler.Name {
						h.HandleMethod = g
					}
				}
			}
		}
	}

	return handlers, nil
}

func checkHandler(h *Handler) error {
	if h.HandleMethod != nil {
		for _, stmt := range h.HandleMethod.Body.List {
			if assignment, ok := stmt.(*ast.AssignStmt); ok {
				for _, expr := range assignment.Rhs {
					if callExpression, ok := expr.(*ast.CallExpr); ok {
						if fun, ok := callExpression.Fun.(*ast.SelectorExpr); ok {
							if pkg, ok := fun.X.(*ast.SelectorExpr); ok {
								if pkg.Sel.Name != "auditor" || fun.Sel.Name != "Record" {
									return errors.New("error")
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func checkFiles(directories []string, logger *zap.Logger) (fail bool, failedHandlerNames []string) {
	handlerDirectory := "./pkg/handlers"
	for _, directory := range directories {
		dir := handlerDirectory + fmt.Sprintf("/%s/", directory)

		files, dirErr := ioutil.ReadDir(dir)
		if dirErr != nil {
			logger.Fatal("reading directory", zap.Error(dirErr))
		}
		for _, file := range files {
			if file.IsDir() || strings.Contains(file.Name(), "test") {
				continue
			}
			if handlers, loadErr := loadHandlersFromFile(dir + file.Name()); loadErr == nil {
				for _, handler := range handlers {
					err := checkHandler(handler)
					if err != nil {
						failedHandlerNames = append(failedHandlerNames, handler.Name)
					}
				}
			} else {
				logger.Fatal("loading models", zap.Error(loadErr))
			}
		}
	}

	if len(failedHandlerNames) > 0 {
		return true, failedHandlerNames
	}

	return false, failedHandlerNames
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logEnv := v.GetString(cli.LoggingEnvFlag)

	logger, err := logging.Config(logEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Track if we have found at least one issue
	fail, failedHandlerNames := checkFiles([]string{"ghcapi", "primeapi"}, logger)

	if fail {
		fmt.Println("Failed handlers:")
		for _, handlerName := range failedHandlerNames {
			fmt.Println(handlerName)
		}
		// There was at least one mismatch, so exit non-zero
		os.Exit(1)
	}
}
