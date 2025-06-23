package internal

import "gopkg.in/yaml.v3"

// Episode represents the structure for episode content that guides users
// through structured instructions, such as introducing a new project.
type Episode struct {
	// Title is the title of the episode or topic
	Title string `json:"title" yaml:"title"`

	// Description is a high-level summary of what the episode covers
	Description string `json:"description" yaml:"description"`

	// Instructions contains special instructions to use for this episode
	Instructions string `json:"instructions" yaml:"instructions"`

	// Voicing defines the intended tone, style, or personality of the instruction content
	Voicing string `json:"voicing" yaml:"voicing"`

	// Include contains the file paths associated with this episode
	Include []string `json:"include" yaml:"include"`
}

// ToYAML converts episode to YAML bytes
func (e *Episode) ToYAML() ([]byte, error) {
	return yaml.Marshal(e)
}
