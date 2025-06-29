package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateEpisodeWithScanner(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create test files
	if err := os.MkdirAll(filepath.Join(tmpDir, "src"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "src", "main.go"), []byte("package main\n\nfunc main() {}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test Project"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .gitignore
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte("*.log\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "debug.log"), []byte("debug info"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create episode with mixed includes
	episode := Episode{
		Title:       "Test Episode",
		Description: "Test description",
		Include: []string{
			filepath.Join(tmpDir, "README.md"),       // individual file
			filepath.Join(tmpDir, "src"),             // directory (should expand)
			filepath.Join(tmpDir, "debug.log"),       // gitignored file (should include)
			filepath.Join(tmpDir, "nonexistent.txt"), // missing file (should warn and skip)
		},
	}

	outputDir := filepath.Join(tmpDir, "output")

	// Generate episode transcript (without OpenAI client, so it creates placeholder)
	err := generateEpisodeTranscript(episode, 1, outputDir, nil, nil, "")
	if err != nil {
		t.Fatalf("generateEpisodeTranscript failed: %v", err)
	}

	// Check that transcript file was created
	transcriptPath := filepath.Join(outputDir, "ep1.md")
	if _, err := os.Stat(transcriptPath); err != nil {
		t.Fatalf("Transcript file not created: %v", err)
	}

	// Read transcript content
	content, err := os.ReadFile(transcriptPath)
	if err != nil {
		t.Fatalf("Failed to read transcript: %v", err)
	}

	transcriptContent := string(content)

	// Should contain episode title
	if !strings.Contains(transcriptContent, "Test Episode") {
		t.Error("Transcript should contain episode title")
	}
}
