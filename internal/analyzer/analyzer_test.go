package analyzer

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata, err := filepath.Abs(filepath.Join("..", "..", "testdata"))
	if err != nil {
		t.Fatalf("filepath.Abs(testdata): %v", err)
	}

	for _, packageName := range []string{"slogcases", "zapcases", "mixedcases"} {
		packageName := packageName

		t.Run(packageName, func(t *testing.T) {
			analysistest.Run(t, testdata, Analyzer, packageName)
		})
	}
}