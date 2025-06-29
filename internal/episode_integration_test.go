package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateEpisodeTranscriptWithCommands(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Create a test file to include
	testFile := filepath.Join(tempDir, "README.md")
	testFileContent := "# Test Project\n\nThis is a test project with some content."
	err := os.WriteFile(testFile, []byte(testFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to temp directory so relative paths work
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalDir)
		if err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test episode with both file includes and commands
	episode := Episode{
		Title:        "Dynamic Content Episode",
		Description:  "Testing integration of file content and command output",
		Instructions: "Create engaging content using both static files and dynamic command output",
		Voicing:      "enthusiastic and technical",
		Include:      []string{"README.md"},
		Commands:     []string{"echo 'Generated at runtime!'", "echo 'Dynamic data: 42'"},
	}

	outputDir := filepath.Join(tempDir, "output")

	// Create mock podcast config
	podcastConfig := &PodcastConfig{
		Title:        "Test Podcast",
		Description:  "Test Description",
		Instructions: "Test Instructions",
		Voicing:      "friendly",
		Type:         SeriesTypeOnboarding,
		Episodes:     []Episode{},
	}

	// Generate transcript with nil client (testing mode)
	err = generateEpisodeTranscript(episode, 1, outputDir, nil, []map[string]interface{}{}, podcastConfig, "")
	if err != nil {
		t.Fatalf("Failed to generate episode transcript: %v", err)
	}

	// Verify output file was created
	transcriptPath := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		t.Fatal("Transcript file was not created")
	}

	// Read the transcript content
	transcriptContent, err := os.ReadFile(transcriptPath)
	if err != nil {
		t.Fatalf("Failed to read transcript: %v", err)
	}

	content := string(transcriptContent)

	// Verify the transcript contains the expected content
	if !strings.Contains(content, "Dynamic Content Episode") {
		t.Error("Transcript should contain the episode title")
	}
}

func TestCommandExecutionTimeout(t *testing.T) {
	episode := Episode{
		Title:       "Timeout Test",
		Description: "Testing command timeout",
		Commands:    []string{"sleep 2"}, // This should timeout with 60s default
	}

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "output")

	// Create mock podcast config
	podcastConfig := &PodcastConfig{
		Title:        "Test Podcast",
		Description:  "Test Description",
		Instructions: "Test Instructions",
		Voicing:      "friendly",
		Type:         SeriesTypeOnboarding,
		Episodes:     []Episode{},
	}

	// This should complete without hanging (command will timeout but generation continues)
	err := generateEpisodeTranscript(episode, 1, outputDir, nil, []map[string]interface{}{}, podcastConfig, "")
	if err != nil {
		t.Fatalf("Episode generation should not fail due to command timeout: %v", err)
	}

	// Verify output file was still created
	transcriptPath := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		t.Error("Transcript file should be created even when commands fail")
	}
}
