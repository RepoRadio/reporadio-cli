package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnore functionality
//
// The scanner now supports .gitignore files to filter out unwanted files and directories
// during repository scanning. This helps exclude build artifacts, dependencies, logs,
// and other files that shouldn't be included in podcast episodes.
//
// Supported gitignore features:
// - Basic file patterns (*.log, temp.txt)
// - Directory patterns (node_modules/, build/)
// - Negation patterns (!important.log)
// - Glob patterns with * and ? wildcards
// - Leading slash patterns (/build)
// - Comments and empty lines (ignored)
// - Nested gitignore files in subdirectories
//
// The implementation follows standard gitignore semantics where:
// - Patterns are processed in order
// - Later patterns can override earlier ones
// - Negated patterns (!pattern) exclude files from being ignored
// - Directory patterns (ending with /) only match directories
// - Patterns without path separators match filenames in any directory

// GitIgnorePattern represents a single pattern from a gitignore file
type GitIgnorePattern struct {
	Pattern   string
	IsDir     bool // ends with /
	IsNegated bool // starts with !
	IsGlob    bool // contains * or ?
}

// GitIgnore holds patterns and provides matching functionality
type GitIgnore struct {
	Patterns []GitIgnorePattern
	BasePath string
}

// Scanner handles repository scanning and content analysis
type Scanner struct{}

// ScanResult contains the results of a repository scan
type ScanResult struct {
	ReadmePath  string
	Files       []FileInfo
	ProjectType string
}

// FileInfo represents information about a discovered file
type FileInfo struct {
	Path      string
	Extension string
	Size      int64
}

// ReadmeStructure represents the parsed structure of a README file
type ReadmeStructure struct {
	Title    string
	Sections []Section
}

// Section represents a section in the README
type Section struct {
	Title   string
	Content string
	Level   int
}

// NewScanner creates a new repository scanner
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanRepository scans a repository directory and returns information about its structure
func (s *Scanner) ScanRepository(rootPath string) (*ScanResult, error) {
	result := &ScanResult{
		Files: make([]FileInfo, 0),
	}

	// Load gitignore and reporadioignore files
	gitignores := make(map[string]*GitIgnore)
	reporadioignores := make(map[string]*GitIgnore)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Load gitignore files as we encounter them
		if !info.IsDir() && info.Name() == ".gitignore" {
			if gitignore, err := loadGitIgnoreFile(path); err == nil {
				gitignores[filepath.Dir(path)] = gitignore
			}
		}

		// Load reporadioignore files as we encounter them
		if !info.IsDir() && info.Name() == ".reporadioignore" {
			if reporadioignore, err := loadGitIgnoreFile(path); err == nil {
				reporadioignores[filepath.Dir(path)] = reporadioignore
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load ignore files: %w", err)
	}

	// Now scan files, respecting gitignore and reporadioignore rules
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if ignored by either gitignore or reporadioignore
		if s.shouldIgnorePath(path, info.IsDir(), gitignores, reporadioignores) {
			if info.IsDir() {
				return filepath.SkipDir // Skip entire directory
			}
			return nil // Skip file
		}

		// Skip directories and hidden files (including .gitignore and .reporadioignore)
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Check for README
		if strings.ToLower(info.Name()) == "readme.md" {
			result.ReadmePath = path
		}

		// Add file info
		fileInfo := FileInfo{
			Path:      path,
			Extension: filepath.Ext(path),
			Size:      info.Size(),
		}
		result.Files = append(result.Files, fileInfo)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Detect project type
	result.ProjectType = s.detectProjectType(result.Files)

	return result, nil
}

// shouldIgnorePath checks if a path should be ignored based on all applicable gitignore and reporadioignore files
func (s *Scanner) shouldIgnorePath(path string, isDir bool, gitignores map[string]*GitIgnore, reporadioignores map[string]*GitIgnore) bool {
	// Check all gitignore files that could apply to this path
	for basePath, gitignore := range gitignores {
		// Only check gitignore files that are in parent directories or the same directory
		relBase, err := filepath.Rel(basePath, path)
		if err != nil || strings.HasPrefix(relBase, "..") {
			continue // This gitignore doesn't apply to this path
		}

		if gitignore.ShouldIgnore(path, isDir) {
			return true
		}
	}

	// Check all reporadioignore files that could apply to this path
	for basePath, reporadioignore := range reporadioignores {
		// Only check reporadioignore files that are in parent directories or the same directory
		relBase, err := filepath.Rel(basePath, path)
		if err != nil || strings.HasPrefix(relBase, "..") {
			continue // This reporadioignore doesn't apply to this path
		}

		if reporadioignore.ShouldIgnore(path, isDir) {
			return true
		}
	}

	return false
}

// GetFilesByExtension returns files with the specified extension
func (sr *ScanResult) GetFilesByExtension(ext string) []FileInfo {
	var files []FileInfo
	for _, file := range sr.Files {
		if file.Extension == ext {
			files = append(files, file)
		}
	}
	return files
}

// ParseReadme parses README content and extracts structure
func (s *Scanner) ParseReadme(content []byte) *ReadmeStructure {
	structure := &ReadmeStructure{
		Sections: make([]Section, 0),
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	var currentSection *Section

	for scanner.Scan() {
		line := scanner.Text()

		// Check for headers
		if strings.HasPrefix(line, "#") {
			// Save previous section
			if currentSection != nil {
				structure.Sections = append(structure.Sections, *currentSection)
			}

			level := 0
			for _, char := range line {
				if char == '#' {
					level++
				} else {
					break
				}
			}

			title := strings.TrimSpace(line[level:])

			// Set main title if this is the first H1
			if level == 1 && structure.Title == "" {
				structure.Title = title
			}

			currentSection = &Section{
				Title:   title,
				Level:   level,
				Content: "",
			}
		} else if currentSection != nil {
			// Add content to current section
			currentSection.Content += line + "\n"
		}
	}

	// Add last section
	if currentSection != nil {
		structure.Sections = append(structure.Sections, *currentSection)
	}

	return structure
}

// GenerateEpisodes automatically generates episode suggestions based on repository content
func (s *Scanner) GenerateEpisodes(result *ScanResult) []Episode {
	var episodes []Episode

	// Episode 1: Project Overview (always include README)
	if result.ReadmePath != "" {
		episodes = append(episodes, Episode{
			Title:        "Project Overview",
			Description:  "An introduction to the project, its purpose, and getting started guide",
			Instructions: "Focus on the README content and provide an overview of what this project does",
			Voicing:      "Friendly and welcoming introduction",
			Include:      []string{"README.md"},
		})
	}

	// Episode 2: Code Structure (if code files exist)
	codeFiles := s.getCodeFiles(result.Files)
	if len(codeFiles) > 0 {
		var paths []string
		for _, file := range codeFiles {
			// Use relative paths from repository root
			relPath, _ := filepath.Rel(filepath.Dir(result.ReadmePath), file.Path)
			paths = append(paths, relPath)
			if len(paths) >= 5 { // Limit to avoid too many paths
				break
			}
		}

		episodes = append(episodes, Episode{
			Title:        "Code Architecture and Structure",
			Description:  fmt.Sprintf("Deep dive into the codebase structure and main %s components", result.ProjectType),
			Instructions: "Analyze the code structure and explain the main components and their relationships",
			Voicing:      "Technical but accessible explanation",
			Include:      paths,
		})
	}

	// Episode 3: Documentation and Examples (if docs exist)
	docFiles := s.getDocumentationFiles(result.Files)
	if len(docFiles) > 0 {
		var paths []string
		for _, file := range docFiles {
			relPath, _ := filepath.Rel(filepath.Dir(result.ReadmePath), file.Path)
			paths = append(paths, relPath)
		}

		episodes = append(episodes, Episode{
			Title:        "Documentation and Examples",
			Description:  "Exploring the project documentation, usage examples, and guides",
			Instructions: "Walk through the documentation and highlight practical examples",
			Voicing:      "Educational and practical guidance",
			Include:      paths,
		})
	}

	return episodes
}

// detectProjectType detects the primary programming language/framework
func (s *Scanner) detectProjectType(files []FileInfo) string {
	counts := make(map[string]int)

	for _, file := range files {
		switch file.Extension {
		case ".go":
			counts["Go"]++
		case ".js", ".ts":
			counts["JavaScript/TypeScript"]++
		case ".py":
			counts["Python"]++
		case ".java":
			counts["Java"]++
		case ".rs":
			counts["Rust"]++
		case ".cpp", ".cc", ".cxx":
			counts["C++"]++
		case ".c":
			counts["C"]++
		}
	}

	// Find the most common type
	maxCount := 0
	projectType := "Mixed"
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			projectType = lang
		}
	}

	return projectType
}

// getCodeFiles returns files that contain source code
func (s *Scanner) getCodeFiles(files []FileInfo) []FileInfo {
	var codeFiles []FileInfo
	codeExtensions := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".py": true,
		".java": true, ".rs": true, ".cpp": true, ".c": true,
		".h": true, ".hpp": true, ".cs": true, ".php": true,
	}

	for _, file := range files {
		if codeExtensions[file.Extension] {
			codeFiles = append(codeFiles, file)
		}
	}

	return codeFiles
}

// getDocumentationFiles returns files that contain documentation
func (s *Scanner) getDocumentationFiles(files []FileInfo) []FileInfo {
	var docFiles []FileInfo

	for _, file := range files {
		if file.Extension == ".md" ||
			strings.Contains(strings.ToLower(file.Path), "doc") ||
			strings.Contains(strings.ToLower(file.Path), "example") {
			docFiles = append(docFiles, file)
		}
	}

	return docFiles
}

// parseGitIgnore parses gitignore content and returns a GitIgnore instance
func parseGitIgnore(content []byte, basePath string) *GitIgnore {
	gitignore := &GitIgnore{
		Patterns: make([]GitIgnorePattern, 0),
		BasePath: basePath,
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pattern := GitIgnorePattern{}

		// Check for negation
		if strings.HasPrefix(line, "!") {
			pattern.IsNegated = true
			line = line[1:]
		}

		// Handle leading slash (absolute path from repo root)
		line = strings.TrimPrefix(line, "/")

		// Check for directory pattern
		if strings.HasSuffix(line, "/") {
			pattern.IsDir = true
			line = strings.TrimSuffix(line, "/")
		}

		// Check for glob patterns
		if strings.Contains(line, "*") || strings.Contains(line, "?") {
			pattern.IsGlob = true
		}

		pattern.Pattern = line
		gitignore.Patterns = append(gitignore.Patterns, pattern)
	}

	return gitignore
}

// ShouldIgnore checks if a path should be ignored based on gitignore patterns
func (g *GitIgnore) ShouldIgnore(path string, isDir bool) bool {
	// Convert to relative path from base path
	relPath, err := filepath.Rel(g.BasePath, path)
	if err != nil {
		relPath = path
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)

	matched := false

	// Process patterns in order, later patterns can override earlier ones
	for _, pattern := range g.Patterns {
		if g.matchesPattern(pattern, relPath, isDir) {
			matched = !pattern.IsNegated
		}
	}

	return matched
}

// matchesPattern checks if a path matches a specific gitignore pattern
func (g *GitIgnore) matchesPattern(pattern GitIgnorePattern, path string, isDir bool) bool {
	// Directory patterns only match directories
	if pattern.IsDir && !isDir {
		return false
	}

	if pattern.IsGlob {
		// Handle different types of glob patterns
		if strings.Contains(pattern.Pattern, "/") {
			// Pattern contains path separator, match against full path
			matched, _ := filepath.Match(pattern.Pattern, path)
			return matched
		} else {
			// Pattern is just a filename pattern, check filename only
			matched, _ := filepath.Match(pattern.Pattern, filepath.Base(path))
			if matched {
				return true
			}
			// Also check if it matches any directory component
			parts := strings.Split(path, "/")
			for _, part := range parts {
				if matched, _ := filepath.Match(pattern.Pattern, part); matched {
					return true
				}
			}
			return false
		}
	} else {
		// Exact match cases
		if pattern.IsDir {
			// Directory pattern: match if path is the directory or under it
			return path == pattern.Pattern || strings.HasPrefix(path, pattern.Pattern+"/")
		} else {
			// File pattern: exact match or basename match
			if path == pattern.Pattern || filepath.Base(path) == pattern.Pattern {
				return true
			}
			// Also check if any directory component matches
			parts := strings.Split(path, "/")
			for _, part := range parts {
				if part == pattern.Pattern {
					return true
				}
			}
			return false
		}
	}
}

// loadGitIgnoreFile loads and parses a gitignore file
func loadGitIgnoreFile(path string) (*GitIgnore, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseGitIgnore(content, filepath.Dir(path)), nil
}

// ResolveIncludePaths resolves a list of file and directory paths into a list of files
// Files are included as explicit overrides (ignoring gitignore rules)
// Directories are expanded to include all non-ignored files within them
// Non-existent paths are skipped with a warning
func (s *Scanner) ResolveIncludePaths(includes []string) ([]string, error) {
	var resolvedPaths []string

	for _, includePath := range includes {
		info, err := os.Stat(includePath)
		if os.IsNotExist(err) {
			// Skip non-existent paths with warning
			fmt.Fprintf(os.Stderr, "Warning: Include path does not exist: %s\n", includePath)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to stat include path %s: %w", includePath, err)
		}

		if info.IsDir() {
			// Expand directory to include all non-ignored files
			dirFiles, err := s.expandDirectory(includePath)
			if err != nil {
				return nil, fmt.Errorf("failed to expand directory %s: %w", includePath, err)
			}
			resolvedPaths = append(resolvedPaths, dirFiles...)
		} else {
			// Include file directly (explicit override)
			resolvedPaths = append(resolvedPaths, includePath)
		}
	}

	return resolvedPaths, nil
}

// expandDirectory finds all non-ignored files within a directory
func (s *Scanner) expandDirectory(dirPath string) ([]string, error) {
	var files []string

	// Load gitignore and reporadioignore files for this directory tree
	gitignores := make(map[string]*GitIgnore)
	reporadioignores := make(map[string]*GitIgnore)

	// First pass: collect ignore files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip paths with errors
		}

		if info.Name() == ".gitignore" {
			if gitignore, err := loadGitIgnoreFile(path); err == nil {
				gitignores[filepath.Dir(path)] = gitignore
			}
		} else if info.Name() == ".reporadioignore" {
			if reporadioignore, err := loadGitIgnoreFile(path); err == nil {
				reporadioignores[filepath.Dir(path)] = reporadioignore
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Second pass: collect non-ignored files
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip paths with errors
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Check if path should be ignored
		if s.shouldIgnorePath(path, false, gitignores, reporadioignores) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
