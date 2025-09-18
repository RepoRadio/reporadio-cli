package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// APIConfig holds configuration for OpenAI-compatible API endpoints
type APIConfig struct {
	// BaseURL is the API base URL (e.g., "https://api.openai.com/v1" or "http://localhost:1234/v1")
	BaseURL string

	// APIKey is the API key for authentication
	APIKey string

	// ChatModel is the model to use for chat completions (e.g., "gpt-4o", "llama-3.1-8b-instruct")
	ChatModel string

	// TTSModel is the model to use for text-to-speech (e.g., "tts-1", "tts-1-hd")
	// Leave empty to disable TTS functionality for non-OpenAI APIs
	TTSModel string

	// TTSVoice is the voice to use for TTS (e.g., "alloy", "echo", "fable")
	TTSVoice string
}

// DefaultAPIConfig returns the default configuration for OpenAI
func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		BaseURL:   "https://api.openai.com/v1",
		APIKey:    "",
		ChatModel: string(openai.GPT4oMini), // Default to mini for cost efficiency
		TTSModel:  string(openai.TTSModel1),
		TTSVoice:  string(openai.VoiceAlloy),
	}
}

// LoadAPIConfig loads API configuration from environment variables and flags
func LoadAPIConfig(baseURL, apiKey, chatModel, ttsModel, ttsVoice string) (*APIConfig, error) {
	config := DefaultAPIConfig()

	// Override with environment variables
	if envAPIKey := os.Getenv("OPENAI_API_KEY"); envAPIKey != "" {
		config.APIKey = envAPIKey
	}
	if envBaseURL := os.Getenv("OPENAI_BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	}
	if envChatModel := os.Getenv("OPENAI_CHAT_MODEL"); envChatModel != "" {
		config.ChatModel = envChatModel
	}
	if envTTSModel := os.Getenv("OPENAI_TTS_MODEL"); envTTSModel != "" {
		config.TTSModel = envTTSModel
	}
	if envTTSVoice := os.Getenv("OPENAI_TTS_VOICE"); envTTSVoice != "" {
		config.TTSVoice = envTTSVoice
	}

	// Override with command line flags (highest priority)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	if apiKey != "" {
		config.APIKey = apiKey
	}
	if chatModel != "" {
		config.ChatModel = chatModel
	}
	if ttsModel != "" {
		config.TTSModel = ttsModel
	}
	if ttsVoice != "" {
		config.TTSVoice = ttsVoice
	}

	// Validate required fields
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Set OPENAI_API_KEY environment variable or use --api-key flag")
	}
	if config.ChatModel == "" {
		return nil, fmt.Errorf("chat model is required. Set OPENAI_CHAT_MODEL environment variable or use --chat-model flag")
	}

	return config, nil
}

// NewClient creates a new OpenAI client with the configured settings
func (c *APIConfig) NewClient() *openai.Client {
	config := openai.DefaultConfig(c.APIKey)

	// Set custom base URL if provided
	if c.BaseURL != "" && c.BaseURL != "https://api.openai.com/v1" {
		config.BaseURL = c.BaseURL
	}

	return openai.NewClientWithConfig(config)
}

// SupportsTTS returns true if the API configuration supports text-to-speech
func (c *APIConfig) SupportsTTS() bool {
	return c.TTSModel != "" && !c.IsLocalAPI()
}

// IsLocalAPI returns true if this appears to be a local API endpoint
func (c *APIConfig) IsLocalAPI() bool {
	return strings.Contains(strings.ToLower(c.BaseURL), "localhost") ||
		strings.Contains(strings.ToLower(c.BaseURL), "127.0.0.1") ||
		strings.Contains(strings.ToLower(c.BaseURL), "0.0.0.0")
}

// GetChatModel returns the chat model as a string suitable for API requests
func (c *APIConfig) GetChatModel() string {
	return c.ChatModel
}

// GetTTSModel returns the TTS model as a string suitable for API requests
func (c *APIConfig) GetTTSModel() string {
	return c.TTSModel
}

// GetTTSVoice returns the TTS voice as a string suitable for API requests
func (c *APIConfig) GetTTSVoice() string {
	return c.TTSVoice
}

// String returns a human-readable representation of the config (without API key)
func (c *APIConfig) String() string {
	apiKeyDisplay := "***"
	if c.APIKey == "" {
		apiKeyDisplay = "<empty>"
	}

	return fmt.Sprintf("API Config: BaseURL=%s, APIKey=%s, ChatModel=%s, TTSModel=%s, TTSVoice=%s",
		c.BaseURL, apiKeyDisplay, c.ChatModel, c.TTSModel, c.TTSVoice)
}
