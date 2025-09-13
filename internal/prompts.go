package internal

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed prompts/*.tmpl
var promptFS embed.FS

type PromptManager struct {
	templates map[string]*template.Template
}

func NewPromptManager() (*PromptManager, error) {
	pm := &PromptManager{
		templates: make(map[string]*template.Template),
	}

	// Check if custom prompts directory exists
	customPromptsDir := filepath.Join(".reporadio", "prompts")
	useCustomPrompts := false
	if _, err := os.Stat(customPromptsDir); err == nil {
		useCustomPrompts = true
	}

	// Get list of embedded templates first
	entries, err := promptFS.ReadDir("prompts")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		var content []byte

		// Try to load from custom directory first if it exists
		if useCustomPrompts {
			customPath := filepath.Join(customPromptsDir, name)
			if customContent, err := os.ReadFile(customPath); err == nil {
				content = customContent
			} else {
				// Fallback to embedded template
				content, err = promptFS.ReadFile("prompts/" + name)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// Use embedded template
			content, err = promptFS.ReadFile("prompts/" + name)
			if err != nil {
				return nil, err
			}
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return nil, err
		}

		pm.templates[name] = tmpl
	}

	return pm, nil
}

func (pm *PromptManager) Execute(templateName string, data interface{}) (string, error) {
	tmpl, exists := pm.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
