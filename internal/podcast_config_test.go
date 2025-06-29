package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPodcastConfigWithCommands(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()
	podcastDir := filepath.Join(tempDir, ".reporadio", "test-podcast")
	err := os.MkdirAll(podcastDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test podcast.yml with commands
	yamlContent := `episodes:
  - title: "Test Episode with Commands"
    description: "Testing dynamic command execution"
    instructions: "Generate a transcript including command output"
    voicing: "friendly and informative"
    include:
      - "README.md"
    commands:
      - "echo 'Hello World'"
      - "date"
      - "pwd"
  - title: "Test Episode without Commands"
    description: "Testing backward compatibility"
    instructions: "Generate a regular transcript"
    voicing: "professional"
    include:
      - "LICENSE"
`

	configPath := filepath.Join(podcastDir, "podcast.yml")
	err = os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to temp directory for loadPodcastConfig to work
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

	// Verify we have two episodes
	if len(config.Episodes) != 2 {
		t.Fatalf("Expected 2 episodes, got %d", len(config.Episodes))
	}

	// Verify first episode has commands
	episode1 := config.Episodes[0]
	if episode1.Title != "Test Episode with Commands" {
		t.Errorf("Expected title 'Test Episode with Commands', got %q", episode1.Title)
	}

	if len(episode1.Commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(episode1.Commands))
	}

	expectedCommands := []string{"echo 'Hello World'", "date", "pwd"}
	for i, expectedCmd := range expectedCommands {
		if episode1.Commands[i] != expectedCmd {
			t.Errorf("Expected command %q, got %q", expectedCmd, episode1.Commands[i])
		}
	}

	// Verify second episode has no commands (backward compatibility)
	episode2 := config.Episodes[1]
	if episode2.Title != "Test Episode without Commands" {
		t.Errorf("Expected title 'Test Episode without Commands', got %q", episode2.Title)
	}

	if len(episode2.Commands) != 0 {
		t.Errorf("Expected no commands for second episode, got %d", len(episode2.Commands))
	}
}

func TestLoadPodcastConfigInvalidYAML(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()
	podcastDir := filepath.Join(tempDir, ".reporadio", "invalid-podcast")
	err := os.MkdirAll(podcastDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create an invalid YAML file
	invalidYAML := `episodes:
  - title: "Test Episode"
    description: "Testing invalid YAML"
    commands:
      - echo 'missing quotes
    invalid_field: [
`

	configPath := filepath.Join(podcastDir, "podcast.yml")
	err = os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to temp directory
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

	// Try to load the invalid configuration
	_, err = loadPodcastConfig("invalid-podcast")
	if err == nil {
		t.Error("Expected error when loading invalid YAML, got nil")
	}
}
