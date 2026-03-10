package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

func extractMessage(expr ast.Expr) (messageSample, bool) {
	switch currentExpr := expr.(type) {
	case *ast.BasicLit:
		return extractBasicLiteral(currentExpr)
	case *ast.BinaryExpr:
		if currentExpr.Op != token.ADD {
			return messageSample{}, false
		}

		left, ok := extractMessage(currentExpr.X)
		if !ok {
			return messageSample{}, false
		}

		right, ok := extractMessage(currentExpr.Y)
		if !ok {
			return messageSample{}, false
		}

		return messageSample{
			text:  left.text + right.text,
			parts: append(append([]string{}, left.parts...), right.parts...),
		}, true
	default:
		return messageSample{}, false
	}
}

func extractBasicLiteral(literal *ast.BasicLit) (messageSample, bool) {
	if literal.Kind != token.STRING {
		return messageSample{}, false
	}

	text, err := strconv.Unquote(literal.Value)
	if err != nil {
		return messageSample{}, false
	}

	return messageSample{
		text:  text,
		parts: []string{text},
	}, true
}