package analyzer

import (
	"unicode"
	"unicode/utf8"
)

const (
	ruleLowercaseStart = "lowercase-start"
	ruleASCIIOnly      = "english-ascii-only"
)

const (
	msgLowercaseStart = "log message must start with a lowercase letter"
	msgASCIIOnly      = "log message must be in English (ASCII only)"
)

type messageSample struct {
	text string
}

type violation struct {
	ruleID  string
	message string
}

func evaluateRules(sample messageSample) []violation {
	violations := make([]violation, 0, 2)

	if violation, ok := checkLowercaseStart(sample.text); ok {
		violations = append(violations, violation)
	}

	if violation, ok := checkASCIIOnly(sample.text); ok {
		violations = append(violations, violation)
	}

	return violations
}

func checkLowercaseStart(text string) (violation, bool) {
	if text == "" {
		return violation{}, false
	}

	firstRune, _ := utf8.DecodeRuneInString(text)
	if unicode.IsUpper(firstRune) {
		return violation{
			ruleID:  ruleLowercaseStart,
			message: msgLowercaseStart,
		}, true
	}

	return violation{}, false
}

func checkASCIIOnly(text string) (violation, bool) {
	for _, currentRune := range text {
		if currentRune < 0x20 || currentRune > 0x7e {
			return violation{
				ruleID:  ruleASCIIOnly,
				message: msgASCIIOnly,
			}, true
		}
	}

	return violation{}, false
}