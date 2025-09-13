package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// PodcastConfig represents the structure of a podcast.yml file
type PodcastConfig struct {
	Title        string     `yaml:"title"`
	Description  string     `yaml:"description"`
	Instructions string     `yaml:"instructions"`
	Voicing      string     `yaml:"voicing"`
	Type         SeriesType `yaml:"type"`
	Episodes     []Episode  `yaml:"episodes"`
}

// loadPodcastConfig loads a podcast configuration from .reporadio/<name>/podcast.yml
func loadPodcastConfig(name string) (*PodcastConfig, error) {
	configPath := filepath.Join(".reporadio", name, "podcast.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read podcast config %s: %w", configPath, err)
	}

	var config PodcastConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse podcast config %s: %w", configPath, err)
	}

	return &config, nil
}

// loadChatContext loads chat entries from .reporadio/<name>/chat.yaml
func loadChatContext(podcastName string) ([]map[string]interface{}, error) {
	chatPath := filepath.Join(".reporadio", podcastName, "chat.yaml")

	if _, err := os.Stat(chatPath); os.IsNotExist(err) {
		return []map[string]interface{}{}, nil
	}

	data, err := os.ReadFile(chatPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chat context: %w", err)
	}

	var chatFile map[string]interface{}
	err = yaml.Unmarshal(data, &chatFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chat context: %w", err)
	}

	if entries, ok := chatFile["entries"].([]interface{}); ok {
		result := make([]map[string]interface{}, len(entries))
		for i, entry := range entries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				result[i] = entryMap
			}
		}
		return result, nil
	}

	return []map[string]interface{}{}, nil
}

// appendToChatContext appends a new episode entry to the chat.yaml file
func appendToChatContext(podcastName string, episodeNum int, title, transcript string) error {
	chatPath := filepath.Join(".reporadio", podcastName, "chat.yaml")

	// Create directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(chatPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create chat context directory: %w", err)
	}

	var chatFile map[string]interface{}

	if data, err := os.ReadFile(chatPath); err == nil {
		if err := yaml.Unmarshal(data, &chatFile); err != nil {
			return fmt.Errorf("failed to parse existing chat context: %w", err)
		}
	} else {
		chatFile = make(map[string]interface{})
	}

	timestamp := time.Now().Format("2006-01-02T15:04:05.000000-07:00")
	newEntry := map[string]interface{}{
		"timestamp": timestamp,
		"role":      "assistant",
		"message":   fmt.Sprintf("Episode %d: %s\n\n%s", episodeNum, title, transcript),
		"step":      "episode",
	}

	if entries, ok := chatFile["entries"].([]interface{}); ok {
		chatFile["entries"] = append(entries, newEntry)
	} else {
		chatFile["entries"] = []interface{}{newEntry}
	}

	data, err := yaml.Marshal(chatFile)
	if err != nil {
		return fmt.Errorf("failed to marshal chat context: %w", err)
	}

	return os.WriteFile(chatPath, data, 0644)
}

// generateEpisodeTranscript generates a transcript for a single episode
func generateEpisodeTranscript(episode Episode, episodeNum int, outputDir string, client *openai.Client, chatEntries []map[string]interface{}, podcastConfig *PodcastConfig, commandTimeoutFlag string) error {
	Debug(nil, "Starting episode transcript generation")
	Debugf(nil, "Episode %d: %s", episodeNum, episode.Title)

	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	Debugf(nil, "Created output directory: %s", outputDir)

	// Resolve include paths using Scanner
	scanner := NewScanner()
	resolvedPaths, err := scanner.ResolveIncludePaths(episode.Include)
	if err != nil {
		return fmt.Errorf("failed to resolve include paths: %w", err)
	}
	Debugf(nil, "Resolved %d files from %d include paths", len(resolvedPaths), len(episode.Include))

	// Read all resolved files
	var fileContents strings.Builder
	for _, filePath := range resolvedPaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to read file %s: %v\n", filePath, err)
			continue
		}

		fileContents.WriteString(fmt.Sprintf("\n--- %s ---\n", filePath))
		fileContents.WriteString(string(content))
		fileContents.WriteString("\n")
	}
	Debugf(nil, "Read %d files successfully", len(resolvedPaths))

	// Execute commands if any are specified
	var commandOutput string
	if len(episode.Commands) > 0 {
		Debug(nil, "Executing episode commands")
		timeout := getCommandTimeoutWithFlag(commandTimeoutFlag)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		results := executeCommands(ctx, episode.Commands)
		commandOutput = formatCommandOutput(results)
		Debugf(nil, "Executed %d commands with %v timeout, got %d bytes of output", len(episode.Commands), timeout, len(commandOutput))
	}

	// If no client provided (for testing), just create a placeholder file
	if client == nil {
		transcriptPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.md", episodeNum))
		placeholderContent := fmt.Sprintf("# %s\n\nPlaceholder transcript for testing", episode.Title)
		return os.WriteFile(transcriptPath, []byte(placeholderContent), 0644)
	}

	// Build context from previous episodes
	var contextBuilder strings.Builder
	if len(chatEntries) > 0 {
		contextBuilder.WriteString("\n--- Previous Episodes Context ---\n")
		for _, entry := range chatEntries {
			if role, ok := entry["role"].(string); ok && role == "assistant" {
				if step, ok := entry["step"].(string); ok && step == "episode" {
					if message, ok := entry["message"].(string); ok {
						contextBuilder.WriteString(message)
						contextBuilder.WriteString("\n")
					}
				}
			}
		}
	}
	Debugf(nil, "Built context from %d chat entries", len(chatEntries))

	// Create the prompt using the template system
	Debug(nil, "Creating prompt from template")
	promptManager, err := NewPromptManager()
	if err != nil {
		return fmt.Errorf("failed to create prompt manager: %w", err)
	}

	// Use episode-specific values, falling back to podcast-level values
	title := episode.Title
	if title == "" {
		title = podcastConfig.Title
	}

	description := episode.Description
	if description == "" {
		description = podcastConfig.Description
	}

	instructions := episode.Instructions
	if instructions == "" {
		instructions = podcastConfig.Instructions
	}

	voicing := episode.Voicing
	if voicing == "" {
		voicing = podcastConfig.Voicing
	}

	templateData := struct {
		Title         string
		Description   string
		Instructions  string
		Voicing       string
		Context       string
		FileContents  string
		CommandOutput string
	}{
		Title:         title,
		Description:   description,
		Instructions:  instructions,
		Voicing:       voicing,
		Context:       contextBuilder.String(),
		FileContents:  fileContents.String(),
		CommandOutput: commandOutput,
	}

	prompt, err := promptManager.Execute("episode_transcript.tmpl", templateData)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	Debug(nil, "Prompt created successfully")

	// Create chat completion request
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// Get response from OpenAI
	Debug(nil, "Generating transcript with OpenAI")
	transcript, err := getChatResponse(client, messages)
	if err != nil {
		return fmt.Errorf("failed to generate transcript: %w", err)
	}
	Debugf(nil, "Generated transcript with %d characters", len(transcript))

	// Write transcript to file
	transcriptPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.md", episodeNum))
	err = os.WriteFile(transcriptPath, []byte(transcript), 0644)
	if err != nil {
		return fmt.Errorf("failed to write transcript file: %w", err)
	}
	Debugf(nil, "Saved transcript to: %s", transcriptPath)

	Debug(nil, "Episode transcript generation complete")
	return nil
}

// generatePodcastTranscripts generates transcripts and optionally audio for all episodes in a podcast config
func generatePodcastTranscripts(podcastName string, config *PodcastConfig, client *openai.Client, generateAudio bool, commandTimeoutFlag string) error {
	Debug(nil, "Starting podcast transcript generation")
	Debugf(nil, "Podcast: %s, Episodes: %d, Audio: %t", podcastName, len(config.Episodes), generateAudio)

	outputDir := filepath.Join(".reporadio", podcastName, "episodes")

	// Load existing chat context
	Debug(nil, "Loading chat context")
	chatEntries, err := loadChatContext(podcastName)
	if err != nil {
		return fmt.Errorf("failed to load chat context: %w", err)
	}
	Debugf(nil, "Loaded %d chat entries", len(chatEntries))

	fmt.Printf("Generating transcripts for %d episodes...\n", len(config.Episodes))
	if generateAudio {
		fmt.Printf("Audio generation enabled\n")
	}

	for i, episode := range config.Episodes {
		episodeNum := i + 1
		fmt.Printf("Generating Episode %d: %s\n", episodeNum, episode.Title)

		err := generateEpisodeTranscript(episode, episodeNum, outputDir, client, chatEntries, config, commandTimeoutFlag)
		if err != nil {
			return fmt.Errorf("failed to generate episode %d: %w", episodeNum, err)
		}

		fmt.Printf("✓ Generated ep%d.md\n", episodeNum)

		// Read generated transcript and append to chat context
		transcriptPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.md", episodeNum))
		transcriptContent, err := os.ReadFile(transcriptPath)
		if err != nil {
			return fmt.Errorf("failed to read generated transcript: %w", err)
		}

		Debug(nil, "Appending to chat context")
		err = appendToChatContext(podcastName, episodeNum, episode.Title, string(transcriptContent))
		if err != nil {
			return fmt.Errorf("failed to append to chat context: %w", err)
		}

		// Generate audio if requested and client is available
		if generateAudio && client != nil {
			fmt.Printf("Generating audio for Episode %d...\n", episodeNum)
			Debug(nil, "Starting audio generation")

			audioPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.mp3", episodeNum))

			err = generateEpisodeAudio(transcriptPath, audioPath, client)
			if err != nil {
				return fmt.Errorf("failed to generate audio for episode %d: %w", episodeNum, err)
			}

			fmt.Printf("✓ Generated ep%d.mp3\n", episodeNum)
		}
	}

	fmt.Printf("\nAll transcripts generated in %s\n", outputDir)
	if generateAudio {
		fmt.Printf("All audio files generated in %s\n", outputDir)
	}
	Debug(nil, "Podcast transcript generation complete")
	return nil
}

var generateCmd = &cobra.Command{
	Use:   "generate [podcast-name]",
	Short: "Generate a specific podcast",
	Long:  `The generate command generates the transcript and audio based on the podcast config.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// TODO: List available podcast configs
		return fmt.Errorf("please specify a podcast name")
	}

	podcastName := args[0]

	// Load podcast config
	config, err := loadPodcastConfig(podcastName)
	if err != nil {
		return err
	}

	// Get flags
	generateAudio, _ := cmd.Flags().GetBool("audio")
	commandTimeoutFlag, _ := cmd.Flags().GetString("command-timeout")

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("please set the OPENAI_API_KEY environment variable")
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Generate transcripts (and optionally audio) for all episodes
	return generatePodcastTranscripts(podcastName, config, client, generateAudio, commandTimeoutFlag)
}

// parseEpisodeNumbers parses a comma-separated list of episode numbers
func parseEpisodeNumbers(episodesList string) ([]int, error) {
	if episodesList == "" {
		return []int{}, nil
	}

	parts := strings.Split(episodesList, ",")
	episodes := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid episode number: %s", part)
		}

		if num < 1 {
			return nil, fmt.Errorf("episode numbers must be positive: %d", num)
		}

		episodes = append(episodes, num)
	}

	return episodes, nil
}

// validateEpisodeExists checks if an episode transcript file exists
func validateEpisodeExists(podcastName string, episodeNum int) error {
	transcriptPath := filepath.Join(".reporadio", podcastName, "episodes", fmt.Sprintf("ep%d.md", episodeNum))

	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		return fmt.Errorf("episode %d transcript not found at %s", episodeNum, transcriptPath)
	}

	return nil
}

// regenerateAudioForEpisodes regenerates audio for specified episodes
func regenerateAudioForEpisodes(podcastName string, episodes []int, client *openai.Client) error {
	outputDir := filepath.Join(".reporadio", podcastName, "episodes")

	// Validate all episodes exist before starting
	for _, episodeNum := range episodes {
		if err := validateEpisodeExists(podcastName, episodeNum); err != nil {
			return err
		}
	}

	fmt.Printf("Regenerating audio for %d episodes...\n", len(episodes))

	for _, episodeNum := range episodes {
		fmt.Printf("Generating audio for Episode %d...\n", episodeNum)

		transcriptPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.md", episodeNum))
		audioPath := filepath.Join(outputDir, fmt.Sprintf("ep%d.mp3", episodeNum))

		err := generateEpisodeAudio(transcriptPath, audioPath, client)
		if err != nil {
			return fmt.Errorf("failed to generate audio for episode %d: %w", episodeNum, err)
		}

		fmt.Printf("✓ Generated ep%d.mp3\n", episodeNum)
	}

	fmt.Printf("\nAudio regeneration complete for %d episodes in %s\n", len(episodes), outputDir)
	return nil
}

var regenerateAudioCmd = &cobra.Command{
	Use:   "regenerate-audio [podcast-name]",
	Short: "Regenerate audio files from existing episode transcripts",
	Long: `The regenerate-audio command generates audio files from existing episode transcript files (.md) 
without regenerating the transcripts themselves. This is useful when you have manually edited 
the transcript files and want to regenerate only the audio.

Examples:
  # Regenerate all audio files
  reporadio regenerate-audio my-podcast
  
  # Regenerate specific episodes
  reporadio regenerate-audio my-podcast --episodes 1,3,5
  
  # Regenerate a single episode
  reporadio regenerate-audio my-podcast --episode 2`,
	Args: cobra.ExactArgs(1),
	RunE: runRegenerateAudio,
}

func runRegenerateAudio(cmd *cobra.Command, args []string) error {
	podcastName := args[0]

	// Load podcast config to get episode count
	config, err := loadPodcastConfig(podcastName)
	if err != nil {
		return err
	}

	// Get flags
	episodesList, err := cmd.Flags().GetString("episodes")
	if err != nil {
		return fmt.Errorf("failed to retrieve 'episodes' flag: %w", err)
	}
	singleEpisode, err := cmd.Flags().GetInt("episode")
	if err != nil {
		return fmt.Errorf("failed to retrieve 'episode' flag: %w", err)
	}

	// Determine which episodes to regenerate
	var episodes []int

	if singleEpisode > 0 {
		// Single episode specified
		episodes = []int{singleEpisode}
	} else if episodesList != "" {
		// Multiple episodes specified
		episodes, err = parseEpisodeNumbers(episodesList)
		if err != nil {
			return fmt.Errorf("error parsing episodes: %w", err)
		}
	} else {
		// No episodes specified, regenerate all
		episodes = make([]int, len(config.Episodes))
		for i := range config.Episodes {
			episodes[i] = i + 1
		}
	}

	if len(episodes) == 0 {
		return fmt.Errorf("no episodes specified for regeneration")
	}

	// Validate episode numbers are within range
	for _, episodeNum := range episodes {
		if episodeNum > len(config.Episodes) {
			return fmt.Errorf("episode %d does not exist (podcast has %d episodes)", episodeNum, len(config.Episodes))
		}
	}

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("please set the OPENAI_API_KEY environment variable")
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Regenerate audio for specified episodes
	return regenerateAudioForEpisodes(podcastName, episodes, client)
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(regenerateAudioCmd)

	// Add flags for the generate command
	generateCmd.Flags().Bool("audio", false, "Generate audio files in addition to transcripts")
	generateCmd.Flags().String("command-timeout", "", "Timeout for command execution (e.g., '30s', '2m', or '120' for seconds)")

	// Add flags for the regenerate-audio command
	regenerateAudioCmd.Flags().String("episodes", "", "Comma-separated list of episode numbers (e.g., '1,3,5')")
	regenerateAudioCmd.Flags().Int("episode", 0, "Single episode number to regenerate")
}
