package analyzer

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestInspectLogCallRecognizesSlogMessages(t *testing.T) {
	t.Parallel()

	pass, calls := mustTypeCheckCalls(t, `package fixture

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

var (
	ctx = context.Background()
	logger = slog.New(slog.NewTextHandler(io.Discard, nil))
)

func example() {
	slog.Info("top level")
	slog.InfoContext(ctx, "context call")
	slog.Log(ctx, slog.LevelInfo, "log call")
	slog.LogAttrs(ctx, slog.LevelInfo, "attrs call")
	logger.Info("method call")
	logger.WarnContext(ctx, "warn context")
	logger.Log(ctx, slog.LevelWarn, "method log")
	logger.LogAttrs(ctx, slog.LevelError, "method attrs")
	logger.With("component", "api").Info("with chain")
	logger.WithGroup("db").ErrorContext(ctx, "group chain")
	fmt.Println("skip")
	slog.SetDefault(logger)
}
`)

	testCases := []struct {
		callText string
		wantMsg  string
		wantOK   bool
	}{
		{callText: `slog.Info("top level")`, wantMsg: `"top level"`, wantOK: true},
		{callText: `slog.InfoContext(ctx, "context call")`, wantMsg: `"context call"`, wantOK: true},
		{callText: `slog.Log(ctx, slog.LevelInfo, "log call")`, wantMsg: `"log call"`, wantOK: true},
		{callText: `slog.LogAttrs(ctx, slog.LevelInfo, "attrs call")`, wantMsg: `"attrs call"`, wantOK: true},
		{callText: `logger.Info("method call")`, wantMsg: `"method call"`, wantOK: true},
		{callText: `logger.WarnContext(ctx, "warn context")`, wantMsg: `"warn context"`, wantOK: true},
		{callText: `logger.Log(ctx, slog.LevelWarn, "method log")`, wantMsg: `"method log"`, wantOK: true},
		{callText: `logger.LogAttrs(ctx, slog.LevelError, "method attrs")`, wantMsg: `"method attrs"`, wantOK: true},
		{callText: `logger.With("component", "api").Info("with chain")`, wantMsg: `"with chain"`, wantOK: true},
		{callText: `logger.WithGroup("db").ErrorContext(ctx, "group chain")`, wantMsg: `"group chain"`, wantOK: true},
		{callText: `fmt.Println("skip")`, wantOK: false},
		{callText: `slog.SetDefault(logger)`, wantOK: false},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.callText, func(t *testing.T) {
			t.Parallel()

			call := calls[testCase.callText]
			if call == nil {
				t.Fatalf("call %q not found", testCase.callText)
			}

			inspectedCall, ok := inspectLogCall(pass, call)
			if ok != testCase.wantOK {
				t.Fatalf("inspectLogCall(%s) ok = %v, want %v", testCase.callText, ok, testCase.wantOK)
			}

			if !testCase.wantOK {
				return
			}

			if inspectedCall.family != loggerFamilySlog {
				t.Fatalf("inspectLogCall(%s) family = %q, want %q", testCase.callText, inspectedCall.family, loggerFamilySlog)
			}

			if got := formatExpr(t, pass.Fset, inspectedCall.message); got != testCase.wantMsg {
				t.Fatalf("inspectLogCall(%s) message = %s, want %s", testCase.callText, got, testCase.wantMsg)
			}
		})
	}
}

func mustTypeCheckCalls(t *testing.T, source string) (*analysis.Pass, map[string]*ast.CallExpr) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "fixture.go", source, 0)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	info := &types.Info{
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
	}

	config := types.Config{Importer: importer.Default()}
	pkg, err := config.Check("fixture", fset, []*ast.File{file}, info)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
	}

	calls := make(map[string]*ast.CallExpr)
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		calls[formatExpr(t, fset, call)] = call
		return true
	})

	return pass, calls
}

func formatExpr(t *testing.T, fset *token.FileSet, expr ast.Node) string {
	t.Helper()

	var buffer bytes.Buffer
	if err := format.Node(&buffer, fset, expr); err != nil {
		t.Fatalf("format.Node: %v", err)
	}

	return buffer.String()
}