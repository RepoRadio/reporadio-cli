package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CommandResult represents the result of executing a shell command
type CommandResult struct {
	Command string
	Output  string
	Success bool
	Error   error
}

// executeCommands runs a list of shell commands sequentially and returns their results
func executeCommands(ctx context.Context, commands []string) []CommandResult {
	results := make([]CommandResult, 0, len(commands))

	for _, command := range commands {
		result := executeCommand(ctx, command)
		results = append(results, result)
	}

	return results
}

// executeCommand runs a single shell command and returns its result
func executeCommand(ctx context.Context, command string) CommandResult {
	// Split command into parts for exec.CommandContext
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return CommandResult{
			Command: command,
			Output:  "",
			Success: false,
			Error:   fmt.Errorf("empty command"),
		}
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	// Capture stdout only (stderr is ignored per requirements)
	output, err := cmd.Output()

	result := CommandResult{
		Command: command,
		Output:  string(output),
		Success: err == nil,
		Error:   err,
	}

	// Log failures but continue execution
	if err != nil {
		log.Printf("Command failed: %s, error: %v", command, err)
	}

	return result
}

// formatCommandOutput formats the output of successful commands with headers
func formatCommandOutput(results []CommandResult) string {
	var output strings.Builder

	for _, result := range results {
		if result.Success && strings.TrimSpace(result.Output) != "" {
			output.WriteString(fmt.Sprintf("=== Command: %s ===\n", result.Command))
			output.WriteString(strings.TrimSpace(result.Output))
			output.WriteString("\n\n")
		}
	}

	return output.String()
}

// getCommandTimeout returns the command timeout duration from environment variable or default
func getCommandTimeout() time.Duration {
	return getCommandTimeoutWithFlag("")
}

// getCommandTimeoutWithFlag returns the command timeout duration from flag, environment variable, or default
func getCommandTimeoutWithFlag(flagValue string) time.Duration {
	// First check command-line flag
	if flagValue != "" {
		if duration, err := parseTimeoutValue(flagValue); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid command timeout flag value '%s', checking environment variable", flagValue)
	}

	// Then check environment variable
	envValue := os.Getenv("REPORADIO_COMMAND_TIMEOUT")
	if envValue != "" {
		if duration, err := parseTimeoutValue(envValue); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid REPORADIO_COMMAND_TIMEOUT value '%s', using default 60 seconds", envValue)
	}

	// Default timeout
	return 60 * time.Second
}

// parseTimeoutValue parses a timeout value from string
func parseTimeoutValue(value string) (time.Duration, error) {
	// Try parsing as duration first (e.g., "30s", "2m")
	if duration, err := time.ParseDuration(value); err == nil {
		return duration, nil
	}

	// Try parsing as seconds (e.g., "30")
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	return 0, fmt.Errorf("invalid timeout format: %s", value)
}

// executeCommandsWithTimeout is a convenience function that executes commands with a default timeout
func executeCommandsWithTimeout(commands []string, timeout time.Duration) []CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return executeCommands(ctx, commands)
}

// executeCommandsWithTimeoutOverride executes commands with a custom timeout
func executeCommandsWithTimeoutOverride(commands []string, timeout time.Duration) []CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return executeCommands(ctx, commands)
}
