package analyzer

import (
	"go/parser"
	"testing"
)

func TestBuildDiagnostics(t *testing.T) {
	t.Parallel()

	expr, err := parser.ParseExpr(`"Starting server"`)
	if err != nil {
		t.Fatalf("ParseExpr: %v", err)
	}

	violations := []violation{
		{ruleID: ruleLowercaseStart, message: msgLowercaseStart},
		{ruleID: ruleLowercaseStart, message: msgLowercaseStart},
		{ruleID: ruleASCIIOnly, message: msgASCIIOnly},
	}

	diagnostics := buildDiagnostics(expr, violations)
	if len(diagnostics) != 2 {
		t.Fatalf("buildDiagnostics() count = %d, want 2", len(diagnostics))
	}

	if diagnostics[0].Category != ruleLowercaseStart {
		t.Fatalf("first diagnostic category = %q, want %q", diagnostics[0].Category, ruleLowercaseStart)
	}

	if diagnostics[0].Message != msgLowercaseStart {
		t.Fatalf("first diagnostic message = %q, want %q", diagnostics[0].Message, msgLowercaseStart)
	}

	if diagnostics[0].Pos != expr.Pos() || diagnostics[0].End != expr.End() {
		t.Fatalf("first diagnostic range = [%d,%d], want [%d,%d]", diagnostics[0].Pos, diagnostics[0].End, expr.Pos(), expr.End())
	}

	if diagnostics[1].Category != ruleASCIIOnly {
		t.Fatalf("second diagnostic category = %q, want %q", diagnostics[1].Category, ruleASCIIOnly)
	}

	if len(diagnostics[0].SuggestedFixes) != 1 {
		t.Fatalf("first diagnostic SuggestedFixes count = %d, want 1", len(diagnostics[0].SuggestedFixes))
	}

	firstEdit := diagnostics[0].SuggestedFixes[0].TextEdits[0]
	if string(firstEdit.NewText) != `"starting server"` {
		t.Fatalf("first diagnostic fix = %s, want %s", string(firstEdit.NewText), `"starting server"`)
	}
}

func TestBuildDiagnosticsNilCases(t *testing.T) {
	t.Parallel()

	if diagnostics := buildDiagnostics(nil, []violation{{ruleID: ruleASCIIOnly, message: msgASCIIOnly}}); diagnostics != nil {
		t.Fatalf("buildDiagnostics(nil, violations) = %#v, want nil", diagnostics)
	}

	expr, err := parser.ParseExpr(`"ok"`)
	if err != nil {
		t.Fatalf("ParseExpr: %v", err)
	}

	if diagnostics := buildDiagnostics(expr, nil); diagnostics != nil {
		t.Fatalf("buildDiagnostics(expr, nil) = %#v, want nil", diagnostics)
	}
}

func TestBuildDiagnosticsSkipsUnsafeSuggestedFixes(t *testing.T) {
	t.Parallel()

	expr, err := parser.ParseExpr("`Starting server`")
	if err != nil {
		t.Fatalf("ParseExpr: %v", err)
	}

	diagnostics := buildDiagnostics(expr, []violation{{ruleID: ruleLowercaseStart, message: msgLowercaseStart}})
	if len(diagnostics) != 1 {
		t.Fatalf("buildDiagnostics() count = %d, want 1", len(diagnostics))
	}

	if len(diagnostics[0].SuggestedFixes) != 0 {
		t.Fatalf("raw string SuggestedFixes count = %d, want 0", len(diagnostics[0].SuggestedFixes))
	}
}
