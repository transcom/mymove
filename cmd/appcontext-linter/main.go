package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	appcontextlinter "github.com/transcom/mymove/pkg/appcontext-linter"
)

func main() {
	singlechecker.Main(appcontextlinter.AppContextAnalyzer)
}
