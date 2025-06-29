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
	"gopkg.in/yaml.v3"
)

const RepoRadioDir = ".reporadio"

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
	if err := createProjectStructure(podcastName, series, chatLog, client, conversationManager.GetMessages()); err != nil {
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

// loadReadmeContent reads and returns the content of README.md
func loadReadmeContent() (string, error) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// getSystemPrompt returns the system prompt using the template system
func getSystemPrompt() string {
	promptManager, err := NewPromptManager()
	if err != nil {
		// Fallback to default prompt if template loading fails
		return "You are a helpful assistant."
	}

	// Prepare template data
	data := struct {
		ReadmeContent string
	}{}

	// Load README content if available
	if readmeContent, err := loadReadmeContent(); err == nil {
		data.ReadmeContent = readmeContent
	}

	prompt, err := promptManager.Execute("system_prompt.tmpl", data)
	if err != nil {
		// Fallback to default prompt if template execution fails
		return "You are a helpful assistant."
	}

	return prompt
}

// getChatResponse sends a message to OpenAI and returns the response
func getChatResponse(client *openai.Client, messages []openai.ChatCompletionMessage) (string, error) {
	request := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: messages,
	}

	DebugOpenAIRequest(nil, messages, string(openai.GPT4o))

	resp, err := client.CreateChatCompletion(context.Background(), request)

	if err != nil {
		Debugf(nil, "OpenAI API error: %v", err)
		return "", err
	}

	DebugOpenAIResponse(nil, resp)

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	Debug(nil, "No response choices returned from OpenAI")
	return "", fmt.Errorf("no response from OpenAI")
}

// extractProjectInformation extracts project information from conversation messages
func extractProjectInformation(client *openai.Client, messages []openai.ChatCompletionMessage) (*Series, error) {
	var series Series
	schema, err := jsonschema.GenerateSchemaForType(series)
	if err != nil {
		return nil, fmt.Errorf("schema generation error: %w", err)
	}

	// Get extraction prompt from template
	promptManager, err := NewPromptManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt manager: %w", err)
	}

	extractionPrompt, err := promptManager.Execute("extract_project_info.tmpl", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	extractionResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: extractionPrompt,
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

// extractEpisodes extracts episodes from conversation and repository context
func extractEpisodes(client *openai.Client, messages []openai.ChatCompletionMessage, scanResult *ScanResult) ([]Episode, error) {
	// Create context message with directory structure
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Repository Structure:\n")
	contextBuilder.WriteString(fmt.Sprintf("Project Type: %s\n", scanResult.ProjectType))

	if scanResult.ReadmePath != "" {
		contextBuilder.WriteString(fmt.Sprintf("README: %s\n", scanResult.ReadmePath))
	}

	contextBuilder.WriteString("\nFiles:\n")
	for _, file := range scanResult.Files {
		contextBuilder.WriteString(fmt.Sprintf("- %s (%s, %d bytes)\n", file.Path, file.Extension, file.Size))
	}

	// Get extraction prompt from template
	promptManager, err := NewPromptManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt manager: %w", err)
	}

	templateData := struct {
		DirectoryStructure  string
		ConversationContext string
	}{
		DirectoryStructure:  contextBuilder.String(),
		ConversationContext: "Previous conversation context included above",
	}

	extractionPrompt, err := promptManager.Execute("extract_episodes.tmpl", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Create episodes wrapper for schema
	episodesWrapper := struct {
		Episodes []Episode `json:"episodes"`
	}{}
	schema, err := jsonschema.GenerateSchemaForType(episodesWrapper)
	if err != nil {
		return nil, fmt.Errorf("schema generation error: %w", err)
	}

	// Create extraction request with repository context
	extractionMessages := append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: extractionPrompt,
	})

	extractionResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: extractionMessages,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "episodes_extraction",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error extracting episodes: %w", err)
	}

	if len(extractionResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI for episodes extraction")
	}

	var result struct {
		Episodes []Episode `json:"episodes"`
	}
	err = schema.Unmarshal(extractionResp.Choices[0].Message.Content, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling episodes: %w", err)
	}

	return result.Episodes, nil
}

func createProjectStructure(projectName string, series *Series, chat *ChatLog, client *openai.Client, messages []openai.ChatCompletionMessage) error {
	Debug(nil, "Starting project structure creation")
	Debugf(nil, "Project name: %s", projectName)

	// Create .reporadio directory
	if err := os.MkdirAll(RepoRadioDir, 0755); err != nil {
		return err
	}
	Debug(nil, "Created .reporadio directory")

	// Create project directory
	projectDir := filepath.Join(RepoRadioDir, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}
	Debugf(nil, "Created project directory: %s", projectDir)

	// Scan repository to get directory structure
	Debug(nil, "Starting repository scan")
	repoScanner := NewScanner()
	scanResult, err := repoScanner.ScanRepository(".")
	if err != nil {
		return fmt.Errorf("failed to scan repository: %w", err)
	}
	Debugf(nil, "Repository scan complete - found %d files, project type: %s", len(scanResult.Files), scanResult.ProjectType)

	// Extract episodes using AI with repository context
	Debug(nil, "Starting episode extraction with OpenAI")
	episodes, err := extractEpisodes(client, messages, scanResult)
	if err != nil {
		return fmt.Errorf("failed to extract episodes: %w", err)
	}
	Debugf(nil, "Episode extraction complete - generated %d episodes", len(episodes))

	// Create consolidated podcast config with series info and episodes
	podcastConfig := struct {
		Title        string     `yaml:"title"`
		Description  string     `yaml:"description"`
		Instructions string     `yaml:"instructions"`
		Voicing      string     `yaml:"voicing"`
		Type         SeriesType `yaml:"type"`
		Episodes     []Episode  `yaml:"episodes"`
	}{
		Title:        series.Title,
		Description:  series.Description,
		Instructions: series.Instructions,
		Voicing:      series.Voicing,
		Type:         series.Type,
		Episodes:     episodes,
	}

	podcastYAML, err := yaml.Marshal(podcastConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal podcast config to YAML: %w", err)
	}

	podcastPath := filepath.Join(projectDir, "podcast.yml")
	if err := os.WriteFile(podcastPath, podcastYAML, 0644); err != nil {
		return fmt.Errorf("failed to write podcast.yml: %w", err)
	}
	Debugf(nil, "Saved podcast config to: %s", podcastPath)

	fmt.Printf("✅ Generated %d episodes based on conversation and repository analysis\n", len(episodes))

	// Save chat log
	chatBytes, err := chat.ToYAML()
	if err != nil {
		return err
	}

	chatPath := filepath.Join(projectDir, "chat.yaml")
	if err := os.WriteFile(chatPath, chatBytes, 0644); err != nil {
		return err
	}
	Debugf(nil, "Saved chat log to: %s", chatPath)

	Debug(nil, "Project structure creation complete")
	return nil
}
