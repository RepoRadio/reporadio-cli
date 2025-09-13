package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new RepoRadio configuration in the current directory",
	Long: `The init command sets up a new RepoRadio configuration in the current directory.
It creates the .reporadio directory structure and copies default template files.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	// Create .reporadio directory
	err := os.MkdirAll(RepoRadioDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create %s directory: %w", RepoRadioDir, err)
	}

	// Create .reporadio/prompts directory and copy templates
	promptsDir := filepath.Join(RepoRadioDir, "prompts")
	err = copyEmbeddedTemplates(promptsDir)
	if err != nil {
		return fmt.Errorf("failed to copy templates: %w", err)
	}

	fmt.Printf("✅ Initialized RepoRadio configuration in %s\n", RepoRadioDir)
	return nil
}

func copyEmbeddedTemplates(targetDir string) error {
	// Create target directory if it doesn't exist
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Read all embedded templates
	entries, err := promptFS.ReadDir("prompts")
	if err != nil {
		return fmt.Errorf("failed to read embedded templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		templateName := entry.Name()
		targetPath := filepath.Join(targetDir, templateName)

		// Skip if file already exists (preserve custom templates)
		if _, err := os.Stat(targetPath); err == nil {
			continue
		}

		// Read embedded template content
		content, err := promptFS.ReadFile("prompts/" + templateName)
		if err != nil {
			return fmt.Errorf("failed to read embedded template %s: %w", templateName, err)
		}

		// Write template to target directory
		err = os.WriteFile(targetPath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write template %s: %w", templateName, err)
		}
	}

	return nil
}

func init() {
	// Add init command to root command
	rootCmd.AddCommand(initCmd)
}
