package internal

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CreateConfig represents the structure of .reporadio/create.yml
type CreateConfig struct {
	Include  []string `yaml:"include,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
}

// CreateContext represents the combined context for chat integration
type CreateContext struct {
	CommandOutput map[string]string // command -> output
	FileContent   map[string]string // path -> content
	Config        *CreateConfig     // original config
	Mode          string            // "explicit" or "default"
}

// Configuration constants
const (
	CreateConfigPath = ".reporadio/create.yml"
	ModeExplicit     = "explicit"
	ModeDefault      = "default"
)

// Default configuration values
var (
	DefaultInclude  = []string{"."}
	DefaultCommands = []string{}
)

// Configuration errors
var (
	ErrInvalidYAML    = errors.New("invalid YAML syntax in create.yml")
	ErrInvalidPattern = errors.New("invalid glob pattern in include list")
	ErrCommandFailed  = errors.New("command execution failed")
	ErrConfigNotFound = errors.New("create.yml not found")
)

// ConfigLoaderImpl implements the ConfigLoader interface
type ConfigLoaderImpl struct{}

// NewConfigLoader creates a new config loader
func NewConfigLoader() *ConfigLoaderImpl {
	return &ConfigLoaderImpl{}
}

// LoadCreateConfig loads configuration from .reporadio/create.yml
// Returns default config if file doesn't exist (backward compatibility)
func (c *ConfigLoaderImpl) LoadCreateConfig() (*CreateConfig, string, error) {
	configPath := CreateConfigPath

	// Check if config file exists
	if !c.HasCreateConfig() {
		// Return default configuration
		return &CreateConfig{
			Include:  append([]string{}, DefaultInclude...),
			Commands: append([]string{}, DefaultCommands...),
		}, ModeDefault, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, "", err
	}

	// Parse YAML
	var config CreateConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, "", ErrInvalidYAML
	}

	// Validate config
	if err := c.ValidateConfig(&config); err != nil {
		return nil, "", err
	}

	return &config, ModeExplicit, nil
}

// ValidateConfig validates the loaded configuration
func (c *ConfigLoaderImpl) ValidateConfig(config *CreateConfig) error {
	// Validate include patterns
	for _, pattern := range config.Include {
		if _, err := filepath.Match(pattern, "test"); err != nil {
			return ErrInvalidPattern
		}
	}

	// Validate commands (ensure no empty strings)
	for _, cmd := range config.Commands {
		if len(cmd) == 0 {
			return errors.New("empty command string")
		}
	}

	return nil
}

// HasCreateConfig checks if .reporadio/create.yml exists
func (c *ConfigLoaderImpl) HasCreateConfig() bool {
	_, err := os.Stat(CreateConfigPath)
	return err == nil
}
