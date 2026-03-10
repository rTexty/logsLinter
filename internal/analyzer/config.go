package analyzer

import (
	"flag"
	"regexp"
	"sort"
	"strings"
)

type Config struct {
	Rules         RuleConfig
	SensitiveData SensitiveDataConfig
}

type RuleConfig struct {
	LowercaseStart bool
	ASCIIOnly      bool
	NoSpecialChars bool
	SensitiveData  bool
}

type SensitiveDataConfig struct {
	AdditionalKeywords []string
}

type Settings struct {
	Rules         RuleSettings          `json:"rules"`
	SensitiveData SensitiveDataSettings `json:"sensitive-data"`
}

type RuleSettings struct {
	LowercaseStart *bool `json:"lowercase-start"`
	ASCIIOnly      *bool `json:"english-ascii-only"`
	NoSpecialChars *bool `json:"no-special-chars-or-emoji"`
	SensitiveData  *bool `json:"no-sensitive-data"`
}

type SensitiveDataSettings struct {
	AdditionalKeywords []string `json:"additional-keywords"`
}

type keywordListFlag struct {
	values []string
}

var defaultConfig = Config{
	Rules: RuleConfig{
		LowercaseStart: true,
		ASCIIOnly:      true,
		NoSpecialChars: true,
		SensitiveData:  true,
	},
	SensitiveData: SensitiveDataConfig{
		AdditionalKeywords: nil,
	},
}

func (settings Settings) ToConfig() Config {
	config := defaultConfig

	if settings.Rules.LowercaseStart != nil {
		config.Rules.LowercaseStart = *settings.Rules.LowercaseStart
	}

	if settings.Rules.ASCIIOnly != nil {
		config.Rules.ASCIIOnly = *settings.Rules.ASCIIOnly
	}

	if settings.Rules.NoSpecialChars != nil {
		config.Rules.NoSpecialChars = *settings.Rules.NoSpecialChars
	}

	if settings.Rules.SensitiveData != nil {
		config.Rules.SensitiveData = *settings.Rules.SensitiveData
	}

	config.SensitiveData.AdditionalKeywords = append(
		[]string{},
		settings.SensitiveData.AdditionalKeywords...,
	)

	return config.normalized()
}

func (config Config) normalized() Config {
	normalized := config
	normalized.SensitiveData.AdditionalKeywords = normalizeKeywords(normalized.SensitiveData.AdditionalKeywords)

	return normalized
}

func (config Config) sensitivePatterns() []*regexp.Regexp {
	patterns := make([]*regexp.Regexp, 0, len(defaultSensitiveKeywords)+len(config.SensitiveData.AdditionalKeywords))

	for _, keyword := range defaultSensitiveKeywords {
		patterns = append(patterns, regexp.MustCompile(sensitiveKeywordPattern(keyword)))
	}

	for _, keyword := range config.SensitiveData.AdditionalKeywords {
		patterns = append(patterns, regexp.MustCompile(sensitiveKeywordPattern(keyword)))
	}

	return patterns
}

func normalizeKeywords(keywords []string) []string {
	if len(keywords) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(keywords))
	normalized := make([]string, 0, len(keywords))

	for _, keyword := range keywords {
		trimmed := strings.TrimSpace(strings.ToLower(keyword))
		if trimmed == "" {
			continue
		}

		if _, ok := seen[trimmed]; ok {
			continue
		}

		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	sort.Strings(normalized)
	return normalized
}

func sensitiveKeywordPattern(keyword string) string {
	return `(^|[^a-z0-9_])` + regexp.QuoteMeta(strings.ToLower(keyword)) + `([^a-z0-9_]|$)`
}

func (list *keywordListFlag) String() string {
	return strings.Join(list.values, ",")
}

func (list *keywordListFlag) Set(value string) error {
	parts := strings.Split(value, ",")
	list.values = append(list.values, parts...)
	list.values = normalizeKeywords(list.values)

	return nil
}

func bindConfigFlags(flagSet *flag.FlagSet, config *Config, keywords *keywordListFlag) {
	flagSet.BoolVar(&config.Rules.LowercaseStart, "lowercase-start", config.Rules.LowercaseStart, "report log messages starting with an uppercase letter")
	flagSet.BoolVar(&config.Rules.ASCIIOnly, "english-ascii-only", config.Rules.ASCIIOnly, "report non-ASCII log messages")
	flagSet.BoolVar(&config.Rules.NoSpecialChars, "no-special-chars-or-emoji", config.Rules.NoSpecialChars, "report decorative punctuation and emoji in log messages")
	flagSet.BoolVar(&config.Rules.SensitiveData, "no-sensitive-data", config.Rules.SensitiveData, "report sensitive keywords in log messages")

	keywords.values = append([]string{}, config.SensitiveData.AdditionalKeywords...)
	flagSet.Var(keywords, "additional-sensitive-keywords", "comma-separated additional sensitive keywords to flag in log messages")
}
