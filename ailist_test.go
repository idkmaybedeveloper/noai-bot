package main

import "testing"

func TestMatchesAI(t *testing.T) {
	cases := []struct {
		trailer string
		want    bool
	}{
		// claude code
		{"Claude Sonnet 4.6 <noreply@anthropic.com>", true},
		// copilot (ide)
		{"Copilot <copilot@github.com>", true},
		{"Copilot <223556219+Copilot@users.noreply.github.com>", true},
		// copilot swe agent
		{"copilot-swe-agent <203248971+copilot-swe-agent@users.noreply.github.com>", true},
		// cursor
		{"Cursor <cursoragent@cursor.com>", true},
		// openai codex cli
		{"Codex <noreply@openai.com>", true},
		// gemini code assist (github app)
		{"gemini-code-assist[bot] <176961590+gemini-code-assist[bot]@users.noreply.github.com>", true},
		// gemini cli
		{"gemini-cli <218195315+gemini-cli@users.noreply.github.com>", true},
		// devin
		{"devin-ai-integration[bot] <158243242+devin-ai-integration[bot]@users.noreply.github.com>", true},
		// aider
		{"aider (claude-3-5-sonnet-20241022) <noreply@aider.chat>", true},
		// generic [bot]@users.noreply.github.com catch-all (e.g. jules or any future bot)
		{"jules-google <235795819+jules-google[bot]@users.noreply.github.com>", true},
		{"SomeFutureBot <999999+somefuture[bot]@users.noreply.github.com>", true},
		// extra patterns from config
		{"MyCorp AI <bot@mycorp.internal>", true},
		// real humans
		{"John Doe <john@example.com>", false},
		{"Jane Smith <jane.smith@corp.dev>", false},
		// tools that don't add co-authored-by (no match expected)
		{"Tabnine <tabnine@tabnine.com>", false},
		{"Amazon Q <amazonq@amazon.com>", false},
	}

	extra := []string{"bot@mycorp.internal"}

	for _, tc := range cases {
		got := matchesAI(tc.trailer, extra)
		if got != tc.want {
			t.Errorf("matchesAI(%q) = %v, want %v", tc.trailer, got, tc.want)
		}
	}
}

func TestExtractEmail(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"some name <user@example.com>", "user@example.com"},
		{"no angle brackets", "no angle brackets"},
		{"<only-email@x.com>", "only-email@x.com"},
		{"weird <<double@x.com>", "double@x.com"},
		{"aider (claude-3-5) <noreply@aider.chat>", "noreply@aider.chat"},
	}
	for _, tc := range cases {
		got := extractEmail(tc.in)
		if got != tc.want {
			t.Errorf("extractEmail(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
