package internal

// ContextBuilderImpl implements the ContextBuilder interface
type ContextBuilderImpl struct {
	commandExecutor   *CommandExecutorImpl
	fileContentLoader *FileContentLoaderImpl
	scanner           *Scanner
}

// NewContextBuilder creates a new context builder
func NewContextBuilder() *ContextBuilderImpl {
	return &ContextBuilderImpl{
		commandExecutor:   NewCommandExecutor(),
		fileContentLoader: NewFileContentLoader(),
		scanner:           NewScanner(),
	}
}

// BuildCreateContext combines command output and file content into context
func (c *ContextBuilderImpl) BuildCreateContext(config *CreateConfig, mode string) (*CreateContext, error) {
	if config == nil {
		return nil, ErrConfigNotFound
	}

	ctx := &CreateContext{
		CommandOutput: make(map[string]string),
		FileContent:   make(map[string]string),
		Config:        config,
		Mode:          mode,
	}

	// Execute commands first
	if len(config.Commands) > 0 {
		commandOutput, err := c.commandExecutor.ExecuteCommands(config.Commands)
		if err != nil {
			return nil, err
		}
		ctx.CommandOutput = commandOutput
	}

	// Scan files based on include patterns
	scanResult, err := c.scanner.ScanWithIncludes(".", config.Include)

	if err != nil {
		return nil, err
	}

	// Load file content
	fileContent, err := c.fileContentLoader.LoadFileContent(scanResult)
	if err != nil {
		return nil, err
	}

	ctx.FileContent = fileContent

	return ctx, nil
}

// GetDefaultConfig returns the default configuration when create.yml missing
func (c *ContextBuilderImpl) GetDefaultConfig() *CreateConfig {
	return &CreateConfig{
		Include:  append([]string{}, DefaultInclude...),
		Commands: append([]string{}, DefaultCommands...),
	}
}
