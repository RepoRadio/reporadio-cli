// Contract tests for FileContentLoader interface
// These tests define the expected behavior and MUST FAIL initially (RED phase)

package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// FileContentLoader interface (contract definition)
// This will fail until implemented
type FileContentLoader interface {
	// LoadFileContent loads content from files matching scan result
	LoadFileContent(scanResult *ScanResult) (map[string]string, error)

	// LoadSingleFile loads content from a single file
	LoadSingleFile(path string) (string, error)
}

// NewFileContentLoader creates a new file content loader
// Implementation now exists in file_content_loader.go
func NewFileContentLoaderForTest() FileContentLoader {
	return NewFileContentLoader()
}

func TestFileContentLoader_LoadFileContent(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string
		scanResult     *ScanResult
		expectError    bool
		validateResult func(t *testing.T, content map[string]string)
	}{
		{
			name: "load multiple files",
			setupFiles: map[string]string{
				"file1.txt": "Content 1",
				"file2.txt": "Content 2",
				"file3.md":  "# Markdown Content",
			},
			scanResult: &ScanResult{
				Files: []FileInfo{
					{Path: "file1.txt", Extension: ".txt", Size: 9},
					{Path: "file2.txt", Extension: ".txt", Size: 9},
					{Path: "file3.md", Extension: ".md", Size: 19},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, content map[string]string) {
				if len(content) != 3 {
					t.Errorf("expected 3 files, got %d", len(content))
				}
				if content["file1.txt"] != "Content 1" {
					t.Errorf("file1.txt content mismatch: got %s", content["file1.txt"])
				}
				if content["file2.txt"] != "Content 2" {
					t.Errorf("file2.txt content mismatch: got %s", content["file2.txt"])
				}
				if content["file3.md"] != "# Markdown Content" {
					t.Errorf("file3.md content mismatch: got %s", content["file3.md"])
				}
			},
		},
		{
			name: "handle missing files gracefully",
			setupFiles: map[string]string{
				"exists.txt": "Content exists",
			},
			scanResult: &ScanResult{
				Files: []FileInfo{
					{Path: "exists.txt", Extension: ".txt", Size: 14},
					{Path: "missing.txt", Extension: ".txt", Size: 0},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, content map[string]string) {
				// Should handle missing files gracefully - only include existing files
				if content["exists.txt"] != "Content exists" {
					t.Errorf("existing file content incorrect: got %s", content["exists.txt"])
				}
				// Missing file should not be in results or should have empty content
				if missingContent, exists := content["missing.txt"]; exists && missingContent != "" {
					t.Errorf("missing file should have empty content, got: %s", missingContent)
				}
			},
		},
		{
			name: "empty scan result",
			setupFiles: map[string]string{
				"ignored.txt": "This file is ignored",
			},
			scanResult: &ScanResult{
				Files: []FileInfo{},
			},
			expectError: false,
			validateResult: func(t *testing.T, content map[string]string) {
				if len(content) != 0 {
					t.Errorf("expected empty content map, got %d items", len(content))
				}
			},
		},
		{
			name: "handle binary files appropriately",
			setupFiles: map[string]string{
				"text.txt":   "Regular text content",
				"binary.bin": "\x00\x01\x02\x03", // Binary content
			},
			scanResult: &ScanResult{
				Files: []FileInfo{
					{Path: "text.txt", Extension: ".txt", Size: 20},
					{Path: "binary.bin", Extension: ".bin", Size: 4},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, content map[string]string) {
				if content["text.txt"] != "Regular text content" {
					t.Errorf("text file content incorrect")
				}
				// Binary content should be handled (converted to string or skipped)
				if _, exists := content["binary.bin"]; exists {
					t.Logf("Binary file included: %q", content["binary.bin"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Setup files with absolute paths for ScanResult
			for path, content := range tt.setupFiles {
				filePath := filepath.Join(tempDir, path)
				os.WriteFile(filePath, []byte(content), 0644)
			}

			// Convert relative paths in ScanResult to absolute paths
			for i := range tt.scanResult.Files {
				if !filepath.IsAbs(tt.scanResult.Files[i].Path) {
					tt.scanResult.Files[i].Path = filepath.Join(tempDir, tt.scanResult.Files[i].Path)
				}
			}

			oldDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(oldDir)

			loader := NewFileContentLoaderForTest()
			content, err := loader.LoadFileContent(tt.scanResult)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validateResult != nil {
				// Convert absolute paths back to relative for validation
				relativeContent := make(map[string]string)
				for path, content := range content {
					if relPath, err := filepath.Rel(tempDir, path); err == nil {
						relativeContent[relPath] = content
					} else {
						relativeContent[path] = content
					}
				}
				tt.validateResult(t, relativeContent)
			}
		})
	}
}

func TestFileContentLoader_LoadSingleFile(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		fileContent     string
		expectError     bool
		expectedContent string
	}{
		{
			name:            "load regular text file",
			fileName:        "test.txt",
			fileContent:     "Hello, World!",
			expectError:     false,
			expectedContent: "Hello, World!",
		},
		{
			name:            "load markdown file",
			fileName:        "README.md",
			fileContent:     "# Title\n\nContent here",
			expectError:     false,
			expectedContent: "# Title\n\nContent here",
		},
		{
			name:            "load empty file",
			fileName:        "empty.txt",
			fileContent:     "",
			expectError:     false,
			expectedContent: "",
		},
		{
			name:        "file does not exist",
			fileName:    "nonexistent.txt",
			fileContent: "", // Not created
			expectError: true,
		},
		{
			name:            "load file with unicode content",
			fileName:        "unicode.txt",
			fileContent:     "Hello 世界 🌍",
			expectError:     false,
			expectedContent: "Hello 世界 🌍",
		},
		{
			name:            "load large text file",
			fileName:        "large.txt",
			fileContent:     strings.Repeat("Line of text\n", 1000),
			expectError:     false,
			expectedContent: strings.Repeat("Line of text\n", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, tt.fileName)

			// Only create file if content is provided
			if tt.fileContent != "" || !tt.expectError {
				os.WriteFile(filePath, []byte(tt.fileContent), 0644)
			}

			loader := NewFileContentLoaderForTest()
			content, err := loader.LoadSingleFile(filePath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if content != tt.expectedContent {
				t.Errorf("content mismatch:\nexpected: %q\ngot: %q", tt.expectedContent, content)
			}
		})
	}
}

func TestFileContentLoader_Integration(t *testing.T) {
	t.Run("load content from realistic scan result", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a realistic project structure
		files := map[string]string{
			"README.md":         "# Test Project\n\nA sample project for testing.",
			"main.go":           "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			"config.json":       "{\n\t\"name\": \"test-project\",\n\t\"version\": \"1.0.0\"\n}",
			"docs/api.md":       "# API Documentation\n\n## Endpoints\n\n- GET /health",
			"internal/utils.go": "package internal\n\n// Helper functions",
		}

		for path, content := range files {
			filePath := filepath.Join(tempDir, path)
			dir := filepath.Dir(filePath)
			os.MkdirAll(dir, 0755)
			os.WriteFile(filePath, []byte(content), 0644)
		}

		// Create scan result with absolute paths
		scanResult := &ScanResult{
			Files: []FileInfo{
				{Path: filepath.Join(tempDir, "README.md"), Extension: ".md", Size: 44},
				{Path: filepath.Join(tempDir, "main.go"), Extension: ".go", Size: 85},
				{Path: filepath.Join(tempDir, "docs/api.md"), Extension: ".md", Size: 50},
			},
		}

		loader := NewFileContentLoaderForTest()
		content, err := loader.LoadFileContent(scanResult)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if len(content) != 3 {
			t.Errorf("expected 3 files loaded, got %d", len(content))
		}

		// Verify specific content
		readmePath := filepath.Join(tempDir, "README.md")
		if !strings.Contains(content[readmePath], "# Test Project") {
			t.Errorf("README.md content incorrect: %s", content[readmePath])
		}

		mainPath := filepath.Join(tempDir, "main.go")
		if !strings.Contains(content[mainPath], "package main") {
			t.Errorf("main.go content incorrect: %s", content[mainPath])
		}

		apiPath := filepath.Join(tempDir, "docs/api.md")
		if !strings.Contains(content[apiPath], "# API Documentation") {
			t.Errorf("docs/api.md content incorrect: %s", content[apiPath])
		}
	})

	t.Run("handle mixed success and failure", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create only some of the files
		os.WriteFile(filepath.Join(tempDir, "exists1.txt"), []byte("Content 1"), 0644)
		os.WriteFile(filepath.Join(tempDir, "exists2.txt"), []byte("Content 2"), 0644)
		// missing.txt is not created

		scanResult := &ScanResult{
			Files: []FileInfo{
				{Path: filepath.Join(tempDir, "exists1.txt"), Extension: ".txt", Size: 9},
				{Path: filepath.Join(tempDir, "missing.txt"), Extension: ".txt", Size: 0},
				{Path: filepath.Join(tempDir, "exists2.txt"), Extension: ".txt", Size: 9},
			},
		}

		loader := NewFileContentLoaderForTest()
		content, err := loader.LoadFileContent(scanResult)

		// Should handle missing files gracefully, not fail entirely
		if err != nil {
			t.Errorf("should handle missing files gracefully, got error: %v", err)
		}

		// Should have content for existing files
		exists1Path := filepath.Join(tempDir, "exists1.txt")
		exists2Path := filepath.Join(tempDir, "exists2.txt")

		if content[exists1Path] != "Content 1" {
			t.Errorf("exists1.txt content incorrect")
		}
		if content[exists2Path] != "Content 2" {
			t.Errorf("exists2.txt content incorrect")
		}
	})
}

func TestFileContentLoader_EdgeCases(t *testing.T) {
	t.Run("nil scan result", func(t *testing.T) {
		loader := NewFileContentLoaderForTest()
		content, err := loader.LoadFileContent(nil)

		if err == nil {
			t.Error("expected error for nil scan result")
		}

		if content != nil {
			t.Error("expected nil content for nil scan result")
		}
	})

	t.Run("empty file path", func(t *testing.T) {
		loader := NewFileContentLoaderForTest()
		content, err := loader.LoadSingleFile("")

		if err == nil {
			t.Error("expected error for empty file path")
		}

		if content != "" {
			t.Error("expected empty content for empty file path")
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		tempDir := t.TempDir()
		dirPath := filepath.Join(tempDir, "subdir")
		os.MkdirAll(dirPath, 0755)

		loader := NewFileContentLoaderForTest()
		content, err := loader.LoadSingleFile(dirPath)

		// Should error when trying to read directory as file
		if err == nil {
			t.Error("expected error when reading directory as file")
		}

		if content != "" {
			t.Error("expected empty content when reading directory as file")
		}
	})
}
