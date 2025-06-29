package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDynamicCommandExecutionEndToEnd(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	podcastDir := filepath.Join(tempDir, ".reporadio", "test-podcast")
	err := os.MkdirAll(podcastDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test files
	readmeContent := "# Test Project\n\nThis is a demonstration project for dynamic command execution."
	readmePath := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	// Create a test script
	scriptContent := `#!/bin/bash
echo "Script executed successfully"
echo "Current directory: $(pwd)"
echo "Files in directory:"
ls -la | head -5
`
	scriptPath := filepath.Join(tempDir, "test-script.sh")
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Create a comprehensive podcast.yml with commands
	yamlContent := `episodes:
  - title: "Dynamic Content Showcase"
    description: "Demonstrating dynamic command execution in podcast generation"
    instructions: "Create an engaging transcript that incorporates both static file content and dynamic command output"
    voicing: "enthusiastic and technical"
    include:
      - "README.md"
    commands:
      - "echo 'Welcome to the dynamic content demo!'"
      - "date"
      - "echo 'Repository statistics:'"
      - "ls -la | wc -l"
      - "./test-script.sh"
      - "echo 'End of dynamic content'"
`

	configPath := filepath.Join(podcastDir, "podcast.yml")
	err = os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to temp directory for the test
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

	// Load the podcast configuration
	config, err := loadPodcastConfig("test-podcast")
	if err != nil {
		t.Fatalf("Failed to load podcast config: %v", err)
	}

	// Verify configuration loaded correctly
	if len(config.Episodes) != 1 {
		t.Fatalf("Expected 1 episode, got %d", len(config.Episodes))
	}

	episode := config.Episodes[0]
	if len(episode.Commands) != 6 {
		t.Fatalf("Expected 6 commands, got %d", len(episode.Commands))
	}

	// Create mock podcast config
	podcastConfig := &PodcastConfig{
		Title:        "Test Podcast",
		Description:  "Test Description",
		Instructions: "Test Instructions",
		Voicing:      "friendly",
		Type:         SeriesTypeOnboarding,
		Episodes:     []Episode{},
	}

	// Generate the episode transcript
	outputDir := filepath.Join(tempDir, "output")
	err = generateEpisodeTranscript(episode, 1, outputDir, nil, []map[string]interface{}{}, podcastConfig, "30s")
	if err != nil {
		t.Fatalf("Failed to generate episode transcript: %v", err)
	}

	// Verify output file was created
	transcriptPath := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		t.Fatal("Transcript file was not created")
	}

	// Read and verify the generated transcript
	transcriptContent, err := os.ReadFile(transcriptPath)
	if err != nil {
		t.Fatalf("Failed to read transcript: %v", err)
	}

	content := string(transcriptContent)

	// Verify the transcript contains expected static content
	if !strings.Contains(content, "Dynamic Content Showcase") {
		t.Error("Transcript should contain the episode title")
	}

	// Note: Since this is a placeholder transcript (nil client),
	// we can't verify the actual command output integration in the final content.
	// But the test verifies that:
	// 1. Configuration loading works with commands
	// 2. Command execution integration doesn't break episode generation
	// 3. Files are created successfully

	t.Logf("Successfully generated episode transcript with dynamic commands")
	t.Logf("Episode title: %s", episode.Title)
	t.Logf("Commands configured: %d", len(episode.Commands))
	t.Logf("Transcript file created: %s", transcriptPath)
}
