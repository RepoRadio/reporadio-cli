package internal

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

// convertTextToSpeech converts text to speech using OpenAI TTS
func convertTextToSpeech(client *openai.Client, text string, filename string) error {
	if client == nil {
		return fmt.Errorf("OpenAI client is required for text-to-speech conversion")
	}

	// Create text-to-speech request
	resp, err := client.CreateSpeech(context.Background(), openai.CreateSpeechRequest{
		Model:          openai.TTSModelGPT4oMini,
		Input:          text,
		Voice:          openai.VoiceAlloy,
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

// generateEpisodeAudio generates audio for a single episode transcript
func generateEpisodeAudio(transcriptPath, audioPath string, client *openai.Client) error {
	if client == nil {
		return fmt.Errorf("OpenAI client is required for audio generation")
	}

	// Read transcript content
	content, err := os.ReadFile(transcriptPath)
	if err != nil {
		return fmt.Errorf("failed to read transcript %s: %w", transcriptPath, err)
	}

	// Convert text to speech
	return convertTextToSpeech(client, string(content), audioPath)
}
