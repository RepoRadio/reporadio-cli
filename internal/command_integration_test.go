package internal

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIntegrateCommandOutput(t *testing.T) {
	// Test episode with commands
	episode := Episode{
		Title:        "Test Episode",
		Description:  "Testing command integration",
		Instructions: "Include command output",
		Voicing:      "technical",
		Include:      []string{},
		Commands:     []string{"echo 'test output'", "echo 'second command'"},
	}

	// Execute commands and get formatted output
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := executeCommands(ctx, episode.Commands)
	commandOutput := formatCommandOutput(results)

	// Verify command output is formatted correctly
	if !strings.Contains(commandOutput, "=== Command: echo 'test output' ===") {
		t.Error("Command output should contain header for first command")
	}
	if !strings.Contains(commandOutput, "test output") {
		t.Error("Command output should contain result of first command")
	}
	if !strings.Contains(commandOutput, "=== Command: echo 'second command' ===") {
		t.Error("Command output should contain header for second command")
	}
	if !strings.Contains(commandOutput, "second command") {
		t.Error("Command output should contain result of second command")
	}
}

func TestGenerateEpisodeWithCommands(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Create a test file to include
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test episode with both includes and commands
	episode := Episode{
		Title:        "Test Episode with Commands",
		Description:  "Testing both file includes and command execution",
		Instructions: "Generate transcript with file content and command output",
		Voicing:      "friendly",
		Include:      []string{testFile},
		Commands:     []string{"echo 'dynamic content'", "date"},
	}

	outputDir := filepath.Join(tempDir, "output")

	// Test with nil client (placeholder mode)
	err = generateEpisodeTranscript(episode, 1, outputDir, nil, []map[string]interface{}{}, "")
	if err != nil {
		t.Fatalf("Failed to generate episode transcript: %v", err)
	}

	// Verify output file was created
	transcriptPath := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		t.Error("Transcript file was not created")
	}
}
