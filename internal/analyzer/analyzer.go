package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer reports log message policy violations in supported logger calls.
var Analyzer = &analysis.Analyzer{
	Name: "logslinter",
	Doc:  "report invalid log messages in supported logger calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			analyzeCall(pass, call)
			return true
		})
	}

	return nil, nil
}

func analyzeCall(pass *analysis.Pass, call *ast.CallExpr) {
	inspectedCall, ok := inspectLogCall(pass, call)
	if !ok {
		return
	}

	sample, ok := extractMessage(inspectedCall.message)
	if !ok {
		return
	}

	violations := evaluateRules(sample)
	if len(violations) == 0 {
		return
	}

	reportDiagnostics(pass, buildDiagnostics(inspectedCall.message, violations))
}
