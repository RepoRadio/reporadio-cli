package internal

import (
	"testing"
)

// TestCreateConfig_BasicStructure tests that CreateConfig struct is properly defined
func TestCreateConfig_BasicStructure(t *testing.T) {
	config := &CreateConfig{
		Include:  []string{"*.md"},
		Commands: []string{"echo test"},
	}

	if len(config.Include) != 1 {
		t.Errorf("Expected 1 include pattern, got %d", len(config.Include))
	}

	if len(config.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(config.Commands))
	}
}

// TestCreateContext_BasicStructure tests that CreateContext struct is properly defined
func TestCreateContext_BasicStructure(t *testing.T) {
	ctx := &CreateContext{
		CommandOutput: make(map[string]string),
		FileContent:   make(map[string]string),
		Config:        &CreateConfig{},
		Mode:          ModeDefault,
	}

	if ctx.CommandOutput == nil {
		t.Error("CommandOutput should be initialized")
	}

	if ctx.FileContent == nil {
		t.Error("FileContent should be initialized")
	}

	if ctx.Mode != ModeDefault {
		t.Errorf("Expected mode %s, got %s", ModeDefault, ctx.Mode)
	}
}

// TestConstants tests that configuration constants are properly defined
func TestConstants(t *testing.T) {
	if CreateConfigPath != ".reporadio/create.yml" {
		t.Errorf("Expected CreateConfigPath to be '.reporadio/create.yml', got %s", CreateConfigPath)
	}

	if ModeExplicit != "explicit" {
		t.Errorf("Expected ModeExplicit to be 'explicit', got %s", ModeExplicit)
	}

	if ModeDefault != "default" {
		t.Errorf("Expected ModeDefault to be 'default', got %s", ModeDefault)
	}
}

// TestDefaultValues tests that default configuration values are properly defined
func TestDefaultValues(t *testing.T) {
	if len(DefaultInclude) != 1 || DefaultInclude[0] != "." {
		t.Errorf("Expected DefaultInclude to be ['.'], got %v", DefaultInclude)
	}

	if len(DefaultCommands) != 0 {
		t.Errorf("Expected DefaultCommands to be empty, got %v", DefaultCommands)
	}
}
