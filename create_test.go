package main

import (
	"testing"
)

func TestCreateCommand(t *testing.T) {
	t.Run("create command exists", func(t *testing.T) {
		if createCmd == nil {
			t.Error("createCmd should not be nil")
		}

		if createCmd.Use != "create [podcast-name]" {
			t.Errorf("expected command use to be 'create [podcast-name]', got %s", createCmd.Use)
		}
	})

	t.Run("create command has description", func(t *testing.T) {
		if createCmd.Short == "" {
			t.Error("createCmd should have a short description")
		}

		if createCmd.Long == "" {
			t.Error("createCmd should have a long description")
		}
	})

	t.Run("runCreate returns no error for basic execution", func(t *testing.T) {
		err := runCreate(createCmd, []string{})
		if err != nil {
			t.Errorf("runCreate should not return error, got: %v", err)
		}
	})

	t.Run("runCreate prints podcast name when provided", func(t *testing.T) {
		args := []string{"my-podcast"}
		err := runCreate(createCmd, args)
		if err != nil {
			t.Errorf("runCreate should not return error, got: %v", err)
		}
		// TODO: Capture and verify output contains "my-podcast"
	})
}
