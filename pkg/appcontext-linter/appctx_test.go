package appcontextlinter

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// this test starts up Test runner (line 18) and looks at tests in appctx_linter_tests and runs linter against those files
// if there are no want statements, the linter moves on to the next statement
// if there are no want statements and the test fails, our linter is failing because nothing is expected.
func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, AppContextAnalyzer, "appctx_linter_tests")
}
