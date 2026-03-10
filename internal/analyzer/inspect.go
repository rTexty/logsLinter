package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const slogPackagePath = "log/slog"
const zapPackagePath = "go.uber.org/zap"

type loggerFamily string

const (
	loggerFamilyUnknown loggerFamily = ""
	loggerFamilySlog    loggerFamily = "slog"
	loggerFamilyZap     loggerFamily = "zap"
)

type logCall struct {
	family  loggerFamily
	message ast.Expr
}

func inspectLogCall(pass *analysis.Pass, call *ast.CallExpr) (logCall, bool) {
	if inspectedCall, ok := inspectSlogCall(pass, call); ok {
		return inspectedCall, true
	}

	if inspectedCall, ok := inspectZapCall(pass, call); ok {
		return inspectedCall, true
	}

	return logCall{}, false
}

func inspectSlogCall(pass *analysis.Pass, call *ast.CallExpr) (logCall, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return logCall{}, false
	}

	messageIndex, ok := slogMessageArgIndex(pass, selector)
	if !ok || messageIndex >= len(call.Args) {
		return logCall{}, false
	}

	return logCall{
		family:  loggerFamilySlog,
		message: call.Args[messageIndex],
	}, true
}

func inspectZapCall(pass *analysis.Pass, call *ast.CallExpr) (logCall, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return logCall{}, false
	}

	messageIndex, ok := zapMessageArgIndex(pass, selector)
	if !ok || messageIndex >= len(call.Args) {
		return logCall{}, false
	}

	return logCall{
		family:  loggerFamilyZap,
		message: call.Args[messageIndex],
	}, true
}

func slogMessageArgIndex(pass *analysis.Pass, selector *ast.SelectorExpr) (int, bool) {
	if isSlogPackageSelector(pass, selector) {
		return slogFunctionMessageArgIndex(selector.Sel.Name)
	}

	if isSlogLoggerMethod(pass, selector) {
		return slogMethodMessageArgIndex(selector.Sel.Name)
	}

	return 0, false
}

func slogFunctionMessageArgIndex(name string) (int, bool) {
	switch name {
	case "Debug", "Info", "Warn", "Error":
		return 0, true
	case "DebugContext", "InfoContext", "WarnContext", "ErrorContext":
		return 1, true
	case "Log", "LogAttrs":
		return 2, true
	default:
		return 0, false
	}
}

func slogMethodMessageArgIndex(name string) (int, bool) {
	switch name {
	case "Debug", "Info", "Warn", "Error":
		return 0, true
	case "DebugContext", "InfoContext", "WarnContext", "ErrorContext":
		return 1, true
	case "Log", "LogAttrs":
		return 2, true
	default:
		return 0, false
	}
}

func zapMessageArgIndex(pass *analysis.Pass, selector *ast.SelectorExpr) (int, bool) {
	if isZapLoggerMethod(pass, selector) {
		return zapLoggerMessageArgIndex(selector.Sel.Name)
	}

	if isZapSugaredLoggerMethod(pass, selector) {
		return zapSugaredLoggerMessageArgIndex(selector.Sel.Name)
	}

	return 0, false
}

func zapLoggerMessageArgIndex(name string) (int, bool) {
	switch name {
	case "Debug", "Info", "Warn", "Error":
		return 0, true
	default:
		return 0, false
	}
}

func zapSugaredLoggerMessageArgIndex(name string) (int, bool) {
	switch name {
	case "Debugw", "Infow", "Warnw", "Errorw":
		return 0, true
	default:
		return 0, false
	}
}

func isSlogPackageSelector(pass *analysis.Pass, selector *ast.SelectorExpr) bool {
	pkgIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
	if !ok {
		return false
	}

	imported := pkgName.Imported()
	return imported != nil && imported.Path() == slogPackagePath
}

func isSlogLoggerMethod(pass *analysis.Pass, selector *ast.SelectorExpr) bool {
	selection := pass.TypesInfo.Selections[selector]
	if selection == nil {
		return false
	}

	return isNamedType(selection.Recv(), slogPackagePath, "Logger")
}

func isZapLoggerMethod(pass *analysis.Pass, selector *ast.SelectorExpr) bool {
	selection := pass.TypesInfo.Selections[selector]
	if selection == nil {
		return false
	}

	return isNamedType(selection.Recv(), zapPackagePath, "Logger")
}

func isZapSugaredLoggerMethod(pass *analysis.Pass, selector *ast.SelectorExpr) bool {
	selection := pass.TypesInfo.Selections[selector]
	if selection == nil {
		return false
	}

	return isNamedType(selection.Recv(), zapPackagePath, "SugaredLogger")
}

func isNamedType(currentType types.Type, packagePath, typeName string) bool {
	for {
		pointer, ok := currentType.(*types.Pointer)
		if !ok {
			break
		}

		currentType = pointer.Elem()
	}

	named, ok := currentType.(*types.Named)
	if !ok {
		return false
	}

	object := named.Obj()
	if object == nil || object.Name() != typeName || object.Pkg() == nil {
		return false
	}

	return object.Pkg().Path() == packagePath
}