package main

import (
	appcontextlinter "github.com/transcom/mymove/pkg/appcontext-linter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(appcontextlinter.ATOAnalyzer)
}
