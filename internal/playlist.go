package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// PlaylistEpisode represents an episode in the playlist
type PlaylistEpisode struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AudioFile      string `json:"audioFile"`
	TranscriptFile string `json:"transcriptFile"`
	Summary        string `json:"summary"`
}

// Playlist represents the main playlist structure
type Playlist struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Episodes    []PlaylistEpisode `json:"episodes"`
}

// PlaylistRoot represents the root structure of the playlist.json file
type PlaylistRoot struct {
	Playlist Playlist `json:"playlist"`
}

// generateEpisodeSummary generates a concise summary from the episode transcript
func generateEpisodeSummary(client *openai.Client, transcript string) (string, error) {
	if client == nil {
		return "Summary unavailable (no OpenAI client)", nil
	}

	// Create a focused prompt for summarization
	prompt := fmt.Sprintf(`Please provide a concise 1-2 sentence summary of the following podcast episode transcript:

%s

Summary:`, transcript)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// Get summary from OpenAI
	summary, err := getChatResponse(client, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate episode summary: %w", err)
	}

	// Clean up the summary (remove extra whitespace and newlines)
	summary = strings.TrimSpace(summary)
	summary = strings.ReplaceAll(summary, "\n", " ")

	return summary, nil
}

// generatePlaylistFile creates the playlist.json file with episode metadata
func generatePlaylistFile(podcastName string, config *PodcastConfig, episodeData []PlaylistEpisode) error {
	// Create the playlist structure
	playlist := PlaylistRoot{
		Playlist: Playlist{
			Title:       config.Title,
			Description: config.Description,
			Episodes:    episodeData,
		},
	}

	// Convert to JSON with proper formatting
	jsonData, err := json.MarshalIndent(playlist, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal playlist data: %w", err)
	}

	// Write to the playlist.json file
	playlistPath := filepath.Join(".reporadio", podcastName, "episodes", "playlist.json")
	err = os.WriteFile(playlistPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write playlist.json: %w", err)
	}

	return nil
}
