package analyzer

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type analyzerRunner struct {
	config             Config
	additionalKeywords keywordListFlag
}

// Analyzer reports log message policy violations in supported logger calls.
var Analyzer = NewAnalyzer(defaultConfig)

func NewAnalyzer(config Config) *analysis.Analyzer {
	runner := &analyzerRunner{
		config: config.normalized(),
	}

	analyzer := &analysis.Analyzer{
		Name: "logslinter",
		Doc:  "report invalid log messages in supported logger calls",
		Run:  runner.run,
		Flags: *flag.NewFlagSet(
			"logslinter",
			flag.ExitOnError,
		),
	}

	bindConfigFlags(&analyzer.Flags, &runner.config, &runner.additionalKeywords)
	return analyzer
}

func (runner *analyzerRunner) run(pass *analysis.Pass) (any, error) {
	currentConfig := runner.config.normalized()
	currentConfig.SensitiveData.AdditionalKeywords = append([]string{}, runner.additionalKeywords.values...)
	currentConfig = currentConfig.normalized()

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			analyzeCall(pass, call, currentConfig)
			return true
		})
	}

	return nil, nil
}

func analyzeCall(pass *analysis.Pass, call *ast.CallExpr, config Config) {
	inspectedCall, ok := inspectLogCall(pass, call)
	if !ok {
		return
	}

	sample, ok := extractMessage(inspectedCall.message)
	if !ok {
		return
	}

	violations := evaluateRules(sample, config)
	if len(violations) == 0 {
		return
	}

	reportDiagnostics(pass, buildDiagnostics(inspectedCall.message, violations))
}
