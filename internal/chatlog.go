package internal

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// ChatEntry represents a single interaction in the chat log
type ChatEntry struct {
	Timestamp time.Time `yaml:"timestamp" json:"timestamp"`
	Role      string    `yaml:"role" json:"role"` // "user" or "assistant"
	Message   string    `yaml:"message" json:"message"`
	Step      string    `yaml:"step,omitempty" json:"step,omitempty"` // e.g., "project_name", "metadata", etc.
}

// ChatLog represents the complete conversation history
type ChatLog struct {
	ProjectName string      `yaml:"project_name" json:"project_name"`
	StartTime   time.Time   `yaml:"start_time" json:"start_time"`
	EndTime     *time.Time  `yaml:"end_time,omitempty" json:"end_time,omitempty"`
	Entries     []ChatEntry `yaml:"entries" json:"entries"`
}

// NewChatLog creates a new chat log for a project
func NewChatLog(projectName string) *ChatLog {
	return &ChatLog{
		ProjectName: projectName,
		StartTime:   time.Now(),
		Entries:     make([]ChatEntry, 0),
	}
}

// AddEntry adds a new entry to the chat log
func (cl *ChatLog) AddEntry(role, message, step string) {
	entry := ChatEntry{
		Timestamp: time.Now(),
		Role:      role,
		Message:   message,
		Step:      step,
	}
	cl.Entries = append(cl.Entries, entry)
}

// Complete marks the chat log as completed
func (cl *ChatLog) Complete() {
	now := time.Now()
	cl.EndTime = &now
}

// ToYAML converts chat log to YAML bytes
func (cl *ChatLog) ToYAML() ([]byte, error) {
	return yaml.Marshal(cl)
}

// FromYAML creates chat log from YAML bytes
func FromYAML(data []byte) (*ChatLog, error) {
	var cl ChatLog
	if err := yaml.Unmarshal(data, &cl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return &cl, nil
}

// Duration returns the duration of the chat session
func (cl *ChatLog) Duration() time.Duration {
	endTime := time.Now()
	if cl.EndTime != nil {
		endTime = *cl.EndTime
	}
	return endTime.Sub(cl.StartTime)
}
