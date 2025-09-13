package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromptManager(t *testing.T) {
	pm, err := NewPromptManager()
	if err != nil {
		t.Fatalf("Failed to create prompt manager: %v", err)
	}

	// Test system prompt template with new structure
	data := struct {
		ReadmeContent string
		CommandOutput map[string]string
		FileContent   map[string]string
		Mode          string
	}{
		ReadmeContent: "Test README content",
		CommandOutput: map[string]string{
			"echo test": "test output",
		},
		FileContent: map[string]string{
			"README.md": "Test README content",
			"main.go":   "package main",
		},
		Mode: "explicit",
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

	if !strings.Contains(prompt, "echo test") {
		t.Error("System prompt should contain command output")
	}

	if !strings.Contains(prompt, "package main") {
		t.Error("System prompt should contain file content")
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

func TestPromptManagerWithCustomTemplates(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "reporadio-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .reporadio/prompts directory
	promptsDir := filepath.Join(".reporadio", "prompts")
	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create a custom system prompt template
	customSystemPrompt := `Custom System Prompt
ReadmeContent: {{.ReadmeContent}}
This is a custom template for testing.`

	customSystemPromptPath := filepath.Join(promptsDir, "system_prompt.tmpl")
	err = os.WriteFile(customSystemPromptPath, []byte(customSystemPrompt), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom system prompt: %v", err)
	}

	// Test prompt manager loads custom templates
	pm, err := NewPromptManager()
	if err != nil {
		t.Fatalf("Failed to create prompt manager: %v", err)
	}

	// Test that system prompt uses custom template
	data := struct {
		ReadmeContent string
	}{
		ReadmeContent: "Test README content",
	}

	prompt, err := pm.Execute("system_prompt.tmpl", data)
	if err != nil {
		t.Fatalf("Failed to execute custom system_prompt.tmpl: %v", err)
	}

	if !strings.Contains(prompt, "Custom System Prompt") {
		t.Error("Should use custom system prompt template")
	}

	if !strings.Contains(prompt, "Test README content") {
		t.Error("Custom template should still process data correctly")
	}

	// Test that embedded templates are still available for non-overridden templates
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
		t.Fatalf("Failed to execute embedded episode_transcript.tmpl: %v", err)
	}

	if episodePrompt == "" {
		t.Error("Episode prompt is empty")
	}

}

func TestPromptManagerFallbackToEmbedded(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "reporadio-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .reporadio/prompts directory but only with some templates
	promptsDir := filepath.Join(".reporadio", "prompts")
	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create only a custom system prompt, leaving others to fallback
	customSystemPrompt := `Custom System Prompt: {{.ReadmeContent}}`
	customSystemPromptPath := filepath.Join(promptsDir, "system_prompt.tmpl")
	err = os.WriteFile(customSystemPromptPath, []byte(customSystemPrompt), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom system prompt: %v", err)
	}

	// Test prompt manager loads custom and falls back to embedded
	pm, err := NewPromptManager()
	if err != nil {
		t.Fatalf("Failed to create prompt manager: %v", err)
	}

	// Test custom system prompt is used
	data := struct {
		ReadmeContent string
	}{
		ReadmeContent: "Test README",
	}

	prompt, err := pm.Execute("system_prompt.tmpl", data)
	if err != nil {
		t.Fatalf("Failed to execute custom system_prompt.tmpl: %v", err)
	}

	if !strings.Contains(prompt, "Custom System Prompt: Test README") {
		t.Error("Should use custom system prompt")
	}

	// Test that other templates fallback to embedded versions
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
		t.Fatalf("Failed to execute embedded episode_transcript.tmpl: %v", err)
	}

	if episodePrompt == "" {
		t.Error("Episode prompt should fallback to embedded version")
	}
}
