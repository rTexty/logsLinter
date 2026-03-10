package analyzer

import "testing"

func TestCheckLowercaseStart(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		text    string
		want    bool
		ruleID  string
		message string
	}{
		{
			name: "empty message is allowed",
			text: "",
			want: false,
		},
		{
			name: "lowercase start passes",
			text: "starting server",
			want: false,
		},
		{
			name: "uppercase start fails",
			text: "Starting server",
			want: true,
			ruleID: ruleLowercaseStart,
			message: msgLowercaseStart,
		},
		{
			name: "non-letter first rune is ignored",
			text: "123 started",
			want: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, ok := checkLowercaseStart(testCase.text)
			if ok != testCase.want {
				t.Fatalf("checkLowercaseStart(%q) ok = %v, want %v", testCase.text, ok, testCase.want)
			}

			if !testCase.want {
				return
			}

			if got.ruleID != testCase.ruleID {
				t.Fatalf("checkLowercaseStart(%q) ruleID = %q, want %q", testCase.text, got.ruleID, testCase.ruleID)
			}

			if got.message != testCase.message {
				t.Fatalf("checkLowercaseStart(%q) message = %q, want %q", testCase.text, got.message, testCase.message)
			}
		})
	}
}

func TestCheckASCIIOnly(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		text    string
		want    bool
		ruleID  string
		message string
	}{
		{
			name: "plain ascii passes",
			text: "failed to connect",
			want: false,
		},
		{
			name: "mixed ascii and cyrillic fails",
			text: "failed ошибка",
			want: true,
			ruleID: ruleASCIIOnly,
			message: msgASCIIOnly,
		},
		{
			name: "cyrillic fails",
			text: "ошибка подключения",
			want: true,
			ruleID: ruleASCIIOnly,
			message: msgASCIIOnly,
		},
		{
			name: "control characters fail",
			text: "line\nbreak",
			want: true,
			ruleID: ruleASCIIOnly,
			message: msgASCIIOnly,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, ok := checkASCIIOnly(testCase.text)
			if ok != testCase.want {
				t.Fatalf("checkASCIIOnly(%q) ok = %v, want %v", testCase.text, ok, testCase.want)
			}

			if !testCase.want {
				return
			}

			if got.ruleID != testCase.ruleID {
				t.Fatalf("checkASCIIOnly(%q) ruleID = %q, want %q", testCase.text, got.ruleID, testCase.ruleID)
			}

			if got.message != testCase.message {
				t.Fatalf("checkASCIIOnly(%q) message = %q, want %q", testCase.text, got.message, testCase.message)
			}
		})
	}
}

func TestCheckNoSpecialCharsOrEmoji(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		text    string
		want    bool
		ruleID  string
		message string
	}{
		{
			name: "plain message passes",
			text: "connection established",
			want: false,
		},
		{
			name: "exclamation mark fails",
			text: "server started!",
			want: true,
			ruleID: ruleNoSpecialChars,
			message: msgNoSpecialChars,
		},
		{
			name: "trailing question mark fails",
			text: "connection lost?",
			want: true,
			ruleID: ruleNoSpecialChars,
			message: msgNoSpecialChars,
		},
		{
			name: "ellipsis fails",
			text: "waiting...",
			want: true,
			ruleID: ruleNoSpecialChars,
			message: msgNoSpecialChars,
		},
		{
			name: "emoji fails",
			text: "server started 🚀",
			want: true,
			ruleID: ruleNoSpecialChars,
			message: msgNoSpecialChars,
		},
		{
			name: "hyphenated technical term passes",
			text: "connecting to db-host",
			want: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, ok := checkNoSpecialCharsOrEmoji(testCase.text)
			if ok != testCase.want {
				t.Fatalf("checkNoSpecialCharsOrEmoji(%q) ok = %v, want %v", testCase.text, ok, testCase.want)
			}

			if !testCase.want {
				return
			}

			if got.ruleID != testCase.ruleID {
				t.Fatalf("checkNoSpecialCharsOrEmoji(%q) ruleID = %q, want %q", testCase.text, got.ruleID, testCase.ruleID)
			}

			if got.message != testCase.message {
				t.Fatalf("checkNoSpecialCharsOrEmoji(%q) message = %q, want %q", testCase.text, got.message, testCase.message)
			}
		})
	}
}

func TestCheckSensitiveData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		sample  messageSample
		want    bool
		ruleID  string
		message string
	}{
		{
			name: "safe authentication wording passes",
			sample: messageSample{text: "user authenticated successfully"},
			want: false,
		},
		{
			name: "password keyword fails",
			sample: messageSample{text: "user password invalid"},
			want: true,
			ruleID: ruleSensitiveData,
			message: msgSensitiveData,
		},
		{
			name: "substring false positive is ignored",
			sample: messageSample{text: "oauth flow started"},
			want: false,
		},
		{
			name: "api key keyword fails",
			sample: messageSample{text: "api_key provided"},
			want: true,
			ruleID: ruleSensitiveData,
			message: msgSensitiveData,
		},
		{
			name: "literal concatenation parts fail",
			sample: messageSample{
				text:  "user password: ",
				parts: []string{"user password: ", ""},
			},
			want: true,
			ruleID: ruleSensitiveData,
			message: msgSensitiveData,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, ok := checkSensitiveData(testCase.sample)
			if ok != testCase.want {
				t.Fatalf("checkSensitiveData(%+v) ok = %v, want %v", testCase.sample, ok, testCase.want)
			}

			if !testCase.want {
				return
			}

			if got.ruleID != testCase.ruleID {
				t.Fatalf("checkSensitiveData(%+v) ruleID = %q, want %q", testCase.sample, got.ruleID, testCase.ruleID)
			}

			if got.message != testCase.message {
				t.Fatalf("checkSensitiveData(%+v) message = %q, want %q", testCase.sample, got.message, testCase.message)
			}
		})
	}
}