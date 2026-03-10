package analyzer

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	ruleLowercaseStart = "lowercase-start"
	ruleASCIIOnly      = "english-ascii-only"
	ruleNoSpecialChars = "no-special-chars-or-emoji"
	ruleSensitiveData  = "no-sensitive-data"
)

const (
	msgLowercaseStart = "log message must start with a lowercase letter"
	msgASCIIOnly      = "log message must be in English (ASCII only)"
	msgNoSpecialChars = "log message must not contain special characters or emoji"
	msgSensitiveData  = "log message may contain sensitive data"
)

var sensitiveDataPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(^|[^a-z0-9_])password([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])passwd([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])token([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])secret([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])api_key([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])apikey([^a-z0-9_]|$)`),
	regexp.MustCompile(`(^|[^a-z0-9_])auth([^a-z0-9_]|$)`),
}

type messageSample struct {
	text  string
	parts []string
}

type violation struct {
	ruleID  string
	message string
}

type ruleDefinition struct {
	ruleID  string
	message string
	check   func(messageSample) (violation, bool)
}

var ruleDefinitions = []ruleDefinition{
	{
		ruleID:  ruleLowercaseStart,
		message: msgLowercaseStart,
		check: func(sample messageSample) (violation, bool) {
			return checkLowercaseStart(sample.text)
		},
	},
	{
		ruleID:  ruleASCIIOnly,
		message: msgASCIIOnly,
		check: func(sample messageSample) (violation, bool) {
			return checkASCIIOnly(sample.text)
		},
	},
	{
		ruleID:  ruleNoSpecialChars,
		message: msgNoSpecialChars,
		check: func(sample messageSample) (violation, bool) {
			return checkNoSpecialCharsOrEmoji(sample.text)
		},
	},
	{
		ruleID:  ruleSensitiveData,
		message: msgSensitiveData,
		check:   checkSensitiveData,
	},
}

func evaluateRules(sample messageSample) []violation {
	violations := make([]violation, 0, 4)

	for _, definition := range ruleDefinitions {
		violation, ok := definition.check(sample)
		if ok {
			violations = append(violations, violation)
		}
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

func checkNoSpecialCharsOrEmoji(text string) (violation, bool) {
	if strings.Contains(text, "!") || strings.Contains(text, "...") || strings.HasSuffix(text, "?") {
		return violation{
			ruleID:  ruleNoSpecialChars,
			message: msgNoSpecialChars,
		}, true
	}

	for _, currentRune := range text {
		if isEmoji(currentRune) {
			return violation{
				ruleID:  ruleNoSpecialChars,
				message: msgNoSpecialChars,
			}, true
		}
	}

	return violation{}, false
}

func checkSensitiveData(sample messageSample) (violation, bool) {
	parts := sample.parts
	if len(parts) == 0 {
		parts = []string{sample.text}
	}

	for _, part := range parts {
		normalized := strings.ToLower(part)

		for _, pattern := range sensitiveDataPatterns {
			if pattern.MatchString(normalized) {
				return violation{
					ruleID:  ruleSensitiveData,
					message: msgSensitiveData,
				}, true
			}
		}
	}

	return violation{}, false
}

func isEmoji(currentRune rune) bool {
	return currentRune >= 0x1f300 && currentRune <= 0x1faff ||
		currentRune >= 0x2600 && currentRune <= 0x27bf
}