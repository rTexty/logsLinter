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

func TestInspectLogCallRecognizesZapMessages(t *testing.T) {
	t.Parallel()

	pass, calls := mustTypeCheckCallsWithImporter(t, `package fixture

import "go.uber.org/zap"

var (
	logger = zap.NewNop()
	sugar = logger.Sugar()
)

func example() {
	logger.Info("logger info")
	logger.Warn("logger warn", zap.String("component", "api"))
	logger.With(zap.String("component", "db")).Error("logger with")
	logger.Sugar().Infow("sugar chain", "component", "api")
	sugar.Debugw("sugar debug", "component", "worker")
	sugar.Errorw("sugar error", "component", "db")
	sugar.Info("skip print style")
	sugar.Infof("skip %s", "format")
}
`, newFixtureImporter(t))

	testCases := []struct {
		callText string
		wantMsg  string
		wantOK   bool
	}{
		{callText: `logger.Info("logger info")`, wantMsg: `"logger info"`, wantOK: true},
		{callText: `logger.Warn("logger warn", zap.String("component", "api"))`, wantMsg: `"logger warn"`, wantOK: true},
		{callText: `logger.With(zap.String("component", "db")).Error("logger with")`, wantMsg: `"logger with"`, wantOK: true},
		{callText: `logger.Sugar().Infow("sugar chain", "component", "api")`, wantMsg: `"sugar chain"`, wantOK: true},
		{callText: `sugar.Debugw("sugar debug", "component", "worker")`, wantMsg: `"sugar debug"`, wantOK: true},
		{callText: `sugar.Errorw("sugar error", "component", "db")`, wantMsg: `"sugar error"`, wantOK: true},
		{callText: `sugar.Info("skip print style")`, wantOK: false},
		{callText: `sugar.Infof("skip %s", "format")`, wantOK: false},
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

			if inspectedCall.family != loggerFamilyZap {
				t.Fatalf("inspectLogCall(%s) family = %q, want %q", testCase.callText, inspectedCall.family, loggerFamilyZap)
			}

			if got := formatExpr(t, pass.Fset, inspectedCall.message); got != testCase.wantMsg {
				t.Fatalf("inspectLogCall(%s) message = %s, want %s", testCase.callText, got, testCase.wantMsg)
			}
		})
	}
}

func mustTypeCheckCalls(t *testing.T, source string) (*analysis.Pass, map[string]*ast.CallExpr) {
	t.Helper()

	return mustTypeCheckCallsWithImporter(t, source, importer.Default())

}

func mustTypeCheckCallsWithImporter(t *testing.T, source string, currentImporter types.Importer) (*analysis.Pass, map[string]*ast.CallExpr) {
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

	config := types.Config{Importer: currentImporter}
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

func newFixtureImporter(t *testing.T) types.Importer {
	t.Helper()

	defaultImporter := importer.Default()
	packages := map[string]*types.Package{
		zapPackagePath: mustTypeCheckPackage(t, zapPackagePath, `package zap

type Field struct{}

type Logger struct{}

type SugaredLogger struct{}

func NewNop() *Logger { return &Logger{} }

func String(string, string) Field { return Field{} }

func (*Logger) Debug(string, ...Field) {}

func (*Logger) Info(string, ...Field) {}

func (*Logger) Warn(string, ...Field) {}

func (*Logger) Error(string, ...Field) {}

func (logger *Logger) With(...Field) *Logger { return logger }

func (*Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

func (*SugaredLogger) Debugw(string, ...any) {}

func (*SugaredLogger) Infow(string, ...any) {}

func (*SugaredLogger) Warnw(string, ...any) {}

func (*SugaredLogger) Errorw(string, ...any) {}

func (*SugaredLogger) Debug(...any) {}

func (*SugaredLogger) Info(...any) {}

func (*SugaredLogger) Warn(...any) {}

func (*SugaredLogger) Error(...any) {}

func (*SugaredLogger) Debugf(string, ...any) {}

func (*SugaredLogger) Infof(string, ...any) {}

func (*SugaredLogger) Warnf(string, ...any) {}

func (*SugaredLogger) Errorf(string, ...any) {}
`),
	}

	return fixtureImporter{
		fallback: defaultImporter,
		packages: packages,
	}
}

func mustTypeCheckPackage(t *testing.T, path, source string) *types.Package {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path+".go", source, 0)
	if err != nil {
		t.Fatalf("ParseFile(%s): %v", path, err)
	}

	config := types.Config{}
	pkg, err := config.Check(path, fset, []*ast.File{file}, nil)
	if err != nil {
		t.Fatalf("Check(%s): %v", path, err)
	}

	return pkg
}

type fixtureImporter struct {
	fallback types.Importer
	packages map[string]*types.Package
}

func (importer fixtureImporter) Import(path string) (*types.Package, error) {
	if pkg, ok := importer.packages[path]; ok {
		return pkg, nil
	}

	return importer.fallback.Import(path)
}

func formatExpr(t *testing.T, fset *token.FileSet, expr ast.Node) string {
	t.Helper()

	var buffer bytes.Buffer
	if err := format.Node(&buffer, fset, expr); err != nil {
		t.Fatalf("format.Node: %v", err)
	}

	return buffer.String()
}
