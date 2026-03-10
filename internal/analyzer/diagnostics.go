package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode/utf8"

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
			Pos:            expr.Pos(),
			End:            expr.End(),
			Category:       currentViolation.ruleID,
			Message:        currentViolation.message,
			SuggestedFixes: buildSuggestedFixes(expr, currentViolation),
		})
	}

	return diagnostics
}

func reportDiagnostics(pass *analysis.Pass, diagnostics []analysis.Diagnostic) {
	for _, diagnostic := range diagnostics {
		pass.Report(diagnostic)
	}
}

func buildSuggestedFixes(expr ast.Expr, currentViolation violation) []analysis.SuggestedFix {
	if currentViolation.ruleID != ruleLowercaseStart {
		return nil
	}

	fixedLiteral, ok := buildLowercaseLiteralReplacement(expr)
	if !ok {
		return nil
	}

	return []analysis.SuggestedFix{{
		Message: "Lowercase the first letter",
		TextEdits: []analysis.TextEdit{{
			Pos:     expr.Pos(),
			End:     expr.End(),
			NewText: []byte(fixedLiteral),
		}},
	}}
}

func buildLowercaseLiteralReplacement(expr ast.Expr) (string, bool) {
	literal, ok := expr.(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING || len(literal.Value) == 0 || literal.Value[0] != '"' {
		return "", false
	}

	text, err := strconv.Unquote(literal.Value)
	if err != nil || text == "" {
		return "", false
	}

	firstRune, runeWidth := utf8.DecodeRuneInString(text)
	if firstRune < 'A' || firstRune > 'Z' {
		return "", false
	}

	lowercased := string(firstRune+('a'-'A')) + text[runeWidth:]
	if lowercased == text {
		return "", false
	}

	return strconv.Quote(lowercased), true
}
