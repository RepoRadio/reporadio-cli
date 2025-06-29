package internal

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestExecuteCommands(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     []CommandResult
	}{
		{
			name:     "single successful command",
			commands: []string{"echo hello"},
			want: []CommandResult{
				{Command: "echo hello", Output: "hello", Success: true},
			},
		},
		{
			name:     "multiple successful commands",
			commands: []string{"echo first", "echo second"},
			want: []CommandResult{
				{Command: "echo first", Output: "first", Success: true},
				{Command: "echo second", Output: "second", Success: true},
			},
		},
		{
			name:     "mix of successful and failed commands",
			commands: []string{"echo success", "nonexistentcommand123", "echo after-failure"},
			want: []CommandResult{
				{Command: "echo success", Output: "success", Success: true},
				{Command: "nonexistentcommand123", Output: "", Success: false},
				{Command: "echo after-failure", Output: "after-failure", Success: true},
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
			ctx := context.Background()
			results := executeCommands(ctx, tt.commands)

			if len(results) != len(tt.want) {
				t.Errorf("executeCommands() returned %d results, want %d", len(results), len(tt.want))
				return
			}

			for i, result := range results {
				expected := tt.want[i]
				if result.Command != expected.Command {
					t.Errorf("result[%d].Command = %q, want %q", i, result.Command, expected.Command)
				}
				if result.Success != expected.Success {
					t.Errorf("result[%d].Success = %v, want %v", i, result.Success, expected.Success)
				}
				if result.Success && strings.TrimSpace(result.Output) != expected.Output {
					t.Errorf("result[%d].Output = %q, want %q", i, strings.TrimSpace(result.Output), expected.Output)
				}
			}
		})
	}
}

func TestExecuteCommandsWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Command that should timeout
	commands := []string{"sleep 1"}
	results := executeCommands(ctx, commands)

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	result := results[0]
	if result.Success {
		t.Error("Expected command to fail due to timeout, but it succeeded")
	}
	if result.Command != "sleep 1" {
		t.Errorf("Expected command 'sleep 1', got %q", result.Command)
	}
}

func TestFormatCommandOutput(t *testing.T) {
	results := []CommandResult{
		{Command: "echo hello", Output: "hello", Success: true},
		{Command: "echo world", Output: "world", Success: true},
		{Command: "failed-cmd", Output: "", Success: false},
	}

	output := formatCommandOutput(results)

	// Should only include successful commands
	if !strings.Contains(output, "=== Command: echo hello ===") {
		t.Error("Output should contain header for first command")
	}
	if !strings.Contains(output, "hello") {
		t.Error("Output should contain result of first command")
	}
	if !strings.Contains(output, "=== Command: echo world ===") {
		t.Error("Output should contain header for second command")
	}
	if !strings.Contains(output, "world") {
		t.Error("Output should contain result of second command")
	}
	if strings.Contains(output, "failed-cmd") {
		t.Error("Output should not contain failed command")
	}
}
