package analyzer

import (
	"reflect"
	"testing"
)

func TestSettingsToConfig(t *testing.T) {
	t.Parallel()

	falseValue := false
	settings := Settings{
		Rules: RuleSettings{
			LowercaseStart: &falseValue,
		},
		SensitiveData: SensitiveDataSettings{
			AdditionalKeywords: []string{" Credential ", "session_id", "credential", ""},
		},
	}

	config := settings.ToConfig()
	if config.Rules.LowercaseStart {
		t.Fatal("LowercaseStart = true, want false")
	}

	if !config.Rules.ASCIIOnly || !config.Rules.NoSpecialChars || !config.Rules.SensitiveData {
		t.Fatal("unspecified rules should keep default enabled state")
	}

	wantKeywords := []string{"credential", "session_id"}
	if !reflect.DeepEqual(config.SensitiveData.AdditionalKeywords, wantKeywords) {
		t.Fatalf("AdditionalKeywords = %#v, want %#v", config.SensitiveData.AdditionalKeywords, wantKeywords)
	}

	patterns := config.sensitivePatterns()
	if len(patterns) != len(defaultSensitiveKeywords)+len(wantKeywords) {
		t.Fatalf("sensitivePatterns count = %d, want %d", len(patterns), len(defaultSensitiveKeywords)+len(wantKeywords))
	}
}

func TestKeywordListFlagSet(t *testing.T) {
	t.Parallel()

	var list keywordListFlag
	if err := list.Set("credential, session_id, credential"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	want := []string{"credential", "session_id"}
	if !reflect.DeepEqual(list.values, want) {
		t.Fatalf("values = %#v, want %#v", list.values, want)
	}
}
