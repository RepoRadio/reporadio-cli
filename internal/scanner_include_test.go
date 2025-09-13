package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveIncludePaths(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create files and directories
	os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "src", "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "src", "helper.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)

	// Create .gitignore to ignore some files
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte("*.log\ntemp/\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "debug.log"), []byte("log content"), 0644)

	scanner := NewScanner()

	// Test with mix of files, directories, and non-existent paths
	includes := []string{
		filepath.Join(tmpDir, "README.md"),       // existing file
		filepath.Join(tmpDir, "src"),             // existing directory
		filepath.Join(tmpDir, "nonexistent.txt"), // non-existent file
		filepath.Join(tmpDir, "missing-dir"),     // non-existent directory
		filepath.Join(tmpDir, "debug.log"),       // gitignored file (should still include)
	}

	resolved, err := scanner.ResolveIncludePaths(includes)
	if err != nil {
		t.Fatalf("ResolveIncludePaths failed: %v", err)
	}

	// Should include: README.md, main.go, helper.go, debug.log
	// Should skip: nonexistent.txt, missing-dir (with warnings)
	expectedCount := 4
	if len(resolved) != expectedCount {
		t.Errorf("Expected %d resolved paths, got %d: %v", expectedCount, len(resolved), resolved)
	}

	// Check that all expected files are present
	expectedFiles := map[string]bool{
		filepath.Join(tmpDir, "README.md"):        true,
		filepath.Join(tmpDir, "src", "main.go"):   true,
		filepath.Join(tmpDir, "src", "helper.go"): true,
		filepath.Join(tmpDir, "debug.log"):        true,
	}

	for _, resolvedPath := range resolved {
		if !expectedFiles[resolvedPath] {
			t.Errorf("Unexpected resolved path: %s", resolvedPath)
		}
		delete(expectedFiles, resolvedPath)
	}

	if len(expectedFiles) > 0 {
		t.Errorf("Missing expected files: %v", expectedFiles)
	}
}
