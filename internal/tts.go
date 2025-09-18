package internal

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

// convertTextToSpeech converts text to speech using OpenAI TTS (legacy function)
func convertTextToSpeech(client *openai.Client, text string, filename string) error {
	// Use default API config for backward compatibility
	apiConfig := DefaultAPIConfig()
	return convertTextToSpeechWithConfig(client, apiConfig, text, filename)
}

// convertTextToSpeechWithConfig converts text to speech using configurable TTS API
func convertTextToSpeechWithConfig(client *openai.Client, apiConfig *APIConfig, text string, filename string) error {
	if client == nil {
		return fmt.Errorf("API client is required for text-to-speech conversion")
	}

	if !apiConfig.SupportsTTS() {
		return fmt.Errorf("text-to-speech is not supported with the current API configuration (endpoint: %s)", apiConfig.BaseURL)
	}

	// Create text-to-speech request
	resp, err := client.CreateSpeech(context.Background(), openai.CreateSpeechRequest{
		Model:          openai.SpeechModel(apiConfig.GetTTSModel()),
		Input:          text,
		Voice:          openai.SpeechVoice(apiConfig.GetTTSVoice()),
		ResponseFormat: openai.SpeechResponseFormatMp3,
		Speed:          1.0,
	})

	if err != nil {
		return fmt.Errorf("error creating speech: %v", err)
	}
	defer func() { _ = resp.Close() }()

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Copy the audio data to file
	_, err = io.Copy(file, resp)
	if err != nil {
		return fmt.Errorf("error writing audio data: %v", err)
	}

	return nil
}

// generateEpisodeAudio generates audio for a single episode transcript (legacy function)
func generateEpisodeAudio(transcriptPath, audioPath string, client *openai.Client) error {
	// Use default API config for backward compatibility
	apiConfig := DefaultAPIConfig()
	return generateEpisodeAudioWithConfig(transcriptPath, audioPath, client, apiConfig)
}

// generateEpisodeAudioWithConfig generates audio for a single episode transcript with API config
func generateEpisodeAudioWithConfig(transcriptPath, audioPath string, client *openai.Client, apiConfig *APIConfig) error {
	if client == nil {
		return fmt.Errorf("API client is required for audio generation")
	}

	if !apiConfig.SupportsTTS() {
		return fmt.Errorf("text-to-speech is not supported with the current API configuration (endpoint: %s)", apiConfig.BaseURL)
	}

	// Read transcript content
	content, err := os.ReadFile(transcriptPath)
	if err != nil {
		return fmt.Errorf("failed to read transcript %s: %w", transcriptPath, err)
	}

	// Convert text to speech
	return convertTextToSpeechWithConfig(client, apiConfig, string(content), audioPath)
}
