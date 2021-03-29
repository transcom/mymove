package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/transcom/mymove/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
