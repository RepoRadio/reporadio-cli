package internal

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCommandTimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		envValue    string
		expected    time.Duration
		shouldError bool
	}{
		{
			name:     "default timeout when no env var set",
			expected: 60 * time.Second,
		},
		{
			name:     "custom timeout from env var",
			envVar:   "REPORADIO_COMMAND_TIMEOUT",
			envValue: "30",
			expected: 30 * time.Second,
		},
		{
			name:     "custom timeout in seconds format",
			envVar:   "REPORADIO_COMMAND_TIMEOUT",
			envValue: "120s",
			expected: 120 * time.Second,
		},
		{
			name:        "invalid timeout format",
			envVar:      "REPORADIO_COMMAND_TIMEOUT",
			envValue:    "invalid",
			expected:    60 * time.Second, // Should fall back to default
			shouldError: false,            // Should log but not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			if tt.envVar != "" {
				originalValue := os.Getenv(tt.envVar)
				defer func() {
					if originalValue == "" {
						os.Unsetenv(tt.envVar)
					} else {
						os.Setenv(tt.envVar, originalValue)
					}
				}()

				if tt.envValue != "" {
					os.Setenv(tt.envVar, tt.envValue)
				} else {
					os.Unsetenv(tt.envVar)
				}
			}

			timeout := getCommandTimeout()
			if timeout != tt.expected {
				t.Errorf("Expected timeout %v, got %v", tt.expected, timeout)
			}
		})
	}
}

func TestExecuteCommandsWithCustomTimeout(t *testing.T) {
	// Set a very short timeout
	originalValue := os.Getenv("REPORADIO_COMMAND_TIMEOUT")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("REPORADIO_COMMAND_TIMEOUT")
		} else {
			os.Setenv("REPORADIO_COMMAND_TIMEOUT", originalValue)
		}
	}()

	os.Setenv("REPORADIO_COMMAND_TIMEOUT", "100ms")

	commands := []string{"sleep 1"} // This should timeout
	timeout := getCommandTimeout()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results := executeCommands(ctx, commands)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Success {
		t.Error("Command should have failed due to timeout")
	}
}

func TestExecuteCommandsWithTimeoutOverride(t *testing.T) {
	// Test that executeCommandsWithTimeoutOverride works correctly
	commands := []string{"echo 'test'"}
	customTimeout := 5 * time.Second

	results := executeCommandsWithTimeoutOverride(commands, customTimeout)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if !results[0].Success {
		t.Error("Command should have succeeded")
	}
}

func TestCommandTimeoutWithFlag(t *testing.T) {
	tests := []struct {
		name      string
		flagValue string
		envValue  string
		expected  time.Duration
	}{
		{
			name:      "flag overrides environment",
			flagValue: "30s",
			envValue:  "60",
			expected:  30 * time.Second,
		},
		{
			name:      "flag value in seconds",
			flagValue: "45",
			expected:  45 * time.Second,
		},
		{
			name:      "invalid flag falls back to env",
			flagValue: "invalid",
			envValue:  "90",
			expected:  90 * time.Second,
		},
		{
			name:      "invalid flag and env falls back to default",
			flagValue: "invalid",
			envValue:  "also-invalid",
			expected:  60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			originalValue := os.Getenv("REPORADIO_COMMAND_TIMEOUT")
			defer func() {
				if originalValue == "" {
					os.Unsetenv("REPORADIO_COMMAND_TIMEOUT")
				} else {
					os.Setenv("REPORADIO_COMMAND_TIMEOUT", originalValue)
				}
			}()

			if tt.envValue != "" {
				os.Setenv("REPORADIO_COMMAND_TIMEOUT", tt.envValue)
			} else {
				os.Unsetenv("REPORADIO_COMMAND_TIMEOUT")
			}

			timeout := getCommandTimeoutWithFlag(tt.flagValue)
			if timeout != tt.expected {
				t.Errorf("Expected timeout %v, got %v", tt.expected, timeout)
			}
		})
	}
}
