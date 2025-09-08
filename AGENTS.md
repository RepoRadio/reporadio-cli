# RepoRadio CLI - Agent Development Guide

Turn your Git repository into a podcast. RepoRadio is an open source CLI tool that generates narrated audio episodes directly from your codebase.

## Project Details

- **Language**: Go 1.23.3
- **Type**: CLI Tool
- **Framework**: Cobra (CLI framework)
- **Dependencies**: OpenAI SDK, YAML parsing

## Development Environment

### Prerequisites

- Go 1.23.3 or later
- OpenAI API key (set as `OPENAI_API_KEY` environment variable)
- Git (for repository analysis)

### Setup Commands

```bash
# Clone and setup
git clone https://github.com/reporadio/reporadio-cli
cd reporadio-cli

# Install dependencies
go mod download

# Set up environment
export OPENAI_API_KEY=sk-...
```

### Package Manager

- **Tool**: Go modules
- **File**: `go.mod`

## Code Quality

### Formatting

- **Tool**: `gofmt`
- **Command**: `go fmt ./...`

### Type Checking

- **Tool**: Go compiler
- **Command**: `go build`

## Testing

### Framework

- **Tool**: Go built-in testing
- **Test Files**: `internal/*_test.go`

### Test Commands

- **Unit Tests**: `go test ./...`
- **Specific Package**: `go test ./internal`
- **Verbose**: `go test -v ./...`
- **Coverage**: `go test -cover ./...`

## Building

### Build Commands

- **Development**: `make run`
- **Production**: `make build`
- **Install**: `make install`
- **Clean**: `make clean`

### Output Directory

- **Location**: `bin/`
- **Binary Name**: `reporadio-cli`

## Conventions

### Directory Structure

```
├── cmd/                 # Command definitions
├── internal/           # Internal packages
├── .reporadio/        # Generated podcast content
├── bin/               # Compiled binaries
├── Makefile           # Build automation
└── main.go           # Application entry point
```

### Naming Conventions

- **Files**: snake_case for test files, camelCase for regular files
- **Functions**: PascalCase for exported, camelCase for internal
- **Packages**: lowercase, single word preferred

### Architectural Patterns

- CLI commands using Cobra framework
- Internal packages for core functionality
- Configuration-driven podcast generation
- OpenAI API integration for content generation

## Pull Request Guidelines

### Requirements

- All tests must pass: `go test ./...`
- Code must be formatted: `go fmt ./...`
- Build must succeed: `go build`

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification:

**Format**: `<type>[optional scope]: <description>`

**Common Types**:
- `feat:` - New feature (MINOR version)
- `fix:` - Bug fix (PATCH version)
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring without feature changes
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks, dependency updates
- `ci:` - CI/CD configuration changes
- `perf:` - Performance improvements

**Examples**:
```
feat(cli): add playlist generation command
fix(audio): resolve TTS generation timeout
docs: update installation instructions
test(scanner): add tests for file inclusion logic
```

**Breaking Changes**: Add `!` after type or include `BREAKING CHANGE:` footer

### Branch Naming

- Feature branches: `feat/description`
- Bug fixes: `fix/description`
- Documentation: `docs/description`

## Troubleshooting

### Common Issues

1. **Missing OpenAI API Key**: Set `OPENAI_API_KEY` environment variable
2. **Build Failures**: Run `go mod tidy` to sync dependencies
3. **Permission Issues**: Ensure `$GOPATH/bin` is in your `$PATH`

### Debugging Tools

- **Verbose Mode**: Use `-v` flag with commands
- **Go Build Tags**: Use for conditional compilation
- **Logging**: Built-in logging for debugging generation process

## Resources

### Documentation

- **Main README**: [README.md](README.md)
- **GitHub Repository**: https://github.com/RepoRadio/reporadio-cli
- **Go Documentation**: `go doc`
- **Conventional Commits**: https://www.conventionalcommits.org/
### Examples

- **Welcome Podcast**: `.reporadio/welcome/episodes/`
- **Contributor Podcast**: `.reporadio/contributor/episodes/`

### Contact

- **Issues**: https://github.com/RepoRadio/reporadio-cli/issues
- **Email**: hello@reporad.io
