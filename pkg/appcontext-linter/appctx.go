package appcontextlinter

import "fmt"

var Analyzer = &analysis.Analyzer{
	Name: "",
	Doc:  "",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	return nil, fmt.Errorf("this is an error")
}
