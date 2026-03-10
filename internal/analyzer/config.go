package analyzer

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

var defaultConfig = Config{
	Rules: RuleConfig{
		LowercaseStart: true,
		ASCIIOnly:      true,
		NoSpecialChars: true,
		SensitiveData:  true,
	},
}

func (config Config) normalized() Config {
	normalized := config

	normalized.SensitiveData.AdditionalKeywords = append(
		[]string{},
		normalized.SensitiveData.AdditionalKeywords...,
	)

	return normalized
}