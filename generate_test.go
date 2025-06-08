package main

import (
	"testing"
)

func TestRunGenerate(t *testing.T) {
	// Test the runGenerate function directly
	err := runGenerate(generateCmd, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test with argument
	err = runGenerate(generateCmd, []string{"show-notes"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
