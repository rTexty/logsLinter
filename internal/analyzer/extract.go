package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func extractMessage(expr ast.Expr) (messageSample, bool) {
	parts, ok := extractMessageParts(expr)
	if !ok {
		return messageSample{}, false
	}

	return newMessageSample(parts), true
}

func extractMessageParts(expr ast.Expr) ([]string, bool) {
	switch currentExpr := expr.(type) {
	case *ast.BasicLit:
		sample, ok := extractBasicLiteral(currentExpr)
		if !ok {
			return nil, false
		}

		return sample.parts, true
	case *ast.BinaryExpr:
		if currentExpr.Op != token.ADD {
			return nil, false
		}

		left, ok := extractMessageParts(currentExpr.X)
		if !ok {
			return nil, false
		}

		right, ok := extractMessageParts(currentExpr.Y)
		if !ok {
			return nil, false
		}

		return append(append([]string{}, left...), right...), true
	default:
		return nil, false
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

func newMessageSample(parts []string) messageSample {
	normalizedParts := append([]string{}, parts...)

	return messageSample{
		text:  strings.Join(normalizedParts, ""),
		parts: normalizedParts,
	}
}
