package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateEpisodeAudio(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test transcript file
	transcriptPath := filepath.Join(tempDir, "test_transcript.md")
	transcriptContent := "# Test Episode\n\nThis is a short test transcript for audio generation."

	err := os.WriteFile(transcriptPath, []byte(transcriptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test transcript: %v", err)
	}

	audioPath := filepath.Join(tempDir, "test_audio.mp3")

	// Test with nil client (should handle gracefully)
	err = generateEpisodeAudio(transcriptPath, audioPath, nil)
	if err == nil {
		t.Error("Expected error with nil client")
	}

	// Check that no audio file was created
	if _, err := os.Stat(audioPath); err == nil {
		t.Error("Did not expect audio file to be created with nil client")
	}
}

func TestConvertTextToSpeech(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	audioPath := filepath.Join(tempDir, "test_speech.mp3")

	// Test with nil client (should handle gracefully)
	err := convertTextToSpeech(nil, "Hello, world!", audioPath)
	if err == nil {
		t.Error("Expected error with nil client")
	}

	// Check that no audio file was created
	if _, err := os.Stat(audioPath); err == nil {
		t.Error("Did not expect audio file to be created with nil client")
	}
}
