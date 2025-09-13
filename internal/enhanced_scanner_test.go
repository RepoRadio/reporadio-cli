// Contract tests for EnhancedScanner interface
// These tests define the expected behavior and MUST FAIL initially (RED phase)

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// EnhancedScanner interface (contract definition)
// This will fail until implemented - extends existing Scanner
type EnhancedScanner interface {
	// ScanWithIncludes scans repository using include patterns
	// Empty patterns slice falls back to current behavior
	ScanWithIncludes(rootPath string, includes []string) (*ScanResult, error)

	// ScanRepository maintains existing contract for backward compatibility
	ScanRepository(rootPath string) (*ScanResult, error)
}

// NewEnhancedScanner creates a new enhanced scanner
// Implementation uses existing Scanner which now has ScanWithIncludes
func NewEnhancedScannerForTest() EnhancedScanner {
	return NewScanner()
}

func TestEnhancedScanner_ScanWithIncludes(t *testing.T) {
	tests := []struct {
		name          string
		includes      []string
		setupFiles    []string
		expectedFiles []string
	}{
		{
			name:     "include specific files",
			includes: []string{"README.md", "main.go"},
			setupFiles: []string{
				"README.md",
				"main.go",
				"test.go",
				"docs/api.md",
			},
			expectedFiles: []string{"README.md", "main.go"},
		},
		{
			name:     "include with glob patterns",
			includes: []string{"*.go", "docs/*.md"},
			setupFiles: []string{
				"main.go",
				"test.go",
				"README.md",
				"docs/api.md",
				"docs/guide.md",
				"docs/internal/private.md",
			},
			expectedFiles: []string{"main.go", "test.go", "docs/api.md", "docs/guide.md"},
		},
		{
			name:     "include directory recursively",
			includes: []string{"src/**/*.go"},
			setupFiles: []string{
				"main.go",
				"src/app.go",
				"src/internal/helper.go",
				"src/internal/deep/util.go",
				"test/test.go",
			},
			expectedFiles: []string{"src/app.go", "src/internal/helper.go", "src/internal/deep/util.go"},
		},
		{
			name:          "empty includes produces no results",
			includes:      []string{},
			setupFiles:    []string{"main.go", "README.md", "test.go"},
			expectedFiles: []string{}, // No files - explicit empty means include nothing
		},
		{
			name:          "no matching files",
			includes:      []string{"*.nonexistent"},
			setupFiles:    []string{"main.go", "README.md"},
			expectedFiles: []string{}, // No matches
		},
		{
			name:     "include current directory pattern",
			includes: []string{"."},
			setupFiles: []string{
				"main.go",
				"README.md",
				"docs/api.md",
				".gitignore",
			},
			expectedFiles: []string{"main.go", "README.md", "docs/api.md"}, // Should not include .gitignore (hidden)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test directory structure
			tempDir := t.TempDir()

			for _, file := range tt.setupFiles {
				filePath := filepath.Join(tempDir, file)
				dir := filepath.Dir(filePath)
				os.MkdirAll(dir, 0755)
				os.WriteFile(filePath, []byte("test content"), 0644)
			}

			// Test the contract (WILL FAIL - no implementation yet)
			scanner := NewEnhancedScannerForTest()
			result, err := scanner.ScanWithIncludes(tempDir, tt.includes)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify expected files are included
			foundFiles := make(map[string]bool)
			for _, fileInfo := range result.Files {
				// Convert absolute path to relative for comparison
				relPath, _ := filepath.Rel(tempDir, fileInfo.Path)
				foundFiles[relPath] = true
			}

			// Check all expected files are present
			for _, expectedFile := range tt.expectedFiles {
				if !foundFiles[expectedFile] {
					t.Errorf("expected file %s not found in scan result", expectedFile)
				}
			}

			// Check no unexpected files are included
			if len(result.Files) != len(tt.expectedFiles) {
				actualFiles := make([]string, 0, len(result.Files))
				for _, fileInfo := range result.Files {
					relPath, _ := filepath.Rel(tempDir, fileInfo.Path)
					actualFiles = append(actualFiles, relPath)
				}
				t.Errorf("expected %d files, got %d files. Expected: %v, Got: %v",
					len(tt.expectedFiles), len(result.Files), tt.expectedFiles, actualFiles)
			}
		})
	}
}

func TestEnhancedScanner_BackwardCompatibility(t *testing.T) {
	// Test that ScanRepository still works as before
	t.Run("ScanRepository maintains existing behavior", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create test files
		testFiles := []string{"README.md", "main.go", "test.go"}
		for _, file := range testFiles {
			filePath := filepath.Join(tempDir, file)
			os.WriteFile(filePath, []byte("content"), 0644)
		}

		scanner := NewEnhancedScannerForTest()
		result, err := scanner.ScanRepository(tempDir)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if len(result.Files) != len(testFiles) {
			t.Errorf("expected %d files, got %d", len(testFiles), len(result.Files))
		}

		// Verify all test files are found
		foundFiles := make(map[string]bool)
		for _, fileInfo := range result.Files {
			relPath, _ := filepath.Rel(tempDir, fileInfo.Path)
			foundFiles[relPath] = true
		}

		for _, expectedFile := range testFiles {
			if !foundFiles[expectedFile] {
				t.Errorf("expected file %s not found in scan result", expectedFile)
			}
		}
	})
}

func TestEnhancedScanner_GitignoreIntegration(t *testing.T) {
	// Test that include patterns work with existing gitignore functionality
	t.Run("includes respect gitignore rules", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create gitignore file
		gitignorePath := filepath.Join(tempDir, ".gitignore")
		os.WriteFile(gitignorePath, []byte("*.tmp\nbuild/\n"), 0644)

		// Create test files
		testFiles := []string{
			"main.go",
			"temp.tmp",     // Should be ignored by gitignore
			"build/out.go", // Should be ignored by gitignore
		}
		for _, file := range testFiles {
			filePath := filepath.Join(tempDir, file)
			dir := filepath.Dir(filePath)
			os.MkdirAll(dir, 0755)
			os.WriteFile(filePath, []byte("content"), 0644)
		}

		scanner := NewEnhancedScannerForTest()
		result, err := scanner.ScanWithIncludes(tempDir, []string{"*.go", "*.tmp"})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Should only find main.go (*.tmp and build/*.go ignored by gitignore)
		if len(result.Files) != 1 {
			actualFiles := make([]string, 0, len(result.Files))
			for _, fileInfo := range result.Files {
				relPath, _ := filepath.Rel(tempDir, fileInfo.Path)
				actualFiles = append(actualFiles, relPath)
			}
			t.Errorf("expected 1 file, got %d files: %v", len(result.Files), actualFiles)
		}

		if len(result.Files) > 0 {
			relPath, _ := filepath.Rel(tempDir, result.Files[0].Path)
			if relPath != "main.go" {
				t.Errorf("expected main.go, got %s", relPath)
			}
		}
	})

	t.Run("includes work with reporadioignore", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create reporadioignore file
		reporadioignorePath := filepath.Join(tempDir, ".reporadioignore")
		os.WriteFile(reporadioignorePath, []byte("internal/\n*.log\n"), 0644)

		// Create test files
		testFiles := []string{
			"main.go",
			"app.log",            // Should be ignored by reporadioignore
			"internal/config.go", // Should be ignored by reporadioignore
			"docs/readme.md",
		}
		for _, file := range testFiles {
			filePath := filepath.Join(tempDir, file)
			dir := filepath.Dir(filePath)
			os.MkdirAll(dir, 0755)
			os.WriteFile(filePath, []byte("content"), 0644)
		}

		scanner := NewEnhancedScannerForTest()
		result, err := scanner.ScanWithIncludes(tempDir, []string{"*.go", "*.md", "*.log"})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Should only find main.go and docs/readme.md (internal/ and *.log ignored)
		expectedFiles := []string{"main.go", "docs/readme.md"}
		if len(result.Files) != len(expectedFiles) {
			actualFiles := make([]string, 0, len(result.Files))
			for _, fileInfo := range result.Files {
				relPath, _ := filepath.Rel(tempDir, fileInfo.Path)
				actualFiles = append(actualFiles, relPath)
			}
			t.Errorf("expected %d files %v, got %d files: %v",
				len(expectedFiles), expectedFiles, len(result.Files), actualFiles)
		}
	})
}

func TestEnhancedScanner_Performance(t *testing.T) {
	t.Run("include patterns optimize scanning", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create many files but only include a few
		for i := 0; i < 100; i++ {
			fileName := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
			os.WriteFile(fileName, []byte("content"), 0644)
		}

		// Create one .go file
		goFile := filepath.Join(tempDir, "main.go")
		os.WriteFile(goFile, []byte("package main"), 0644)

		scanner := NewEnhancedScannerForTest()

		// Only include .go files - should be more efficient than scanning all files
		result, err := scanner.ScanWithIncludes(tempDir, []string{"*.go"})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Should only find the one .go file
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d", len(result.Files))
		}

		if len(result.Files) == 1 {
			relPath, _ := filepath.Rel(tempDir, result.Files[0].Path)
			if relPath != "main.go" {
				t.Errorf("expected main.go, got %s", relPath)
			}
		}
	})
}
