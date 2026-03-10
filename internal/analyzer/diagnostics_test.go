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
