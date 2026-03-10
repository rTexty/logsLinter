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

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"token",
	"secret",
	"api_key",
	"apikey",
	"auth",
}

type messageSample struct {
	text  string
	parts []string
}

type violation struct {
	ruleID  string
	message string
}

func evaluateRules(sample messageSample, config Config) []violation {
	violations := make([]violation, 0, 4)

	if config.Rules.LowercaseStart {
		if violation, ok := checkLowercaseStart(sample.text); ok {
			violations = append(violations, violation)
		}
	}

	if config.Rules.ASCIIOnly {
		if violation, ok := checkASCIIOnly(sample.text); ok {
			violations = append(violations, violation)
		}
	}

	if config.Rules.NoSpecialChars {
		if violation, ok := checkNoSpecialCharsOrEmoji(sample.text); ok {
			violations = append(violations, violation)
		}
	}

	if config.Rules.SensitiveData {
		if violation, ok := checkSensitiveDataWithPatterns(sample, config.sensitivePatterns()); ok {
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
	return checkSensitiveDataWithPatterns(sample, defaultConfig.sensitivePatterns())
}

func checkSensitiveDataWithPatterns(sample messageSample, patterns []*regexp.Regexp) (violation, bool) {
	parts := sample.parts
	if len(parts) == 0 {
		parts = []string{sample.text}
	}

	for _, part := range parts {
		normalized := strings.ToLower(part)

		for _, pattern := range patterns {
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
