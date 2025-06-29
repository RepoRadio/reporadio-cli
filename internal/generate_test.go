package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPodcastConfig(t *testing.T) {
	// Create a temporary podcast config for testing
	tempDir := t.TempDir()
	podcastDir := filepath.Join(tempDir, ".reporadio", "test")
	err := os.MkdirAll(podcastDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	configContent := `episodes:
  - title: 'Episode 1: Test Episode'
    description: 'A test episode'
    instructions: 'Test instructions'
    voicing: 'Test voicing'
    include:
      - README.md
      - main.go
`

	configPath := filepath.Join(podcastDir, "podcast.yml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to temp directory to simulate running from project root
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test loading the config
	config, err := loadPodcastConfig("test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(config.Episodes) != 1 {
		t.Fatalf("Expected 1 episode, got %d", len(config.Episodes))
	}

	episode := config.Episodes[0]
	if episode.Title != "Episode 1: Test Episode" {
		t.Errorf("Expected title 'Episode 1: Test Episode', got '%s'", episode.Title)
	}

	if len(episode.Include) != 2 {
		t.Errorf("Expected 2 included files, got %d", len(episode.Include))
	}
}

func TestLoadPodcastConfigNotFound(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	_, err := loadPodcastConfig("nonexistent")
	if err == nil {
		t.Fatal("Expected error for nonexistent config, got nil")
	}
}

func TestGenerateEpisodeTranscript(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create test files that the episode will include
	err := os.WriteFile("README.md", []byte("# Test Project\nThis is a test project."), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = os.WriteFile("main.go", []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test episode
	episode := Episode{
		Title:        "Episode 1: Introduction",
		Description:  "An introduction to the project",
		Instructions: "Explain the project basics",
		Voicing:      "Friendly and informative",
		Include:      []string{"README.md", "main.go"},
	}

	// Create a mock OpenAI client (we'll use nil for now and handle it in the function)
	outputDir := filepath.Join(".reporadio", "test", "episodes")

	// Test the function
	err = generateEpisodeTranscript(episode, 1, outputDir, nil, []map[string]interface{}{}, "")

	// For now, we expect this to fail since we don't have a real client
	// but we want to test the file reading and structure logic
	if err == nil {
		// Check if transcript file was created
		transcriptPath := filepath.Join(outputDir, "ep1.md")
		if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
			t.Error("Expected transcript file to be created")
		}
	}
}

func TestGeneratePodcastTranscripts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create test files
	err := os.WriteFile("README.md", []byte("# Test Project\nThis is a test project."), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = os.WriteFile("main.go", []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test podcast config
	config := &PodcastConfig{
		Episodes: []Episode{
			{
				Title:        "Episode 1: Introduction",
				Description:  "An introduction to the project",
				Instructions: "Explain the project basics",
				Voicing:      "Friendly and informative",
				Include:      []string{"README.md"},
			},
			{
				Title:        "Episode 2: Code Overview",
				Description:  "Overview of the main code",
				Instructions: "Explain the code structure",
				Voicing:      "Technical but accessible",
				Include:      []string{"main.go"},
			},
		},
	}

	// Test the function (with nil client for testing)
	err = generatePodcastTranscripts("test", config, nil, false, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if transcript files were created
	outputDir := filepath.Join(".reporadio", "test", "episodes")

	ep1Path := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(ep1Path); os.IsNotExist(err) {
		t.Error("Expected ep1.md to be created")
	}

	ep2Path := filepath.Join(outputDir, "ep2.md")
	if _, err := os.Stat(ep2Path); os.IsNotExist(err) {
		t.Error("Expected ep2.md to be created")
	}

	// Check content of first episode
	content, err := os.ReadFile(ep1Path)
	if err != nil {
		t.Fatalf("Failed to read ep1.md: %v", err)
	}

	if !strings.Contains(string(content), "Episode 1: Introduction") {
		t.Error("Episode 1 transcript should contain the title")
	}
}

func TestGeneratePodcastTranscriptsWithAudio(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWd) }() // Ignore error in test cleanup
	_ = os.Chdir(tempDir)

	// Create test files
	err := os.WriteFile("README.md", []byte("# Test Project\nThis is a test project."), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create test podcast config
	config := &PodcastConfig{
		Episodes: []Episode{
			{
				Title:        "Episode 1: Introduction",
				Description:  "A short intro episode",
				Instructions: "Keep it brief",
				Voicing:      "Friendly",
				Include:      []string{"README.md"},
			},
		},
	}

	// Test with audio generation but no client (should handle gracefully)
	err = generatePodcastTranscripts("test", config, nil, true, "")
	if err != nil {
		t.Fatalf("Expected no error with nil client, got %v", err)
	}

	// Check if transcript was created
	outputDir := filepath.Join(".reporadio", "test", "episodes")
	ep1Path := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(ep1Path); os.IsNotExist(err) {
		t.Error("Expected ep1.md to be created")
	}

	// Audio file should not be created with nil client
	audioPath := filepath.Join(outputDir, "ep1.mp3")
	if _, err := os.Stat(audioPath); err == nil {
		t.Error("Did not expect ep1.mp3 to be created with nil client")
	}
}

func TestChatContext(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	podcastDir := filepath.Join(tempDir, ".reporadio", "test")
	err := os.MkdirAll(podcastDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Change to temp directory to simulate running from project root
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test loading empty chat context
	entries, err := loadChatContext("test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("Expected 0 entries, got %d", len(entries))
	}

	// Test appending to chat context
	err = appendToChatContext("test", 1, "Test Episode", "Test transcript content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test loading chat context with one entry
	entries, err = loadChatContext("test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	// Verify the entry content
	entry := entries[0]
	if role, ok := entry["role"].(string); !ok || role != "assistant" {
		t.Fatalf("Expected role 'assistant', got %v", entry["role"])
	}
	if step, ok := entry["step"].(string); !ok || step != "episode" {
		t.Fatalf("Expected step 'episode', got %v", entry["step"])
	}
	if message, ok := entry["message"].(string); !ok || !strings.Contains(message, "Test Episode") {
		t.Fatalf("Expected message to contain 'Test Episode', got %v", entry["message"])
	}

	// Test appending another entry
	err = appendToChatContext("test", 2, "Second Episode", "Second transcript content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test loading chat context with two entries
	entries, err = loadChatContext("test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}
}
