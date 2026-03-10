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