// Contract tests for ConfigLoader interface
// These tests define the expected behavior and MUST FAIL initially (RED phase)

package internal

import (
	"os"
	"path/filepath"
	"testing"
)

// ConfigLoader interface (contract definition)
// This will fail until implemented
type ConfigLoader interface {
	// LoadCreateConfig loads configuration from .reporadio/create.yml
	// Returns default config if file doesn't exist (backward compatibility)
	LoadCreateConfig() (*CreateConfig, string, error) // config, mode, error

	// ValidateConfig validates the loaded configuration
	ValidateConfig(config *CreateConfig) error

	// HasCreateConfig checks if .reporadio/create.yml exists
	HasCreateConfig() bool
}

// NewConfigLoader creates a new config loader
// Implementation now exists in config.go
func NewConfigLoaderForTest() ConfigLoader {
	return NewConfigLoader()
}

func TestConfigLoader_LoadCreateConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectError    bool
		expectedMode   string
		expectedConfig *CreateConfig
	}{
		{
			name: "valid config with include and commands",
			configContent: `include:
  - "README.md"
  - "docs/*.md"
commands:
  - "go doc -all > api.md"
  - "tree > structure.txt"`,
			expectError:  false,
			expectedMode: ModeExplicit,
			expectedConfig: &CreateConfig{
				Include:  []string{"README.md", "docs/*.md"},
				Commands: []string{"go doc -all > api.md", "tree > structure.txt"},
			},
		},
		{
			name: "config with only include",
			configContent: `include:
  - "*.go"
  - "internal/**/*.go"`,
			expectError:  false,
			expectedMode: ModeExplicit,
			expectedConfig: &CreateConfig{
				Include:  []string{"*.go", "internal/**/*.go"},
				Commands: []string{},
			},
		},
		{
			name: "config with only commands",
			configContent: `commands:
  - "make docs"
  - "npm run build"`,
			expectError:  false,
			expectedMode: ModeExplicit,
			expectedConfig: &CreateConfig{
				Include:  []string{},
				Commands: []string{"make docs", "npm run build"},
			},
		},
		{
			name:          "empty config file",
			configContent: "",
			expectError:   false,
			expectedMode:  ModeExplicit,
			expectedConfig: &CreateConfig{
				Include:  []string{},
				Commands: []string{},
			},
		},
		{
			name:          "invalid YAML",
			configContent: "invalid: yaml: content [",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test directory
			tempDir := t.TempDir()
			reporadioDir := filepath.Join(tempDir, ".reporadio")
			os.MkdirAll(reporadioDir, 0755)
			configPath := filepath.Join(reporadioDir, "create.yml")

			// Write test config
			os.WriteFile(configPath, []byte(tt.configContent), 0644)

			// Change to test directory
			oldDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(oldDir)

			// Test the contract
			loader := NewConfigLoaderForTest()
			config, mode, err := loader.LoadCreateConfig()

			// Assertions
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

			if mode != tt.expectedMode {
				t.Errorf("mode mismatch: got %s, want %s", mode, tt.expectedMode)
			}

			// Compare configurations
			if !equalStringSlices(config.Include, tt.expectedConfig.Include) {
				t.Errorf("Include mismatch: got %v, want %v", config.Include, tt.expectedConfig.Include)
			}

			if !equalStringSlices(config.Commands, tt.expectedConfig.Commands) {
				t.Errorf("Commands mismatch: got %v, want %v", config.Commands, tt.expectedConfig.Commands)
			}
		})
	}
}

func TestConfigLoader_LoadCreateConfig_NoFile(t *testing.T) {
	// Test default behavior when no create.yml exists
	tempDir := t.TempDir()

	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	loader := NewConfigLoaderForTest()
	config, mode, err := loader.LoadCreateConfig()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if mode != ModeDefault {
		t.Errorf("expected mode %s, got %s", ModeDefault, mode)
	}

	// Should return default values
	expectedInclude := []string{"."}
	if !equalStringSlices(config.Include, expectedInclude) {
		t.Errorf("Include mismatch: got %v, want %v", config.Include, expectedInclude)
	}

	expectedCommands := []string{}
	if !equalStringSlices(config.Commands, expectedCommands) {
		t.Errorf("Commands mismatch: got %v, want %v", config.Commands, expectedCommands)
	}
}

func TestConfigLoader_HasCreateConfig(t *testing.T) {
	tests := []struct {
		name        string
		fileExists  bool
		expectedHas bool
	}{
		{
			name:        "config file exists",
			fileExists:  true,
			expectedHas: true,
		},
		{
			name:        "config file does not exist",
			fileExists:  false,
			expectedHas: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.fileExists {
				reporadioDir := filepath.Join(tempDir, ".reporadio")
				os.MkdirAll(reporadioDir, 0755)
				configPath := filepath.Join(reporadioDir, "create.yml")
				os.WriteFile(configPath, []byte("include: []"), 0644)
			}

			oldDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(oldDir)

			loader := NewConfigLoaderForTest()
			has := loader.HasCreateConfig()

			if has != tt.expectedHas {
				t.Errorf("HasCreateConfig() = %v, want %v", has, tt.expectedHas)
			}
		})
	}
}

func TestConfigLoader_ValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *CreateConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &CreateConfig{
				Include:  []string{"*.go", "docs/*.md"},
				Commands: []string{"go fmt", "make test"},
			},
			expectError: false,
		},
		{
			name: "invalid glob pattern",
			config: &CreateConfig{
				Include: []string{"[invalid"},
			},
			expectError: true,
		},
		{
			name: "empty command string",
			config: &CreateConfig{
				Commands: []string{"", "make test"},
			},
			expectError: true,
		},
		{
			name: "valid empty config",
			config: &CreateConfig{
				Include:  []string{},
				Commands: []string{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewConfigLoaderForTest()
			err := loader.ValidateConfig(tt.config)

			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

// Helper function for comparing string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
