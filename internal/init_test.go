package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test the init command
	err := runInit(nil, []string{})
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify .reporadio directory was created
	repoRadioPath := filepath.Join(tmpDir, RepoRadioDir)
	if _, err := os.Stat(repoRadioPath); os.IsNotExist(err) {
		t.Error(".reporadio directory was not created")
	}

	// Verify .reporadio/prompts directory was created
	promptsPath := filepath.Join(repoRadioPath, "prompts")
	if _, err := os.Stat(promptsPath); os.IsNotExist(err) {
		t.Error(".reporadio/prompts directory was not created")
	}

	// Verify all template files were copied
	expectedTemplates := []string{
		"system_prompt.tmpl",
		"extract_episodes.tmpl",
		"extract_project_info.tmpl",
		"episode_transcript.tmpl",
	}

	for _, templateName := range expectedTemplates {
		templatePath := filepath.Join(promptsPath, templateName)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Errorf("Template %s was not copied to .reporadio/prompts/", templateName)
		}

		// Verify the file is not empty
		info, err := os.Stat(templatePath)
		if err != nil {
			t.Errorf("Cannot stat copied template %s: %v", templateName, err)
		} else if info.Size() == 0 {
			t.Errorf("Copied template %s is empty", templateName)
		}
	}
}

func TestRunInitAlreadyExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Pre-create the .reporadio directory
	repoRadioPath := filepath.Join(tmpDir, RepoRadioDir)
	err := os.MkdirAll(repoRadioPath, 0755)
	if err != nil {
		t.Fatalf("Failed to pre-create .reporadio directory: %v", err)
	}

	// Test the init command when directory already exists
	err = runInit(nil, []string{})
	if err != nil {
		t.Fatalf("runInit failed when .reporadio already exists: %v", err)
	}

	// Verify prompts directory was still created
	promptsPath := filepath.Join(repoRadioPath, "prompts")
	if _, err := os.Stat(promptsPath); os.IsNotExist(err) {
		t.Error(".reporadio/prompts directory was not created when .reporadio already existed")
	}
}

func TestRunInitWithExistingPrompts(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Pre-create the directory structure and a custom template
	repoRadioPath := filepath.Join(tmpDir, RepoRadioDir)
	promptsPath := filepath.Join(repoRadioPath, "prompts")
	err := os.MkdirAll(promptsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to pre-create directory structure: %v", err)
	}

	// Create a custom template file
	customTemplatePath := filepath.Join(promptsPath, "system_prompt.tmpl")
	customContent := "This is a custom system prompt template"
	err = os.WriteFile(customTemplatePath, []byte(customContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom template: %v", err)
	}

	// Test the init command
	err = runInit(nil, []string{})
	if err != nil {
		t.Fatalf("runInit failed when prompts directory already exists: %v", err)
	}

	// Verify the custom template was preserved (not overwritten)
	content, err := os.ReadFile(customTemplatePath)
	if err != nil {
		t.Fatalf("Failed to read template after init: %v", err)
	}

	if string(content) != customContent {
		t.Error("Custom template was overwritten - should preserve existing files")
	}

	// Verify other templates were still copied
	otherTemplates := []string{
		"extract_episodes.tmpl",
		"extract_project_info.tmpl", 
		"episode_transcript.tmpl",
	}

	for _, templateName := range otherTemplates {
		templatePath := filepath.Join(promptsPath, templateName)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Errorf("Template %s was not copied when other templates existed", templateName)
		}
	}
}

func TestCopyEmbeddedTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	promptsPath := filepath.Join(tmpDir, "prompts")

	// Test the copyEmbeddedTemplates function
	err := copyEmbeddedTemplates(promptsPath)
	if err != nil {
		t.Fatalf("copyEmbeddedTemplates failed: %v", err)
	}

	// Verify all templates were copied
	expectedTemplates := []string{
		"system_prompt.tmpl",
		"extract_episodes.tmpl", 
		"extract_project_info.tmpl",
		"episode_transcript.tmpl",
	}

	for _, templateName := range expectedTemplates {
		templatePath := filepath.Join(promptsPath, templateName)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Errorf("Template %s was not copied", templateName)
		}

		// Verify content matches embedded template
		embeddedContent, err := promptFS.ReadFile("prompts/" + templateName)
		if err != nil {
			t.Errorf("Failed to read embedded template %s: %v", templateName, err)
			continue
		}

		copiedContent, err := os.ReadFile(templatePath)
		if err != nil {
			t.Errorf("Failed to read copied template %s: %v", templateName, err)
			continue
		}

		if string(embeddedContent) != string(copiedContent) {
			t.Errorf("Template %s content does not match embedded version", templateName)
		}
	}
}

func TestCopyEmbeddedTemplatesTargetDirNotExist(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	promptsPath := filepath.Join(tmpDir, "nonexistent", "prompts")

	// Test copyEmbeddedTemplates with non-existent directory
	err := copyEmbeddedTemplates(promptsPath)
	if err != nil {
		t.Fatalf("copyEmbeddedTemplates should create target directory: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(promptsPath); os.IsNotExist(err) {
		t.Error("Target directory was not created")
	}

	// Verify templates were copied
	templatePath := filepath.Join(promptsPath, "system_prompt.tmpl")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Error("Template was not copied to newly created directory")
	}
}

func TestCopyEmbeddedTemplatesPreserveExisting(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	promptsPath := filepath.Join(tmpDir, "prompts")
	err := os.MkdirAll(promptsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create a custom template that should be preserved
	customTemplatePath := filepath.Join(promptsPath, "system_prompt.tmpl")
	customContent := "This is my custom system prompt"
	err = os.WriteFile(customTemplatePath, []byte(customContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom template: %v", err)
	}

	// Test copyEmbeddedTemplates
	err = copyEmbeddedTemplates(promptsPath)
	if err != nil {
		t.Fatalf("copyEmbeddedTemplates failed: %v", err)
	}

	// Verify custom template was preserved
	content, err := os.ReadFile(customTemplatePath)
	if err != nil {
		t.Fatalf("Failed to read custom template: %v", err)
	}

	if string(content) != customContent {
		t.Error("Custom template should be preserved, not overwritten")
	}

	// Verify other templates were still copied
	otherTemplate := filepath.Join(promptsPath, "extract_episodes.tmpl")
	if _, err := os.Stat(otherTemplate); os.IsNotExist(err) {
		t.Error("Other templates should still be copied when some exist")
	}
}