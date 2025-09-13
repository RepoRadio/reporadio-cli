package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Test interface compliance
func TestCommandExecutorInterface(t *testing.T) {
	var _ CommandExecutor = NewCommandExecutor()
}

func TestCommandExecutor_ExecuteCommand(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		expectError    bool
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "successful command",
			command:        "echo 'hello world'",
			expectError:    false,
			expectedStdout: "hello world",
			expectedStderr: "",
		},
		{
			name:        "command not found",
			command:     "nonexistent_command_12345",
			expectError: true,
		},
		{
			name:        "command with error",
			command:     "sh -c 'exit 1'",
			expectError: true,
		},
		{
			name:           "command with output to stderr",
			command:        "sh -c 'echo stderr message >&2'",
			expectError:    false,
			expectedStdout: "",
			expectedStderr: "stderr message",
		},
		{
			name:           "simple command with no output",
			command:        "true",
			expectError:    false,
			expectedStdout: "",
			expectedStderr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewCommandExecutor()
			stdout, stderr, err := executor.ExecuteCommand(tt.command)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.expectedStdout != "" && stdout != tt.expectedStdout {
				t.Errorf("expected stdout %q, got %q", tt.expectedStdout, stdout)
			}

			if tt.expectedStderr != "" && stderr != tt.expectedStderr {
				t.Errorf("expected stderr %q, got %q", tt.expectedStderr, stderr)
			}
		})
	}
}

func TestCommandExecutor_ExecuteCommands(t *testing.T) {
	tests := []struct {
		name         string
		commands     []string
		expectError  bool
		validateFunc func(t *testing.T, results map[string]string)
	}{
		{
			name:     "multiple successful commands",
			commands: []string{"echo 'first'", "echo 'second'"},
			validateFunc: func(t *testing.T, results map[string]string) {
				if len(results) != 2 {
					t.Errorf("expected 2 results, got %d", len(results))
				}
				if results["echo 'first'"] != "first" {
					t.Errorf("expected 'first', got %q", results["echo 'first'"])
				}
				if results["echo 'second'"] != "second" {
					t.Errorf("expected 'second', got %q", results["echo 'second'"])
				}
			},
		},
		{
			name:     "empty commands list",
			commands: []string{},
			validateFunc: func(t *testing.T, results map[string]string) {
				if len(results) != 0 {
					t.Errorf("expected 0 results, got %d", len(results))
				}
			},
		},
		{
			name:        "one command fails",
			commands:    []string{"echo 'success'", "nonexistent_command_12345"},
			expectError: true,
		},
		{
			name:     "commands with file redirection",
			commands: []string{"echo 'test content' > temp_test_file.txt", "cat temp_test_file.txt"},
			validateFunc: func(t *testing.T, results map[string]string) {
				defer os.Remove("temp_test_file.txt") // Clean up

				if len(results) != 2 {
					t.Errorf("expected 2 results, got %d", len(results))
				}
				if results["cat temp_test_file.txt"] != "test content" {
					t.Errorf("expected 'test content', got %q", results["cat temp_test_file.txt"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewCommandExecutor()
			results, err := executor.ExecuteCommands(tt.commands)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, results)
			}
		})
	}
}

func TestCommandExecutor_ExecuteCommandsDetailed(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     []CommandResult
	}{
		{
			name:     "single successful command",
			commands: []string{"echo hello"},
			want: []CommandResult{
				{Command: "echo hello", Stdout: "hello", Success: true},
			},
		},
		{
			name:     "multiple successful commands",
			commands: []string{"echo first", "echo second"},
			want: []CommandResult{
				{Command: "echo first", Stdout: "first", Success: true},
				{Command: "echo second", Stdout: "second", Success: true},
			},
		},
		{
			name:     "mix of successful and failed commands",
			commands: []string{"echo success", "nonexistentcommand123", "echo after-failure"},
			want: []CommandResult{
				{Command: "echo success", Stdout: "success", Success: true},
				{Command: "nonexistentcommand123", Stdout: "", Success: false},
				{Command: "echo after-failure", Stdout: "after-failure", Success: true},
			},
		},
		{
			name:     "empty commands",
			commands: []string{},
			want:     []CommandResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewCommandExecutor()
			results := executor.ExecuteCommandsDetailed(tt.commands)

			if len(results) != len(tt.want) {
				t.Errorf("expected %d results, got %d", len(tt.want), len(results))
				return
			}

			for i, want := range tt.want {
				if i >= len(results) {
					t.Errorf("missing result at index %d", i)
					continue
				}

				got := results[i]
				if got.Command != want.Command {
					t.Errorf("command mismatch at index %d: want %q, got %q", i, want.Command, got.Command)
				}
				if got.Success != want.Success {
					t.Errorf("success mismatch at index %d: want %v, got %v", i, want.Success, got.Success)
				}
				if want.Success && got.Stdout != want.Stdout {
					t.Errorf("stdout mismatch at index %d: want %q, got %q", i, want.Stdout, got.Stdout)
				}
			}
		})
	}
}

func TestCommandExecutor_ExecuteCommandsWithTimeout(t *testing.T) {
	executor := NewCommandExecutor()

	t.Run("commands complete within timeout", func(t *testing.T) {
		commands := []string{"echo 'fast'", "echo 'quick'"}
		results, err := executor.ExecuteCommandsWithTimeout(commands, 5*time.Second)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("command times out", func(t *testing.T) {
		commands := []string{"sleep 2"}
		_, err := executor.ExecuteCommandsWithTimeout(commands, 100*time.Millisecond)

		if err == nil {
			t.Error("expected timeout error but got none")
		}

		if !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("expected timeout-related error, got: %v", err)
		}
	})
}

func TestCommandExecutor_Integration(t *testing.T) {
	t.Run("command generates file for inclusion", func(t *testing.T) {
		tempDir := t.TempDir()
		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		executor := NewCommandExecutor()

		// Command that creates a file
		commands := []string{
			"echo 'Generated content' > generated.txt",
			"cat generated.txt",
		}

		results, err := executor.ExecuteCommands(commands)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}

		if results["cat generated.txt"] != "Generated content" {
			t.Errorf("expected 'Generated content', got %q", results["cat generated.txt"])
		}

		// Verify file exists
		if _, err := os.Stat("generated.txt"); os.IsNotExist(err) {
			t.Error("generated file should exist")
		}
	})

	t.Run("commands execute in repository root context", func(t *testing.T) {
		tempDir := t.TempDir()
		oldDir, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldDir)

		// Create a subdir to test working directory
		subDir := filepath.Join(tempDir, "subdir")
		os.MkdirAll(subDir, 0755)

		executor := NewCommandExecutor()

		commands := []string{"pwd"}
		results, err := executor.ExecuteCommands(commands)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		currentDir, _ := os.Getwd()
		if results["pwd"] != currentDir {
			t.Errorf("expected pwd to return %q, got %q", currentDir, results["pwd"])
		}
	})
}

func TestCommandExecutor_FormatCommandOutput(t *testing.T) {
	executor := NewCommandExecutor()

	results := []CommandResult{
		{Command: "echo hello", Stdout: "hello", Success: true},
		{Command: "false", Stdout: "", Success: false},
		{Command: "echo world", Stdout: "world", Success: true},
		{Command: "echo empty", Stdout: "", Success: true}, // Should be filtered out
	}

	output := executor.FormatCommandOutput(results)

	expectedLines := []string{
		"=== Command: echo hello ===",
		"hello",
		"=== Command: echo world ===",
		"world",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("expected output to contain %q, got:\n%s", line, output)
		}
	}

	// Should not contain failed command
	if strings.Contains(output, "false") {
		t.Error("output should not contain failed commands")
	}

	// Should not contain empty successful command
	if strings.Contains(output, "echo empty") {
		t.Error("output should not contain empty successful commands")
	}
}

func TestCommandTimeout(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "default timeout when no env var set",
			envValue: "",
			expected: 60 * time.Second,
		},
		{
			name:     "custom timeout from env var",
			envValue: "30s",
			expected: 30 * time.Second,
		},
		{
			name:     "custom timeout in seconds format",
			envValue: "45",
			expected: 45 * time.Second,
		},
		{
			name:     "invalid timeout format",
			envValue: "invalid",
			expected: 60 * time.Second, // Should fallback to default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("REPORADIO_COMMAND_TIMEOUT", tt.envValue)
				defer os.Unsetenv("REPORADIO_COMMAND_TIMEOUT")
			} else {
				os.Unsetenv("REPORADIO_COMMAND_TIMEOUT")
			}

			timeout := getCommandTimeout()
			if timeout != tt.expected {
				t.Errorf("expected timeout %v, got %v", tt.expected, timeout)
			}
		})
	}
}

func TestNewCommandExecutorWithTimeout(t *testing.T) {
	customTimeout := 30 * time.Second
	executor := NewCommandExecutorWithTimeout(customTimeout)

	if executor.timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, executor.timeout)
	}
}
