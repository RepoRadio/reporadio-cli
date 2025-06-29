package internal

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEpisodeWithCommands(t *testing.T) {
	// Test episode with commands field
	episode := Episode{
		Title:        "Test Episode",
		Description:  "A test episode with commands",
		Instructions: "Test instructions",
		Voicing:      "friendly",
		Include:      []string{"README.md"},
		Commands:     []string{"echo 'hello'", "date", "ls -la"},
	}

	// Test YAML marshaling
	yamlBytes, err := episode.ToYAML()
	if err != nil {
		t.Fatalf("Failed to marshal episode to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Episode
	err = yaml.Unmarshal(yamlBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal episode from YAML: %v", err)
	}

	// Verify all fields including commands
	if unmarshaled.Title != episode.Title {
		t.Errorf("Expected title %q, got %q", episode.Title, unmarshaled.Title)
	}

	if len(unmarshaled.Commands) != len(episode.Commands) {
		t.Errorf("Expected %d commands, got %d", len(episode.Commands), len(unmarshaled.Commands))
	}

	for i, cmd := range episode.Commands {
		if unmarshaled.Commands[i] != cmd {
			t.Errorf("Expected command %q, got %q", cmd, unmarshaled.Commands[i])
		}
	}
}

func TestEpisodeWithoutCommands(t *testing.T) {
	// Test episode without commands field (should still work)
	episode := Episode{
		Title:        "Test Episode",
		Description:  "A test episode without commands",
		Instructions: "Test instructions",
		Voicing:      "friendly",
		Include:      []string{"README.md"},
		// Commands field omitted
	}

	// Test YAML marshaling
	yamlBytes, err := episode.ToYAML()
	if err != nil {
		t.Fatalf("Failed to marshal episode to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Episode
	err = yaml.Unmarshal(yamlBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal episode from YAML: %v", err)
	}

	// Commands should be nil/empty
	if len(unmarshaled.Commands) != 0 {
		t.Errorf("Expected no commands, got %d", len(unmarshaled.Commands))
	}
}

func TestEpisodeJSONSerialization(t *testing.T) {
	// Test JSON serialization with commands
	episode := Episode{
		Title:        "Test Episode",
		Description:  "A test episode with commands",
		Instructions: "Test instructions",
		Voicing:      "friendly",
		Include:      []string{"README.md"},
		Commands:     []string{"echo 'test'"},
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(episode)
	if err != nil {
		t.Fatalf("Failed to marshal episode to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Episode
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal episode from JSON: %v", err)
	}

	// Verify commands field
	if len(unmarshaled.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(unmarshaled.Commands))
	}

	if unmarshaled.Commands[0] != "echo 'test'" {
		t.Errorf("Expected command %q, got %q", "echo 'test'", unmarshaled.Commands[0])
	}
}
