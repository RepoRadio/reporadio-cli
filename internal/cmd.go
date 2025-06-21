package internal

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reporadio",
	Short: "A CLI tool that helps users turn their codebases into podcast episodes",
	Long: `RepoRadio is a command-line tool that guides users through creating 
podcast metadata, episode lists, and chat history from a codebase.

It provides interactive onboarding and content generation capabilities
to transform your repository into engaging podcast content.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	rootCmd.AddCommand(createCmd)
}
