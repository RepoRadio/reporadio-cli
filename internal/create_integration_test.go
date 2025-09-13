package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateIntegration_SystemPromptWithContext(t *testing.T) {
	// Setup test directory with config and files
	tempDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .reporadio directory and config
	repoRadioDir := filepath.Join(tempDir, ".reporadio")
	os.MkdirAll(repoRadioDir, 0755)

	configContent := `include:
  - "README.md"
  - "main.go"
commands:
  - "echo 'Hello from command'"
  - "pwd"`

	configPath := filepath.Join(repoRadioDir, "create.yml")
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Create test files
	os.WriteFile("README.md", []byte("# Test Project\nThis is a test"), 0644)
	os.WriteFile("main.go", []byte("package main\n\nfunc main() {}"), 0644)

	// Load config and build context
	config, mode, err := loadCreateConfigWithFallback()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if mode != ModeExplicit {
		t.Errorf("Expected mode %s, got %s", ModeExplicit, mode)
	}

	contextBuilder := NewContextBuilder()
	context, err := contextBuilder.BuildCreateContext(config, mode)
	if err != nil {
		t.Fatalf("Failed to build context: %v", err)
	}

	// Generate system prompt
	systemPrompt := getSystemPromptWithContext(context)

	// Verify the prompt contains expected content
	if !strings.Contains(systemPrompt, "# Test Project") {
		t.Error("System prompt should contain README content")
	}

	if !strings.Contains(systemPrompt, "Hello from command") {
		t.Error("System prompt should contain command output")
	}

	if !strings.Contains(systemPrompt, "package main") {
		t.Error("System prompt should contain main.go content")
	}

	if !strings.Contains(systemPrompt, "File: main.go") {
		t.Error("System prompt should label file content")
	}

	// Verify command output structure
	if !strings.Contains(systemPrompt, "Command: echo 'Hello from command'") {
		t.Error("System prompt should label command output")
	}
}

func TestCreateIntegration_DefaultConfigFallback(t *testing.T) {
	// Setup test directory without config
	tempDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create only a README
	os.WriteFile("README.md", []byte("# Default Project"), 0644)

	// Load config (should fallback to default)
	config, mode, err := loadCreateConfigWithFallback()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if mode != ModeDefault {
		t.Errorf("Expected mode %s, got %s", ModeDefault, mode)
	}

	contextBuilder := NewContextBuilder()
	context, err := contextBuilder.BuildCreateContext(config, mode)
	if err != nil {
		t.Fatalf("Failed to build context: %v", err)
	}

	// Generate system prompt
	systemPrompt := getSystemPromptWithContext(context)

	// Should contain README content from default "." include
	if !strings.Contains(systemPrompt, "# Default Project") {
		t.Error("System prompt should contain README content with default config")
	}

	// Should not contain command output (default has no commands)
	if strings.Contains(systemPrompt, "Here's output from setup commands") {
		t.Error("System prompt should not contain command output section with default config")
	}
}
