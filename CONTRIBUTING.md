# Contributing to mongo-migrate

Thank you for considering contributing to mongo-migrate! This document provides guidelines and information for contributors.

## 🎯 How to Contribute

### Reporting Bugs

1. **Search existing issues** to avoid duplicates
2. **Use the bug report template** when creating new issues
3. **Include reproduction steps** and system information
4. **Add relevant labels** (bug, enhancement, etc.)

### Suggesting Enhancements

1. **Check existing feature requests** first
2. **Use the feature request template** 
3. **Describe the use case** and expected behavior
4. **Consider backward compatibility** implications

### Contributing Code

1. **Fork the repository** and create a feature branch
2. **Follow coding standards** (see below)
3. **Add tests** for new functionality
4. **Update documentation** as needed
5. **Submit a pull request** with a clear description

## 🛠️ Development Setup

### Prerequisites

- Go 1.24 or later
- MongoDB 4.4+ for testing
- Git
- Make (optional)

### Local Development

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/mongo-migrate.git
cd mongo-migrate

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build the binary
go build -o mongo-migrate .

# Run with development MongoDB
docker run -d -p 27017:27017 --name mongo-dev mongo:latest

# Test the CLI
export MONGO_URL=mongodb://localhost:27017
export MONGO_DATABASE=test_db
./mongo-migrate status
```

### Project Structure

```
mongo-migrate/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command and global configuration
│   ├── ai.go              # AI analysis commands
│   ├── cert.go            # Certificate management
│   └── *.go               # Other command implementations
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go      # Environment variable handling
│   └── migration/         # Migration engine core
│       ├── engine.go      # Migration execution logic
│       └── types.go       # Data structures
├── migrations/            # Example migrations
├── .github/               # GitHub Actions workflows
├── docs/                  # Additional documentation
└── homebrew/              # Homebrew formula
```

## 📋 Coding Standards

### Go Style Guide

Follow the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) and these additional guidelines:

#### Code Organization
- Use meaningful package names
- Keep functions focused and small
- Prefer composition over inheritance
- Use interfaces for testability

#### Naming Conventions
```go
// Good: Descriptive names
type DatabaseAnalyzer struct {
    provider AIProvider
    config   *Config
}

func (da *DatabaseAnalyzer) AnalyzeSchema(ctx context.Context) error {
    // Implementation
}

// Avoid: Abbreviated or unclear names
type DA struct {
    p AI
    c *Cfg
}
```

#### Error Handling
```go
// Good: Wrap errors with context
if err := db.Connect(); err != nil {
    return fmt.Errorf("failed to connect to database: %w", err)
}

// Good: Check errors immediately
result, err := someOperation()
if err != nil {
    return err
}
processResult(result)
```

#### Documentation
```go
// Package config provides configuration management for mongo-migrate.
// It handles loading from environment variables and supports multiple
// cloud providers and AI service integrations.
package config

// Config holds all configuration options for mongo-migrate.
// It supports MongoDB connection settings, AI provider configurations,
// Google Docs integration, and SSL/TLS settings.
type Config struct {
    // MongoURL is the MongoDB connection string
    MongoURL string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
    // Database name to connect to (required)
    Database string `env:"MONGO_DATABASE,required"`
}
```

### CLI Commands

Follow Cobra best practices:
- Use consistent command naming
- Provide helpful descriptions and examples
- Validate input parameters
- Handle errors gracefully

```go
var exampleCmd = &cobra.Command{
    Use:   "example [args]",
    Short: "Brief description of the command",
    Long: `Detailed description of the command with examples.
    
This command does X, Y, and Z. It's useful for:
- Use case 1
- Use case 2

Examples:
  mongo-migrate example --flag value
  mongo-migrate example --detailed-analysis`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation with proper error handling
        return nil
    },
}
```

## 🧪 Testing Guidelines

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:  "valid input",
            input: validInput,
            want:  expectedOutput,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test component interactions
- **CLI tests**: Test command-line interface behavior
- **End-to-end tests**: Test complete workflows

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/config/...

# Run with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

## 🔄 Pull Request Process

### Before Submitting

1. **Run tests locally**: `go test ./...`
2. **Run linter**: `golangci-lint run`
3. **Check formatting**: `gofmt -s -w .`
4. **Update documentation** if needed
5. **Test manually** with different configurations

### Pull Request Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature causing existing functionality to break)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] All tests pass locally

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings introduced
```

### Review Process

1. **Automated checks** must pass (CI/CD pipeline)
2. **Code review** by maintainers
3. **Testing** in different environments
4. **Documentation** review if applicable
5. **Final approval** and merge

## 🏷️ Versioning and Releases

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Process

1. **Update VERSION file**: Bump version number
2. **Update CHANGELOG.md**: Document all changes
3. **Create GitHub release**: Tag and publish
4. **Homebrew formula**: Auto-updated via CI/CD
5. **Docker images**: Auto-built and published
6. **Go modules**: Auto-indexed by pkg.go.dev

## 🎨 Adding New Features

### AI Analysis Features

When adding new AI analysis capabilities:

1. **Add to `AIProvider` interface**: If needed
2. **Implement for all providers**: OpenAI, Gemini, Claude
3. **Add new command**: In `cmd/ai.go`
4. **Update documentation**: README and AI_ANALYSIS.md
5. **Add tests**: Unit and integration tests

### Migration Features

For new migration functionality:

1. **Extend `Migration` interface**: If needed
2. **Update `Engine` struct**: Add new methods
3. **Add CLI commands**: For user interaction
4. **Backward compatibility**: Ensure existing migrations work
5. **Database tests**: Test with real MongoDB

### Certificate Features

For certificate management:

1. **Add to `cmd/cert.go`**: New subcommands
2. **Platform support**: Ensure cross-platform compatibility
3. **Error handling**: Graceful failure modes
4. **Documentation**: Update troubleshooting guides

## 📚 Documentation

### Types of Documentation

1. **Code documentation**: Go doc comments
2. **API documentation**: Auto-generated from code
3. **User documentation**: README, guides, examples
4. **Developer documentation**: This file, architecture docs

### Writing Guidelines

- **Be clear and concise**: Avoid jargon
- **Include examples**: Show actual usage
- **Keep updated**: Documentation should match code
- **Test examples**: Ensure they actually work

### Documentation Structure

```markdown
# Title
Brief overview of the feature/component.

## Overview
What this does and why it's useful.

## Usage
### Basic Usage
Simple example with minimal configuration.

### Advanced Usage
Complex example with all options.

## Configuration
All available options with descriptions.

## Examples
Real-world scenarios and solutions.

## Troubleshooting
Common issues and solutions.
```

## 🐛 Debugging and Troubleshooting

### Debug Mode

Enable debug logging:
```bash
export DEBUG=true
export LOG_LEVEL=debug
mongo-migrate ai analyze --provider openai
```

### Common Issues

1. **MongoDB connection failures**: Check connection string and credentials
2. **AI API errors**: Verify API keys and rate limits  
3. **Certificate issues**: Check SSL/TLS configuration
4. **Build failures**: Ensure Go version and dependencies

### Development Tools

- **MongoDB Compass**: GUI for database inspection
- **Postman**: API testing for AI providers
- **Docker**: Isolated testing environments
- **golangci-lint**: Code quality checking

## 🤝 Community Guidelines

### Code of Conduct

- **Be respectful**: Treat everyone with kindness
- **Be inclusive**: Welcome all skill levels and backgrounds
- **Be constructive**: Provide helpful feedback
- **Be patient**: Remember we're all learning

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community chat
- **Pull Requests**: Code contributions and reviews

### Getting Help

1. **Check documentation**: README, guides, and code comments
2. **Search issues**: Someone might have asked already
3. **Ask questions**: Create a discussion or issue
4. **Join community**: Engage with other contributors

## 📄 License

By contributing to mongo-migrate, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to mongo-migrate! 🚀