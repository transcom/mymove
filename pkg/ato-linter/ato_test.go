package atolinter

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	t.Skip("skip for now, we wont be using this as official scanning for code issues")
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, ATOAnalyzer, "ato_linter_tests")
}
