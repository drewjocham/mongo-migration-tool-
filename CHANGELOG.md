# Changelog

All notable changes to mongo-migrate will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-08

### Added
- **AI-Powered Database Analysis**: Multi-provider AI support (OpenAI, Google Gemini, Anthropic Claude)
- **Comprehensive Migration System**: Version-controlled database migrations with up/down support
- **Oplog & Replication Analysis**: Deep analysis of MongoDB replication health and oplog patterns
- **Change Stream Optimization**: Real-time data processing pattern analysis
- **Google Docs Integration**: Export professional reports directly to Google Docs
- **Certificate Management**: Debug and fix SSL/TLS certificate issues
- **Cloud Provider Support**: Optimized for STACKIT, AWS, Azure, GCP
- **CLI Interface**: Intuitive command-line interface built with Cobra

### Features

#### Database Migrations
- Version-controlled schema changes
- Rollback capabilities
- Migration status tracking
- Force migration marking
- CI/CD integration support

#### AI Analysis Commands
- `mongo-migrate ai analyze` - Comprehensive database analysis
- `mongo-migrate ai schema` - Schema-focused analysis
- `mongo-migrate ai performance` - Performance optimization recommendations
- `mongo-migrate ai oplog` - Oplog and replication health analysis
- `mongo-migrate ai changestream` - Change stream configuration optimization

#### Certificate Tools
- `mongo-migrate cert diagnose` - Certificate verification diagnostics
- `mongo-migrate cert check` - Host certificate validation
- `mongo-migrate cert fix` - Automatic certificate issue resolution

#### Google Docs Export
- Professional report formatting
- Automatic sharing with team members
- Folder organization in Google Drive
- Custom document titles and metadata

### Configuration
- Environment variable support
- Multiple cloud provider configurations
- AI provider settings (OpenAI, Gemini, Claude)
- Google Docs integration settings
- SSL/TLS configuration options

### Documentation
- Comprehensive README with installation instructions
- AI analysis guide with examples
- Configuration reference
- API documentation with Go doc comments

### Technical Details
- Built with Go 1.24
- Uses Cobra CLI framework
- MongoDB Go Driver integration
- Multi-provider AI client support
- Google APIs integration
- Docker containerization support

## [Unreleased]

### Planned
- Anthropic Claude API integration completion
- Historical analysis comparison
- Integration with monitoring tools
- Custom analysis prompts
- Automated recommendations tracking

---

## Release Notes

### Installation

#### Homebrew
```bash
brew tap jocham/mongo-migrate
brew install mongo-migrate
```

#### Go Install
```bash
go install github.com/jocham/mongo-migrate@latest
```

#### Docker
```bash
docker run --rm ghcr.io/jocham/mongo-migrate:latest --help
```

### Breaking Changes
- None (initial release)

### Dependencies
- Go 1.24 or later
- MongoDB 4.4 or later
- Valid AI provider API key (for AI analysis features)
- Google Cloud credentials (for Google Docs integration)

### Compatibility
- macOS (ARM64, AMD64)
- Linux (ARM64, AMD64) 
- Windows (AMD64)
- Docker containers
- Kubernetes deployments