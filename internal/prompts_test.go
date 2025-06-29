package internal

import (
	"strings"
	"testing"
)

func TestPromptManager(t *testing.T) {
	pm, err := NewPromptManager()
	if err != nil {
		t.Fatalf("Failed to create prompt manager: %v", err)
	}

	// Test system prompt template
	data := struct {
		ReadmeContent string
	}{
		ReadmeContent: "Test README content",
	}

	prompt, err := pm.Execute("system_prompt.tmpl", data)
	if err != nil {
		t.Fatalf("Failed to execute system_prompt.tmpl: %v", err)
	}

	if prompt == "" {
		t.Error("System prompt is empty")
	}

	if !strings.Contains(prompt, "Test README content") {
		t.Error("System prompt should contain README content")
	}

	// Test episode transcript template
	episodeData := struct {
		Title         string
		Description   string
		Instructions  string
		Voicing       string
		Context       string
		FileContents  string
		CommandOutput string
	}{
		Title:         "Test Episode",
		Description:   "Test Description",
		Instructions:  "Test Instructions",
		Voicing:       "Conversational",
		Context:       "Previous episode context",
		FileContents:  "Test file contents",
		CommandOutput: "=== Command: echo test ===\ntest output",
	}

	episodePrompt, err := pm.Execute("episode_transcript.tmpl", episodeData)
	if err != nil {
		t.Fatalf("Failed to execute episode_transcript.tmpl: %v", err)
	}

	if episodePrompt == "" {
		t.Error("Episode prompt is empty")
	}

	if !strings.Contains(episodePrompt, "Test Episode") {
		t.Error("Episode prompt should contain episode title")
	}
}
