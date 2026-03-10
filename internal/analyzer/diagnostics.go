package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type diagnosticKey struct {
	ruleID string
	pos    int
}

func buildDiagnostics(expr ast.Expr, violations []violation) []analysis.Diagnostic {
	if expr == nil || len(violations) == 0 {
		return nil
	}

	diagnostics := make([]analysis.Diagnostic, 0, len(violations))
	seen := make(map[diagnosticKey]struct{}, len(violations))

	for _, currentViolation := range violations {
		key := diagnosticKey{
			ruleID: currentViolation.ruleID,
			pos:    int(expr.Pos()),
		}

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Pos:      expr.Pos(),
			End:      expr.End(),
			Category: currentViolation.ruleID,
			Message:  currentViolation.message,
		})
	}

	return diagnostics
}

func reportDiagnostics(pass *analysis.Pass, diagnostics []analysis.Diagnostic) {
	for _, diagnostic := range diagnostics {
		pass.Report(diagnostic)
	}
}