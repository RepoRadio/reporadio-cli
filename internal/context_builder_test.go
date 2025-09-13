// Contract tests for ContextBuilder interface
// These tests define the expected behavior and MUST FAIL initially (RED phase)

package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ContextBuilder interface (contract definition)
// This will fail until implemented
type ContextBuilder interface {
	// BuildCreateContext combines command output and file content into context
	BuildCreateContext(config *CreateConfig, mode string) (*CreateContext, error)

	// GetDefaultConfig returns the default configuration when create.yml missing
	GetDefaultConfig() *CreateConfig
}

// NewContextBuilder creates a new context builder
// Implementation now exists in context_builder.go
func NewContextBuilderForTest() ContextBuilder {
	return NewContextBuilder()
}

func TestContextBuilder_BuildCreateContext(t *testing.T) {
	tests := []struct {
		name           string
		config         *CreateConfig
		mode           string
		setupFiles     map[string]string // path -> content
		expectError    bool
		validateResult func(t *testing.T, ctx *CreateContext)
	}{
		{
			name: "build context with config and commands",
			config: &CreateConfig{
				Include:  []string{"*.md", "*.txt"},
				Commands: []string{"echo 'test output'", "echo 'second command'"},
			},
			mode: ModeExplicit,
			setupFiles: map[string]string{
				"README.md": "# Test Project\nThis is a test",
				"notes.txt": "Some notes here",
				"ignore.go": "package main", // Should be excluded by include patterns
			},
			expectError: false,
			validateResult: func(t *testing.T, ctx *CreateContext) {
				if ctx.Mode != ModeExplicit {
					t.Errorf("expected mode %s, got %s", ModeExplicit, ctx.Mode)
				}
				if len(ctx.CommandOutput) != 2 {
					t.Errorf("expected 2 command outputs, got %d", len(ctx.CommandOutput))
				}
				if len(ctx.FileContent) != 2 {
					t.Errorf("expected 2 file contents, got %d", len(ctx.FileContent))
				}
				if ctx.FileContent["README.md"] != "# Test Project\nThis is a test" {
					t.Errorf("README.md content mismatch: got %s", ctx.FileContent["README.md"])
				}
				if !strings.Contains(ctx.CommandOutput["echo 'test output'"], "test output") {
					t.Errorf("command output missing expected content")
				}
			},
		},
		{
			name: "build context with default config",
			config: &CreateConfig{
				Include:  []string{"."},
				Commands: []string{},
			},
			mode: ModeDefault,
			setupFiles: map[string]string{
				"README.md": "Default behavior",
				"main.go":   "package main",
			},
			expectError: false,
			validateResult: func(t *testing.T, ctx *CreateContext) {
				if ctx.Mode != ModeDefault {
					t.Errorf("expected mode %s, got %s", ModeDefault, ctx.Mode)
				}
				if len(ctx.CommandOutput) != 0 {
					t.Errorf("expected no command output, got %d", len(ctx.CommandOutput))
				}
				// Should include all files (not filtered by pattern since "." includes all)
				if len(ctx.FileContent) < 2 {
					t.Errorf("expected at least 2 files, got %d", len(ctx.FileContent))
				}
			},
		},
		{
			name: "command execution fails",
			config: &CreateConfig{
				Include:  []string{"*.md"},
				Commands: []string{"nonexistent_command_12345"},
			},
			mode:        ModeExplicit,
			setupFiles:  map[string]string{"README.md": "test"},
			expectError: true,
		},
		{
			name: "empty config produces empty context",
			config: &CreateConfig{
				Include:  []string{},
				Commands: []string{},
			},
			mode:        ModeExplicit,
			setupFiles:  map[string]string{"README.md": "test"},
			expectError: false,
			validateResult: func(t *testing.T, ctx *CreateContext) {
				if len(ctx.CommandOutput) != 0 {
					t.Errorf("expected no command output, got %d", len(ctx.CommandOutput))
				}
				if len(ctx.FileContent) != 0 {
					t.Errorf("expected no file content, got %d", len(ctx.FileContent))
				}
			},
		},
		{
			name: "commands generate files that are then included",
			config: &CreateConfig{
				Include:  []string{"*.md", "generated.txt"},
				Commands: []string{"echo 'Generated content' > generated.txt"},
			},
			mode: ModeExplicit,
			setupFiles: map[string]string{
				"README.md": "Original content",
			},
			expectError: false,
			validateResult: func(t *testing.T, ctx *CreateContext) {
				// Should have both original file and generated file
				if len(ctx.FileContent) != 2 {
					t.Errorf("expected 2 files, got %d", len(ctx.FileContent))
				}
				if ctx.FileContent["README.md"] != "Original content" {
					t.Errorf("README.md content incorrect")
				}
				if !strings.Contains(ctx.FileContent["generated.txt"], "Generated content") {
					t.Errorf("generated.txt content missing: %s", ctx.FileContent["generated.txt"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test directory
			tempDir := t.TempDir()

			for path, content := range tt.setupFiles {
				filePath := filepath.Join(tempDir, path)
				dir := filepath.Dir(filePath)
				os.MkdirAll(dir, 0755)
				os.WriteFile(filePath, []byte(content), 0644)
			}

			// Change to test directory
			oldDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(oldDir)

			// Test the contract (WILL FAIL - no implementation yet)
			builder := NewContextBuilderForTest()
			ctx, err := builder.BuildCreateContext(tt.config, tt.mode)

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

			// Verify config is preserved in context
			if ctx.Config != tt.config {
				t.Error("config not preserved in context")
			}

			if tt.validateResult != nil {
				tt.validateResult(t, ctx)
			}
		})
	}
}

func TestContextBuilder_GetDefaultConfig(t *testing.T) {
	builder := NewContextBuilder() // This will panic
	defaultConfig := builder.GetDefaultConfig()

	if len(defaultConfig.Include) != 1 || defaultConfig.Include[0] != "." {
		t.Errorf("expected default include ['.'], got %v", defaultConfig.Include)
	}

	if len(defaultConfig.Commands) != 0 {
		t.Errorf("expected empty commands, got %v", defaultConfig.Commands)
	}
}

func TestContextBuilder_Integration(t *testing.T) {
	t.Run("full workflow with commands and includes", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create initial files
		setupFiles := map[string]string{
			"README.md":   "# Project\nDescription here",
			"main.go":     "package main\n\nfunc main() {}",
			"docs/api.md": "# API Documentation",
		}

		for path, content := range setupFiles {
			filePath := filepath.Join(tempDir, path)
			dir := filepath.Dir(filePath)
			os.MkdirAll(dir, 0755)
			os.WriteFile(filePath, []byte(content), 0644)
		}

		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		// Config that generates additional content and includes specific patterns
		config := &CreateConfig{
			Include: []string{"*.md", "structure.txt"},
			Commands: []string{
				"echo 'Project structure:' > structure.txt",
				"echo 'Files: README.md, main.go, docs/api.md' >> structure.txt",
			},
		}

		builder := NewContextBuilderForTest()
		ctx, err := builder.BuildCreateContext(config, ModeExplicit)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Should have command outputs
		if len(ctx.CommandOutput) != 2 {
			t.Errorf("expected 2 command outputs, got %d", len(ctx.CommandOutput))
		}

		// Should have file contents (README.md, docs/api.md, structure.txt)
		expectedFiles := 3
		if len(ctx.FileContent) != expectedFiles {
			actualFiles := make([]string, 0, len(ctx.FileContent))
			for path := range ctx.FileContent {
				actualFiles = append(actualFiles, path)
			}
			t.Errorf("expected %d files, got %d: %v", expectedFiles, len(ctx.FileContent), actualFiles)
		}

		// Verify specific content
		if !strings.Contains(ctx.FileContent["README.md"], "# Project") {
			t.Errorf("README.md content incorrect: %s", ctx.FileContent["README.md"])
		}

		if !strings.Contains(ctx.FileContent["structure.txt"], "Project structure:") {
			t.Errorf("structure.txt content incorrect: %s", ctx.FileContent["structure.txt"])
		}
	})

	t.Run("handles missing include files gracefully", func(t *testing.T) {
		tempDir := t.TempDir()

		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		// Config that includes files that don't exist
		config := &CreateConfig{
			Include:  []string{"nonexistent.md", "missing/*.txt"},
			Commands: []string{},
		}

		builder := NewContextBuilderForTest()
		ctx, err := builder.BuildCreateContext(config, ModeExplicit)

		// Should not error for missing files (graceful handling)
		if err != nil {
			t.Errorf("should handle missing files gracefully, got error: %v", err)
		}

		// Should have empty file content
		if len(ctx.FileContent) != 0 {
			t.Errorf("expected empty file content for missing files, got %d", len(ctx.FileContent))
		}
	})
}

func TestContextBuilder_ErrorHandling(t *testing.T) {
	t.Run("invalid config causes error", func(t *testing.T) {
		builder := NewContextBuilderForTest()

		// Pass nil config
		_, err := builder.BuildCreateContext(nil, ModeDefault)

		if err == nil {
			t.Error("expected error for nil config")
		}
	})

	t.Run("invalid mode handled gracefully", func(t *testing.T) {
		config := &CreateConfig{
			Include:  []string{},
			Commands: []string{},
		}

		builder := NewContextBuilderForTest()

		// Use invalid mode
		ctx, err := builder.BuildCreateContext(config, "invalid_mode")

		// Should either error or handle gracefully
		if err == nil && ctx.Mode != "invalid_mode" {
			t.Errorf("mode should be preserved or cause error, got mode: %s", ctx.Mode)
		}
	})
}
