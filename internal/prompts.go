package internal

import (
	"bytes"
	"embed"
	"fmt"
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
	
	entries, err := promptFS.ReadDir("prompts")
	if err != nil {
		return nil, err
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		content, err := promptFS.ReadFile("prompts/" + name)
		if err != nil {
			return nil, err
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
