package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlaylistEpisodeStructure(t *testing.T) {
	// Test the PlaylistEpisode struct and JSON marshaling
	episode := PlaylistEpisode{
		ID:             "ep1",
		Title:          "Test Episode",
		Description:    "A test episode description",
		AudioFile:      "ep1.mp3",
		TranscriptFile: "ep1.md",
		Summary:        "This is a test episode summary.",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(episode)
	if err != nil {
		t.Fatalf("Failed to marshal PlaylistEpisode: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledEpisode PlaylistEpisode
	err = json.Unmarshal(jsonData, &unmarshaledEpisode)
	if err != nil {
		t.Fatalf("Failed to unmarshal PlaylistEpisode: %v", err)
	}

	// Verify the unmarshaled data
	if unmarshaledEpisode.ID != episode.ID {
		t.Errorf("Expected ID %s, got %s", episode.ID, unmarshaledEpisode.ID)
	}
	if unmarshaledEpisode.Title != episode.Title {
		t.Errorf("Expected Title %s, got %s", episode.Title, unmarshaledEpisode.Title)
	}
	if unmarshaledEpisode.Summary != episode.Summary {
		t.Errorf("Expected Summary %s, got %s", episode.Summary, unmarshaledEpisode.Summary)
	}
}

func TestPlaylistStructure(t *testing.T) {
	// Test the Playlist and PlaylistRoot structures
	episodes := []PlaylistEpisode{
		{
			ID:             "ep1",
			Title:          "Episode 1",
			Description:    "First episode",
			AudioFile:      "ep1.mp3",
			TranscriptFile: "ep1.md",
			Summary:        "Summary of episode 1",
		},
		{
			ID:             "ep2",
			Title:          "Episode 2",
			Description:    "Second episode",
			AudioFile:      "ep2.mp3",
			TranscriptFile: "ep2.md",
			Summary:        "Summary of episode 2",
		},
	}

	playlist := Playlist{
		Title:       "Test Podcast",
		Description: "A test podcast description",
		Episodes:    episodes,
	}

	playlistRoot := PlaylistRoot{
		Playlist: playlist,
	}

	// Test JSON marshaling with proper indentation
	jsonData, err := json.MarshalIndent(playlistRoot, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal PlaylistRoot: %v", err)
	}

	// Verify the JSON structure contains expected fields
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "\"playlist\"") {
		t.Error("Expected JSON to contain 'playlist' field")
	}
	if !strings.Contains(jsonStr, "\"episodes\"") {
		t.Error("Expected JSON to contain 'episodes' field")
	}
	if !strings.Contains(jsonStr, "\"Test Podcast\"") {
		t.Error("Expected JSON to contain podcast title")
	}
	if !strings.Contains(jsonStr, "\"Episode 1\"") {
		t.Error("Expected JSON to contain episode title")
	}
}

func TestGenerateEpisodeSummaryWithNilClient(t *testing.T) {
	// Test generateEpisodeSummary with nil client
	transcript := "This is a test transcript for episode summary generation."
	
	summary, err := generateEpisodeSummary(nil, transcript)
	if err != nil {
		t.Fatalf("Expected no error with nil client, got %v", err)
	}
	
	expectedSummary := "Summary unavailable (no OpenAI client)"
	if summary != expectedSummary {
		t.Errorf("Expected summary '%s', got '%s'", expectedSummary, summary)
	}
}

func TestGenerateEpisodeSummaryTextCleanup(t *testing.T) {
	// Test that the summary cleanup works correctly
	// We can't easily test the OpenAI API without mocking, but we can test the text processing logic
	// by examining the function's behavior with different inputs
	
	// Test with nil client to verify the basic flow
	transcript := "This is a test transcript with\nmultiple lines\nand whitespace.   "
	
	summary, err := generateEpisodeSummary(nil, transcript)
	if err != nil {
		t.Fatalf("Expected no error with nil client, got %v", err)
	}
	
	// Verify the summary is cleaned up (no extra whitespace)
	if strings.HasPrefix(summary, " ") || strings.HasSuffix(summary, " ") {
		t.Error("Expected summary to be trimmed of whitespace")
	}
}

func TestGeneratePlaylistFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create test podcast config
	config := &PodcastConfig{
		Title:       "Test Podcast",
		Description: "A test podcast for playlist generation",
	}

	// Create test episode data
	episodeData := []PlaylistEpisode{
		{
			ID:             "ep1",
			Title:          "Episode 1: Introduction",
			Description:    "An introduction episode",
			AudioFile:      "ep1.mp3",
			TranscriptFile: "ep1.md",
			Summary:        "This episode introduces the podcast.",
		},
		{
			ID:             "ep2",
			Title:          "Episode 2: Deep dive",
			Description:    "A deep dive episode",
			AudioFile:      "ep2.mp3",
			TranscriptFile: "ep2.md",
			Summary:        "This episode goes deeper into the topic.",
		},
	}

	// Test generatePlaylistFile
	err := generatePlaylistFile("test", config, episodeData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that playlist.json was created
	playlistPath := filepath.Join(".reporadio", "test", "episodes", "playlist.json")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		t.Error("Expected playlist.json to be created")
	}

	// Read and verify the playlist.json content
	content, err := os.ReadFile(playlistPath)
	if err != nil {
		t.Fatalf("Failed to read playlist.json: %v", err)
	}

	// Verify JSON structure
	var playlistRoot PlaylistRoot
	err = json.Unmarshal(content, &playlistRoot)
	if err != nil {
		t.Fatalf("Failed to unmarshal playlist.json: %v", err)
	}

	// Verify playlist data
	if playlistRoot.Playlist.Title != config.Title {
		t.Errorf("Expected title '%s', got '%s'", config.Title, playlistRoot.Playlist.Title)
	}
	if playlistRoot.Playlist.Description != config.Description {
		t.Errorf("Expected description '%s', got '%s'", config.Description, playlistRoot.Playlist.Description)
	}
	if len(playlistRoot.Playlist.Episodes) != len(episodeData) {
		t.Errorf("Expected %d episodes, got %d", len(episodeData), len(playlistRoot.Playlist.Episodes))
	}

	// Verify episode data
	for i, episode := range playlistRoot.Playlist.Episodes {
		expected := episodeData[i]
		if episode.ID != expected.ID {
			t.Errorf("Episode %d: expected ID '%s', got '%s'", i, expected.ID, episode.ID)
		}
		if episode.Title != expected.Title {
			t.Errorf("Episode %d: expected title '%s', got '%s'", i, expected.Title, episode.Title)
		}
		if episode.Summary != expected.Summary {
			t.Errorf("Episode %d: expected summary '%s', got '%s'", i, expected.Summary, episode.Summary)
		}
	}
}

func TestGeneratePlaylistFileWithEmptyEpisodes(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create test podcast config
	config := &PodcastConfig{
		Title:       "Empty Podcast",
		Description: "A podcast with no episodes",
	}

	// Test with empty episode data
	episodeData := []PlaylistEpisode{}

	// Test generatePlaylistFile
	err := generatePlaylistFile("empty", config, episodeData)
	if err != nil {
		t.Fatalf("Expected no error with empty episodes, got %v", err)
	}

	// Check that playlist.json was created
	playlistPath := filepath.Join(".reporadio", "empty", "episodes", "playlist.json")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		t.Error("Expected playlist.json to be created even with empty episodes")
	}

	// Read and verify the playlist.json content
	content, err := os.ReadFile(playlistPath)
	if err != nil {
		t.Fatalf("Failed to read playlist.json: %v", err)
	}

	// Verify JSON structure
	var playlistRoot PlaylistRoot
	err = json.Unmarshal(content, &playlistRoot)
	if err != nil {
		t.Fatalf("Failed to unmarshal playlist.json: %v", err)
	}

	// Verify empty episodes array
	if len(playlistRoot.Playlist.Episodes) != 0 {
		t.Errorf("Expected 0 episodes, got %d", len(playlistRoot.Playlist.Episodes))
	}
}

func TestGeneratePodcastTranscriptsWithPlaylist(t *testing.T) {
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
		Title:       "Test Podcast with Playlist",
		Description: "A test podcast that generates playlist.json",
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

	// Check that playlist.json was created
	playlistPath := filepath.Join(".reporadio", "test", "episodes", "playlist.json")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		t.Error("Expected playlist.json to be created")
	}

	// Read and verify the playlist.json content
	content, err := os.ReadFile(playlistPath)
	if err != nil {
		t.Fatalf("Failed to read playlist.json: %v", err)
	}

	// Verify JSON structure
	var playlistRoot PlaylistRoot
	err = json.Unmarshal(content, &playlistRoot)
	if err != nil {
		t.Fatalf("Failed to unmarshal playlist.json: %v", err)
	}

	// Verify playlist metadata
	if playlistRoot.Playlist.Title != config.Title {
		t.Errorf("Expected title '%s', got '%s'", config.Title, playlistRoot.Playlist.Title)
	}
	if playlistRoot.Playlist.Description != config.Description {
		t.Errorf("Expected description '%s', got '%s'", config.Description, playlistRoot.Playlist.Description)
	}

	// Verify episode count
	if len(playlistRoot.Playlist.Episodes) != len(config.Episodes) {
		t.Errorf("Expected %d episodes, got %d", len(config.Episodes), len(playlistRoot.Playlist.Episodes))
	}

	// Verify episode data
	for i, episode := range playlistRoot.Playlist.Episodes {
		expectedID := "ep" + string(rune(i+1+'0'))
		if episode.ID != expectedID {
			t.Errorf("Episode %d: expected ID '%s', got '%s'", i, expectedID, episode.ID)
		}
		if episode.Title != config.Episodes[i].Title {
			t.Errorf("Episode %d: expected title '%s', got '%s'", i, config.Episodes[i].Title, episode.Title)
		}
		if episode.Description != config.Episodes[i].Description {
			t.Errorf("Episode %d: expected description '%s', got '%s'", i, config.Episodes[i].Description, episode.Description)
		}
		if episode.TranscriptFile != "ep"+string(rune(i+1+'0'))+".md" {
			t.Errorf("Episode %d: expected transcript file 'ep%d.md', got '%s'", i, i+1, episode.TranscriptFile)
		}
		// With nil client, summary should be the unavailable message
		if episode.Summary != "Summary unavailable (no OpenAI client)" {
			t.Errorf("Episode %d: expected summary 'Summary unavailable (no OpenAI client)', got '%s'", i, episode.Summary)
		}
	}
}

func TestGeneratePodcastTranscriptsWithPlaylistAndAudio(t *testing.T) {
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

	// Create test podcast config
	config := &PodcastConfig{
		Title:       "Test Podcast with Audio",
		Description: "A test podcast that generates audio and playlist.json",
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

	// Check that playlist.json was created
	playlistPath := filepath.Join(".reporadio", "test", "episodes", "playlist.json")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		t.Error("Expected playlist.json to be created")
	}

	// Read and verify the playlist.json content
	content, err := os.ReadFile(playlistPath)
	if err != nil {
		t.Fatalf("Failed to read playlist.json: %v", err)
	}

	// Verify JSON structure
	var playlistRoot PlaylistRoot
	err = json.Unmarshal(content, &playlistRoot)
	if err != nil {
		t.Fatalf("Failed to unmarshal playlist.json: %v", err)
	}

	// Verify that audio file is empty since generation failed with nil client
	if len(playlistRoot.Playlist.Episodes) > 0 {
		episode := playlistRoot.Playlist.Episodes[0]
		if episode.AudioFile != "" {
			t.Errorf("Expected empty audio file with nil client, got '%s'", episode.AudioFile)
		}
	}
}