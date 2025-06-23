package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanRepository_WithReporadioIgnore(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"README.md",
		"main.go",
		"test.log",    // will be ignored by .gitignore
		"docs/api.md", // will be ignored by .reporadioignore
		"src/helper.go",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)
		if dir != tmpDir {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
		err := os.WriteFile(fullPath, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Create .gitignore
	gitignoreContent := "*.log\n"
	err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Create .reporadioignore
	reporadioIgnoreContent := "docs/\n"
	err = os.WriteFile(filepath.Join(tmpDir, ".reporadioignore"), []byte(reporadioIgnoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .reporadioignore: %v", err)
	}

	scanner := NewScanner()
	result, err := scanner.ScanRepository(tmpDir)

	if err != nil {
		t.Fatalf("ScanRepository failed: %v", err)
	}

	// Should include README.md, main.go, src/helper.go
	// Should exclude test.log (gitignore) and docs/api.md (reporadioignore)
	expectedFiles := []string{"README.md", "main.go", "src/helper.go"}

	if len(result.Files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(result.Files))
	}

	foundFiles := make(map[string]bool)
	for _, file := range result.Files {
		rel, _ := filepath.Rel(tmpDir, file.Path)
		foundFiles[rel] = true
	}

	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected file %s not found", expected)
		}
	}

	// Verify excluded files are not present
	excludedFiles := []string{"test.log", "docs/api.md"}
	for _, excluded := range excludedFiles {
		if foundFiles[excluded] {
			t.Errorf("File %s should have been ignored", excluded)
		}
	}
}
