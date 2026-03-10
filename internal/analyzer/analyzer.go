package analyzer

import "golang.org/x/tools/go/analysis"

// Analyzer reports log message policy violations in supported logger calls.
var Analyzer = &analysis.Analyzer{
	Name: "logslinter",
	Doc:  "report invalid log messages in supported logger calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	return nil, nil
}