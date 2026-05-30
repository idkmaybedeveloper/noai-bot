package main

import (
	"regexp"
	"strings"
)

// Co-Authored-By trailer format: "Name <email>"
// two strategies:
// exact email match against confirmed bot addresses
// [bot]@users.noreply.github.com regex catches any future Gh App bots

var knownAIEmails = []string{
	// claude code - ref: github.com/anthropics/claude-code also https://github.com/golang/go/pull/77645/changes/2089b9cc6a5184cad20c31d87dc5bc0d3e1c0288
	"noreply@anthropic.com",

	// gh copilot (ide via git.addAICoAuthor setting)
	"copilot@github.com",
	"223556219+copilot@users.noreply.github.com",

	// gh copilot swe agent (ref: docs+sweagentd)
	"203248971+copilot-swe-agent@users.noreply.github.com",

	// cursor
	"cursoragent@cursor.com",

	// openai codex cli
	"noreply@openai.com",

	// github code assist (GitHub App, enterprise)
	"176961590+gemini-code-assist[bot]@users.noreply.github.com",

	// github cli (npm google/gemini-cli)
	"218195315+gemini-cli@users.noreply.github.com",

	// Devin (Cognition AI)
	"158243242+devin-ai-integration[bot]@users.noreply.github.com",

	// Aider (--attribute-co-authored-by, on by default)
	"noreply@aider.chat",
	// Jules (Google) - user ID unconfirmed, covered by reBotNoreply for now
}

var knownAINames = []string{
	"github copilot",
	"claude code",
	"cursor ai",
	"gemini code assist",
	"gemini cli",
	"devin",
}

// catches any Gh App bot using the standard noreply pattern
var reBotNoreply = regexp.MustCompile(`\[bot\]@users\.noreply\.github\.com$`)

func matchesAI(coAuthor string, extraPatterns []string) bool {
	lower := strings.ToLower(coAuthor)
	email := extractEmail(lower)

	for _, e := range knownAIEmails {
		if email == e {
			return true
		}
	}

	if reBotNoreply.MatchString(email) {
		return true
	}

	for _, n := range knownAINames {
		if strings.Contains(lower, n) {
			return true
		}
	}

	for _, p := range extraPatterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}

	return false
}

func extractEmail(trailer string) string {
	start := strings.LastIndex(trailer, "<")
	end := strings.LastIndex(trailer, ">")
	if start == -1 || end == -1 || end < start {
		return trailer
	}
	return trailer[start+1 : end]
}
