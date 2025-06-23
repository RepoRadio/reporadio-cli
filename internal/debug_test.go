package internal

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDebugDisabledByDefault(t *testing.T) {
	var buf bytes.Buffer

	Debug(&buf, "test message")

	if buf.Len() > 0 {
		t.Error("Debug should be disabled by default")
	}
}

func TestDebugEnabled(t *testing.T) {
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	var buf bytes.Buffer

	Debug(&buf, "test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected debug message, got: %s", output)
	}
}

func TestDebugf(t *testing.T) {
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	var buf bytes.Buffer

	Debugf(&buf, "tokens: %d", 150)

	output := buf.String()
	if !strings.Contains(output, "tokens: 150") {
		t.Error("Expected formatted debug message")
	}
}

func TestDebugOpenAIRequest(t *testing.T) {
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	var buf bytes.Buffer

	messages := []map[string]string{
		{"role": "user", "content": "test message"},
	}

	DebugOpenAIRequest(&buf, messages, "gpt-4")

	output := buf.String()
	if !strings.Contains(output, "OpenAI Request") {
		t.Error("Expected OpenAI Request in debug output")
	}
}
