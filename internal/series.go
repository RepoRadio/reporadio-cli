package internal

import "gopkg.in/yaml.v3"

// SeriesType represents the type of series content
type SeriesType string

const (
	SeriesTypeOnboarding SeriesType = "onboarding"
	SeriesTypeChangelog  SeriesType = "changelog"
)

// Series represents onboarding content that guides users through structured instructions
type Series struct {
	Title        string     `json:"title" yaml:"title"`
	Description  string     `json:"description" yaml:"description"`
	Instructions string     `json:"instructions" yaml:"instructions"`
	Voicing      string     `json:"voicing" yaml:"voicing"`
	Type         SeriesType `json:"type" yaml:"type"`
}

// ToYAML converts episode to YAML bytes
func (e *Series) ToYAML() ([]byte, error) {
	return yaml.Marshal(e)
}
