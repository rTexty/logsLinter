package plugin

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPluginSettingsAffectAnalyzer(t *testing.T) {
	testdata, err := filepath.Abs(filepath.Join("..", "testdata"))
	if err != nil {
		t.Fatalf("filepath.Abs(testdata): %v", err)
	}

	plugin, err := New(map[string]any{
		"rules": map[string]any{
			"lowercase-start": false,
		},
		"sensitive-data": map[string]any{
			"additional-keywords": []string{"credential"},
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("BuildAnalyzers() count = %d, want 1", len(analyzers))
	}

	analysistest.Run(t, testdata, analyzers[0], "configcases")
}
