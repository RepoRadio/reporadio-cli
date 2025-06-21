package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/spf13/cobra"
)

// ConversationManager manages an interactive conversation with error recovery
type ConversationManager struct {
	client     *openai.Client
	scanner    *bufio.Scanner
	chatLog    *ChatLog
	messages   []openai.ChatCompletionMessage
	maxRetries int
}

// NewConversationManager creates a new conversation manager
func NewConversationManager(client *openai.Client, scanner *bufio.Scanner, chatLog *ChatLog, systemPrompt string) *ConversationManager {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	return &ConversationManager{
		client:     client,
		scanner:    scanner,
		chatLog:    chatLog,
		messages:   messages,
		maxRetries: 3,
	}
}

// AddError adds an error to the conversation context for recovery
func (cm *ConversationManager) AddError(err error) {
	errorMessage := fmt.Sprintf("An error occurred: %v. Let's restart our conversation with this context in mind.", err)
	cm.messages = append(cm.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: errorMessage,
	})
	cm.chatLog.AddEntry("system", errorMessage, "error")
}

// GetResponse gets a response from the assistant with retry logic
func (cm *ConversationManager) GetResponse() (string, error) {
	var lastErr error

	for attempt := 0; attempt < cm.maxRetries; attempt++ {
		response, err := getChatResponse(cm.client, cm.messages)
		if err == nil {
			return response, nil
		}

		lastErr = err
		if attempt < cm.maxRetries-1 {
			fmt.Printf("⚠️  Error (attempt %d/%d): %v\n", attempt+1, cm.maxRetries, err)
			fmt.Println("🔄 Retrying...")
		}
	}

	return "", fmt.Errorf("failed after %d attempts, last error: %w", cm.maxRetries, lastErr)
}

// AddUserMessage adds a user message to the conversation
func (cm *ConversationManager) AddUserMessage(content string) {
	cm.messages = append(cm.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	cm.chatLog.AddEntry("user", content, "conversation")
}

// AddAssistantMessage adds an assistant message to the conversation
func (cm *ConversationManager) AddAssistantMessage(content string) {
	cm.messages = append(cm.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})
	cm.chatLog.AddEntry("assistant", content, "conversation")
}

// GetMessages returns the current conversation messages
func (cm *ConversationManager) GetMessages() []openai.ChatCompletionMessage {
	return cm.messages
}

// RunConversation runs the interactive conversation loop with error recovery
func (cm *ConversationManager) RunConversation() error {
	// Start with initial message
	cm.AddUserMessage("Help me onboard")

	assistantMessage, err := cm.GetResponse()
	if err != nil {
		cm.AddError(err)
		fmt.Printf("❌ Failed to start conversation: %v\n", err)
		fmt.Println("🔄 Attempting to restart...")

		// Try to restart the conversation
		assistantMessage, err = cm.GetResponse()
		if err != nil {
			return fmt.Errorf("failed to restart conversation: %w", err)
		}
	}

	fmt.Printf("📻: %s\n\n", assistantMessage)
	cm.AddAssistantMessage(assistantMessage)

	// Main conversation loop
	for {
		fmt.Print("👤: ")

		if !cm.scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(cm.scanner.Text())

		if strings.ToLower(userInput) == "quit" {
			fmt.Println("👋 Thanks for using RepoRad.io! Goodbye!")
			break
		}

		if userInput == "" {
			continue
		}

		cm.AddUserMessage(userInput)

		assistantMessage, err := cm.GetResponse()
		if err != nil {
			cm.AddError(err)
			fmt.Printf("❌ Error in conversation: %v\n", err)
			fmt.Println("🔄 Attempting to continue with error context...")

			// Try to continue the conversation with error context
			assistantMessage, err = cm.GetResponse()
			if err != nil {
				fmt.Printf("❌ Failed to recover from error: %v\n", err)
				continue
			}
		}

		if strings.Contains(strings.ToUpper(assistantMessage), "SETUP COMPLETE") {
			fmt.Printf("📻: %s\n\n", assistantMessage)
			fmt.Println("✅ Configuration completed!")
			cm.chatLog.AddEntry("assistant", assistantMessage, "completion")
			break
		}

		fmt.Printf("📻: %s\n\n", assistantMessage)
		cm.AddAssistantMessage(assistantMessage)

	}

	return nil
}

var createCmd = &cobra.Command{
	Use:   "create <podcast-name>",
	Short: "Create podcast content from your codebase",
	Long: `The create command helps you generate podcast content from your repository.
It guides you through the process of creating episodes based on your codebase.`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func runCreate(cmd *cobra.Command, args []string) error {
	podcastName := args[0]

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("please set the OPENAI_API_KEY environment variable")
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Create chat log
	chatLog := NewChatLog(podcastName)

	// Get system prompt
	systemPrompt := getSystemPrompt()

	// Create conversation manager
	conversationManager := NewConversationManager(client, scanner, chatLog, systemPrompt)

	fmt.Printf("{📻} RepoRad.io Onboarding Assistant - %s\n", podcastName)
	fmt.Println("=====================================")
	fmt.Println()

	// Run the conversation with error recovery
	if err := conversationManager.RunConversation(); err != nil {
		return fmt.Errorf("conversation error: %w", err)
	}

	// Complete the chat log
	chatLog.Complete()

	// Extract project information with error recovery
	series, err := extractProjectInformation(client, conversationManager.GetMessages())
	if err != nil {
		fmt.Println("🔄 Retrying project information extraction...")

		// Retry extraction with error context
		series, err = extractProjectInformation(client, conversationManager.GetMessages())
		if err != nil {
			return fmt.Errorf("failed to extract project information: %w", err)
		}
	}

	// Create project structure
	if err := createProjectStructure(podcastName, series, chatLog); err != nil {
		return fmt.Errorf("error creating project structure: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func init() {
	// Add flags for the create command
	// TODO: Add specific flags as needed
}

// loadPromptFromFile reads the system prompt from a file
func loadPromptFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Remove the HTML comment at the beginning if present
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Skip the first line if it's an HTML comment
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "<!--") {
		lines = lines[1:]
	}

	return strings.Join(lines, "\n"), nil
}

// loadReadmeContent reads and returns the content of README.md
func loadReadmeContent() (string, error) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// getSystemPrompt returns the system prompt, either from file or default
func getSystemPrompt() string {
	var prompt string

	// Try to load from PROMPT.md file first
	if p, err := loadPromptFromFile("PROMPT.md"); err == nil && strings.TrimSpace(p) != "" {
		prompt = strings.TrimSpace(p)
	} else {
		// Fallback to default prompt
		prompt = "You are a helpful assistant."
	}

	// Add README content to provide context about the project
	if readmeContent, err := loadReadmeContent(); err == nil {
		prompt += "\n\nProject README:\n" + readmeContent
	}

	return prompt
}

// getChatResponse sends a message to OpenAI and returns the response
func getChatResponse(client *openai.Client, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")
}

// extractProjectInformation extracts project information from conversation messages
func extractProjectInformation(client *openai.Client, messages []openai.ChatCompletionMessage) (*Series, error) {
	var series Series
	schema, err := jsonschema.GenerateSchemaForType(series)
	if err != nil {
		return nil, fmt.Errorf("schema generation error: %w", err)
	}

	extractionResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "Based on our conversation, extract the project information including name, description, tags, and author.",
		}),
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "extraction",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error extracting project data: %w", err)
	}

	// Check that Choices is not empty array
	if len(extractionResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI for project extraction")
	}

	err = schema.Unmarshal(extractionResp.Choices[0].Message.Content, &series)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling project data: %w", err)
	}

	return &series, nil
}

func createProjectStructure(projectName string, series *Series, chat *ChatLog) error {
	// Create .reporadio directory
	reporadioDir := ".reporadio"
	if err := os.MkdirAll(reporadioDir, 0755); err != nil {
		return err
	}

	// Create project directory
	projectDir := filepath.Join(reporadioDir, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}

	// Create project episodes directory
	episodesDir := filepath.Join(projectDir, "episodes")
	if err := os.MkdirAll(episodesDir, 0755); err != nil {
		return err
	}

	// Save episode data as YAML
	seriesBytes, err := series.ToYAML()
	if err != nil {
		return err
	}

	episodePath := filepath.Join(projectDir, "episode.yaml")
	if err := os.WriteFile(episodePath, seriesBytes, 0644); err != nil {
		return err
	}

	// Save chat log
	chatBytes, err := chat.ToYAML()
	if err != nil {
		return err
	}

	chatPath := filepath.Join(projectDir, "chat.yaml")
	if err := os.WriteFile(chatPath, chatBytes, 0644); err != nil {
		return err
	}

	return nil
}
