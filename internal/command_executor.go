package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CommandExecutor interface for command execution
type CommandExecutor interface {
	// ExecuteCommands runs all commands and returns output map
	ExecuteCommands(commands []string) (map[string]string, error)

	// ExecuteCommand runs a single command
	ExecuteCommand(command string) (stdout, stderr string, err error)

	// ExecuteCommandsWithTimeout runs commands with timeout
	ExecuteCommandsWithTimeout(commands []string, timeout time.Duration) (map[string]string, error)
}

// CommandResult represents the result of executing a shell command
type CommandResult struct {
	Command string
	Stdout  string
	Stderr  string
	Success bool
	Error   error
}

// CommandExecutorImpl implements the CommandExecutor interface
type CommandExecutorImpl struct {
	timeout time.Duration
}

// NewCommandExecutor creates a new command executor with default timeout
func NewCommandExecutor() *CommandExecutorImpl {
	return &CommandExecutorImpl{
		timeout: getCommandTimeout(),
	}
}

// NewCommandExecutorWithTimeout creates a command executor with custom timeout
func NewCommandExecutorWithTimeout(timeout time.Duration) *CommandExecutorImpl {
	return &CommandExecutorImpl{
		timeout: timeout,
	}
}

// ExecuteCommand runs a single command
// Returns stdout, stderr, and error
func (c *CommandExecutorImpl) ExecuteCommand(command string) (stdout, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.executeCommandWithContext(ctx, command)
}

// executeCommandWithContext runs a command with a specific context
func (c *CommandExecutorImpl) executeCommandWithContext(ctx context.Context, command string) (stdout, stderr string, err error) {
	// Use shell to execute command to support shell features like redirection, pipes, etc.
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// Capture stdout and stderr separately using pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	// Read from both pipes concurrently
	stdoutCh := make(chan []byte, 1)
	stderrCh := make(chan []byte, 1)
	errCh := make(chan error, 2)

	go func() {
		data, err := io.ReadAll(stdoutPipe)
		if err != nil {
			errCh <- err
		}
		stdoutCh <- data
	}()

	go func() {
		data, err := io.ReadAll(stderrPipe)
		if err != nil {
			errCh <- err
		}
		stderrCh <- data
	}()

	// Wait for command to complete
	cmdErr := cmd.Wait()

	// Get the results
	stdoutData := <-stdoutCh
	stderrData := <-stderrCh

	// Check for read errors
	select {
	case readErr := <-errCh:
		if cmdErr == nil {
			cmdErr = readErr
		}
	default:
		// No read errors
	}

	return strings.TrimSpace(string(stdoutData)), strings.TrimSpace(string(stderrData)), cmdErr
}

// ExecuteCommands runs all commands in the configuration
// Returns combined output and any execution errors
func (c *CommandExecutorImpl) ExecuteCommands(commands []string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.executeCommandsWithContext(ctx, commands)
}

// ExecuteCommandsWithTimeout runs commands with a specific timeout
func (c *CommandExecutorImpl) ExecuteCommandsWithTimeout(commands []string, timeout time.Duration) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.executeCommandsWithContext(ctx, commands)
}

// executeCommandsWithContext runs commands with a specific context
func (c *CommandExecutorImpl) executeCommandsWithContext(ctx context.Context, commands []string) (map[string]string, error) {
	results := make(map[string]string)

	for _, command := range commands {
		stdout, stderr, err := c.executeCommandWithContext(ctx, command)
		if err != nil {
			log.Printf("Command failed: %s, error: %v", command, err)
			return nil, err
		}

		// Combine stdout and stderr for the result
		output := stdout
		if stderr != "" {
			if output != "" {
				output += "\n" + stderr
			} else {
				output = stderr
			}
		}

		results[command] = output
	}

	return results, nil
}

// ExecuteCommandsDetailed runs commands and returns detailed results
func (c *CommandExecutorImpl) ExecuteCommandsDetailed(commands []string) []CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	results := make([]CommandResult, 0, len(commands))

	for _, command := range commands {
		stdout, stderr, err := c.executeCommandWithContext(ctx, command)

		result := CommandResult{
			Command: command,
			Stdout:  stdout,
			Stderr:  stderr,
			Success: err == nil,
			Error:   err,
		}

		// Log failures but continue execution
		if err != nil {
			log.Printf("Command failed: %s, error: %v", command, err)
		}

		results = append(results, result)
	}

	return results
}

// FormatCommandOutput formats the output of successful commands with headers
func (c *CommandExecutorImpl) FormatCommandOutput(results []CommandResult) string {
	var output strings.Builder

	for _, result := range results {
		if result.Success && strings.TrimSpace(result.Stdout) != "" {
			output.WriteString(fmt.Sprintf("=== Command: %s ===\n", result.Command))
			output.WriteString(strings.TrimSpace(result.Stdout))
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
