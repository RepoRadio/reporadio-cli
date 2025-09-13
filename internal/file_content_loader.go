package internal

import (
	"os"
)

// FileContentLoaderImpl implements the FileContentLoader interface
type FileContentLoaderImpl struct{}

// NewFileContentLoader creates a new file content loader
func NewFileContentLoader() *FileContentLoaderImpl {
	return &FileContentLoaderImpl{}
}

// LoadFileContent loads content from files matching scan result
func (f *FileContentLoaderImpl) LoadFileContent(scanResult *ScanResult) (map[string]string, error) {
	if scanResult == nil {
		return nil, ErrConfigNotFound // Reuse existing error or create new one
	}

	content := make(map[string]string)

	for _, fileInfo := range scanResult.Files {
		fileContent, err := f.LoadSingleFile(fileInfo.Path)
		if err != nil {
			// Handle missing files gracefully - skip them rather than failing entirely
			Debugf(nil, "Skipping file %s: %v", fileInfo.Path, err)
			continue
		}

		content[fileInfo.Path] = fileContent
	}

	return content, nil
}

// LoadSingleFile loads content from a single file
func (f *FileContentLoaderImpl) LoadSingleFile(path string) (string, error) {
	if path == "" {
		return "", ErrConfigNotFound
	}

	// Check if path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "", ErrConfigNotFound // Return error for directories
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
