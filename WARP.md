# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Repository Overview

This is **mongo-essential** (formerly mongo-migration-tool), a MongoDB migration and AI-powered database analysis tool built in Go. It provides Liquibase/Flyway-style migrations for MongoDB with additional AI analysis capabilities using OpenAI GPT, Google Gemini, and Anthropic Claude.

**Available as both a Go library and CLI tool:**
- **Library**: Import `github.com/jocham/mongo-essential/migration` and `github.com/jocham/mongo-essential/config` in your Go projects
- **CLI**: Install via Homebrew, go install, or download binary for standalone use

## Build and Development Commands

### Essential Development Commands

```bash
# Install dependencies and set up development environment
make dev-setup

# Build the binary
make build

# Build for all platforms (CI/Release)
make build-all

# Run tests
make test
make test-library
make test-coverage
make test-examples

# Code quality
make lint
make format
make vet

# Run all CI checks locally
make ci-test
```

### Migration Commands

```bash
# Build and check migration status
make migration-status

# Create new migration (requires DESC parameter)
make create-migration DESC="add user indexes"

# Run pending migrations
make migration-up

# Rollback to specific version (requires VERSION parameter)
make migration-down VERSION="20231201_001"
```

### Docker Development

```bash
# Build and run with Docker
make docker-build
make docker-run

# Use docker-compose for full stack testing
make docker-compose-up
make docker-compose-down

# Start local MongoDB for testing
make db-up
make db-down
```

## Architecture Overview

### Core Structure

**CLI Architecture**: Built with Cobra framework, commands are in `cmd/` directory
- `cmd/root.go`: Main command setup, MongoDB connection, configuration loading
- Individual command files: `up.go`, `down.go`, `status.go`, `create.go`, etc.

**Migration Engine**: Core logic in `internal/migration/`
- `engine.go`: Migration execution engine with up/down/status/force operations
- `types.go`: Migration interfaces and data structures
- Migrations stored in `migrations/` directory as Go files implementing the `Migration` interface

**Configuration**: Environment-based configuration in `config/` (public package)
- Supports MongoDB connections, SSL/TLS, AI providers, Google Docs integration
- Uses godotenv for `.env` files and caarlos0/env for parsing

**Public Library API**: Exposed packages for external use
- `migration/`: Public migration engine and interfaces  
- `config/`: Configuration management and loading
- Comprehensive documentation and examples in `examples/` directory

### Key Components

1. **Migration Interface**: Each migration implements `Version()`, `Description()`, `Up()`, and `Down()` methods
2. **Migration Records**: Tracked in MongoDB collection (`schema_migrations` by default)
3. **AI Analysis**: Multi-provider support (OpenAI, Gemini, Claude) for database insights
4. **Certificate Management**: SSL/TLS troubleshooting tools for cloud providers like STACKIT

### Binary Names and Module Path

- **Module Path**: `github.com/jocham/mongo-essential`
- **Binary Name**: `mongo-migrate` (in Makefile) or `mongo-essential` (in CLI commands)
- **Docker Image**: `mongo-migration-tool`

## Configuration

### Environment Configuration

Configuration is primarily environment-based using `.env` files:

```bash
# Copy example configuration
cp .env.example .env
```

Key configuration categories:
- **MongoDB Connection**: `MONGO_URL`, `MONGO_DATABASE`, `MONGO_USERNAME`, `MONGO_PASSWORD`
- **SSL/TLS**: `MONGO_SSL_ENABLED`, `MONGO_SSL_INSECURE` (important for cloud providers)
- **AI Analysis**: `AI_ENABLED`, `AI_PROVIDER`, `OPENAI_API_KEY`, etc.
- **Google Docs**: `GOOGLE_DOCS_ENABLED`, `GOOGLE_CREDENTIALS_PATH`

### Migration Configuration

- **Migrations Path**: `MIGRATIONS_PATH` (default: `./migrations`)
- **Collection Name**: `MIGRATIONS_COLLECTION` (default: `schema_migrations`)

## Testing and Quality

### Current Test Setup

- **No test files currently exist** in the codebase
- CI pipeline includes test, lint, and security scanning jobs
- Uses `golangci-lint` for linting
- Coverage reporting with Codecov

### Test Development Patterns

When adding tests, follow these patterns:
- Use standard Go testing package
- Test migration engine functionality with test MongoDB instances
- Mock AI providers for analysis testing
- Integration tests with docker-compose setup

## Development Workflow

### Adding New Migrations

1. **Create Migration**: Use `make create-migration DESC="description"`
2. **Implement Interface**: Follow pattern in `migrations/20231201_001_create_users_collection.go`
3. **Register Migration**: Currently manual registration in `cmd/root.go:loadMigrations()`

### AI Feature Development

- AI providers configured through interfaces
- Multi-provider support: OpenAI GPT-4o, Google Gemini 1.5, Claude 3.5 Sonnet
- Google Docs integration for report generation
- Environment-based API key management

### Cloud Provider Support

- **STACKIT Cloud**: Primary focus with SSL/TLS certificate management
- **Multi-cloud**: AWS, Azure, GCP support mentioned
- **Certificate Troubleshooting**: Dedicated `cert` command for debugging SSL issues

## CI/CD Pipeline

### GitHub Actions Workflow

- **Test Job**: Go tests with coverage
- **Lint Job**: golangci-lint validation
- **Security Job**: Gosec security scanning
- **Build Job**: Multi-platform binary builds (Linux, macOS, Windows × AMD64/ARM64)
- **Docker Job**: Multi-architecture container builds
- **Release Job**: GitHub releases with checksums
- **Homebrew Job**: Automatic Homebrew formula updates

### Release Process

1. **Create Release**: GitHub release triggers full pipeline
2. **Multi-platform Builds**: Binaries for all supported platforms
3. **Container Images**: Published to GitHub Container Registry
4. **Homebrew**: Automatic formula update in `jocham/homebrew-mongo-essential`
5. **Documentation**: pkg.go.dev updates

## Development Notes

### Migration System Design

- **Version-based**: Migrations use timestamp-based versions (`YYYYMMDD_NNN`)
- **Bidirectional**: Both up and down migrations required
- **State Tracking**: Migration status stored in MongoDB collection
- **Force Option**: Can mark migrations as applied without execution
- **Checksum Validation**: MD5 checksums for migration integrity

### AI Integration Architecture

- **Provider Abstraction**: Support for multiple AI providers
- **Analysis Types**: Schema, performance, oplog, change stream analysis
- **Report Generation**: Professional reports via Google Docs API
- **Flexible Configuration**: Environment-based provider selection

### Go Module Considerations

- **Go Version**: Currently 1.25.0 (cutting edge)
- **Key Dependencies**: Cobra CLI, MongoDB driver, AI SDK packages, Google APIs
- **Build Constraints**: CGO disabled for static binaries
