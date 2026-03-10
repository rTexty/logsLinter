package analyzer

import (
	"go/ast"
	"go/parser"
	"reflect"
	"testing"
)

func TestExtractMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		expr      string
		wantOK    bool
		wantText  string
		wantParts []string
	}{
		{
			name:      "string literal",
			expr:      `"starting server"`,
			wantOK:    true,
			wantText:  "starting server",
			wantParts: []string{"starting server"},
		},
		{
			name:      "raw string literal",
			expr:      "`starting server`",
			wantOK:    true,
			wantText:  "starting server",
			wantParts: []string{"starting server"},
		},
		{
			name:      "nested concatenation",
			expr:      `"start" + "ing" + " server"`,
			wantOK:    true,
			wantText:  "starting server",
			wantParts: []string{"start", "ing", " server"},
		},
		{
			name:   "variable is skipped",
			expr:   `message`,
			wantOK: false,
		},
		{
			name:   "function call is skipped",
			expr:   `buildMessage()`,
			wantOK: false,
		},
		{
			name:   "formatted expression is skipped",
			expr:   `fmt.Sprintf("hello %s", name)`,
			wantOK: false,
		},
		{
			name:   "literal plus variable is skipped",
			expr:   `"password: " + password`,
			wantOK: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			expr := mustParseExpr(t, testCase.expr)
			got, ok := extractMessage(expr)
			if ok != testCase.wantOK {
				t.Fatalf("extractMessage(%q) ok = %v, want %v", testCase.expr, ok, testCase.wantOK)
			}

			if !testCase.wantOK {
				return
			}

			if got.text != testCase.wantText {
				t.Fatalf("extractMessage(%q) text = %q, want %q", testCase.expr, got.text, testCase.wantText)
			}

			if !reflect.DeepEqual(got.parts, testCase.wantParts) {
				t.Fatalf("extractMessage(%q) parts = %#v, want %#v", testCase.expr, got.parts, testCase.wantParts)
			}
		})
	}
}

func mustParseExpr(t *testing.T, expr string) ast.Expr {
	t.Helper()

	parsedExpr, err := parser.ParseExpr(expr)
	if err != nil {
		t.Fatalf("ParseExpr(%q): %v", expr, err)
	}

	return parsedExpr
}