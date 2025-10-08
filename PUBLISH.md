# Publication Checklist

This guide walks you through publishing mongo-migrate to Homebrew and Go package documentation.

## 🚀 Pre-Publication Checklist

### ✅ Code Quality
- [ ] All tests pass: `go test ./...`
- [ ] Linting passes: `golangci-lint run`
- [ ] Code is formatted: `gofmt -s -w .`
- [ ] Go mod is clean: `go mod tidy && go mod verify`
- [ ] Build succeeds: `go build -o mongo-migrate .`

### ✅ Documentation
- [ ] README.md is comprehensive and up-to-date
- [ ] AI_ANALYSIS.md has examples for all features
- [ ] CHANGELOG.md documents all changes
- [ ] Go doc comments are present on all public APIs
- [ ] CONTRIBUTING.md guides are clear

### ✅ Configuration
- [ ] go.mod has correct module path: `github.com/jocham/mongo-migrate`
- [ ] VERSION file contains correct version (e.g., `v1.0.0`)
- [ ] LICENSE file is present (MIT)
- [ ] .env.example has all configuration options

### ✅ Testing
- [ ] Manual testing with real MongoDB instance
- [ ] AI providers tested (OpenAI, Gemini)
- [ ] Certificate utilities tested
- [ ] Google Docs integration tested (optional)
- [ ] Cross-platform compatibility verified

## 📦 GitHub Repository Setup

### 1. Create GitHub Repository

```bash
# If not already done, create repository on GitHub:
# https://github.com/new
# Repository name: mongo-migrate
# Description: MongoDB migration tool with AI-powered database analysis
# Public repository
# Initialize with README: No (we have our own)
```

### 2. Push Code to GitHub

```bash
# Initialize git if needed
git init
git branch -M main

# Add remote
git remote add origin https://github.com/jocham/mongo-migrate.git

# Add all files
git add .

# Commit
git commit -m "Initial release v1.0.0

- AI-powered database analysis with OpenAI, Gemini, Claude
- Comprehensive migration system with up/down support  
- Oplog and replication health analysis
- Change stream optimization recommendations
- Google Docs integration for professional reports
- Certificate management and troubleshooting tools
- Cross-platform support (macOS, Linux, Windows)"

# Push to GitHub
git push -u origin main
```

### 3. Configure Repository Settings

On GitHub, go to your repository settings:

- [ ] **General**: Add description and topics
  - Description: "MongoDB migration tool with AI-powered database analysis"
  - Topics: `mongodb`, `database`, `migration`, `ai`, `golang`, `cli`, `devops`

- [ ] **Security**: Enable security features
  - [ ] Dependency graph
  - [ ] Dependabot alerts
  - [ ] Code scanning (GitHub CodeQL)

- [ ] **Pages**: Enable GitHub Pages (optional)
  - Source: Deploy from a branch
  - Branch: main, /docs folder

## 🏠 Homebrew Publication

### 1. Create Homebrew Tap Repository

```bash
# Create new repository on GitHub:
# Repository name: homebrew-mongo-migrate
# Description: Homebrew tap for mongo-migrate
# Public repository

# Clone the tap repository
git clone https://github.com/jocham/homebrew-mongo-migrate.git
cd homebrew-mongo-migrate

# Create Formula directory
mkdir -p Formula

# Copy the formula
cp ../mongo-migrate/homebrew/mongo-migrate.rb Formula/
```

### 2. Update Formula with Release Info

After creating the first GitHub release:

```bash
# Get the release tarball SHA256
TARBALL_URL="https://github.com/jocham/mongo-migrate/archive/v1.0.0.tar.gz"
SHA256=$(curl -sL "$TARBALL_URL" | sha256sum | cut -d' ' -f1)

# Update the formula
sed -i "s/SHA256_PLACEHOLDER/$SHA256/" Formula/mongo-migrate.rb
```

### 3. Test Formula Locally

```bash
# Install from local tap
brew install --build-from-source ./Formula/mongo-migrate.rb

# Test installation
mongo-migrate --version
mongo-migrate --help

# Test certificate command (doesn't need MongoDB)
mongo-migrate cert diagnose

# Uninstall for clean testing
brew uninstall mongo-migrate
```

### 4. Publish Homebrew Tap

```bash
# Commit and push the formula
git add Formula/mongo-migrate.rb
git commit -m "Add mongo-migrate formula v1.0.0"
git push origin main
```

### 5. Make Formula Available

Users can now install with:
```bash
brew tap jocham/mongo-migrate
brew install mongo-migrate
```

## 📚 Go Package Documentation

### 1. Ensure Module is Public

```bash
# Verify module path in go.mod
grep "module" go.mod
# Should output: module github.com/jocham/mongo-migrate

# Tag and push the release
git tag v1.0.0
git push origin v1.0.0
```

### 2. Trigger Go Module Indexing

```bash
# Request module indexing (run after GitHub release)
curl -X POST "https://proxy.golang.org/github.com/jocham/mongo-migrate/@v/v1.0.0.info"

# Check if module is available
go list -m github.com/jocham/mongo-migrate@v1.0.0
```

### 3. Verify pkg.go.dev

- Wait 5-10 minutes after tagging
- Visit: https://pkg.go.dev/github.com/jocham/mongo-migrate
- Documentation should automatically appear

### 4. Improve pkg.go.dev Presentation

Add these to improve documentation:

```go
// Package main provides the mongo-migrate CLI tool.
//
// mongo-migrate is a comprehensive MongoDB migration and database analysis tool
// with AI-powered insights, similar to Liquibase/Flyway for MongoDB.
//
// Key features:
//   - Version-controlled database migrations
//   - AI-powered database analysis (OpenAI, Gemini, Claude)
//   - Oplog and replication health monitoring
//   - Change stream optimization
//   - Google Docs integration
//   - Certificate troubleshooting
//
// Example usage:
//   go install github.com/jocham/mongo-migrate@latest
//   mongo-migrate ai analyze --provider openai
package main
```

## 🎯 Release Process

### 1. Create GitHub Release

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0

Initial release with AI-powered database analysis"
git push origin v1.0.0
```

On GitHub:
- [ ] Go to Releases → Create a new release
- [ ] Choose tag: v1.0.0
- [ ] Release title: "mongo-migrate v1.0.0"
- [ ] Description: Copy from CHANGELOG.md
- [ ] Upload binary assets (optional, CI/CD will do this)
- [ ] Publish release

### 2. Automated Release Process

The GitHub Actions workflow will automatically:
- [ ] Build binaries for multiple platforms
- [ ] Create Docker images
- [ ] Update Homebrew formula
- [ ] Generate release notes
- [ ] Upload release artifacts

## 🔍 Post-Publication Verification

### Test Homebrew Installation
```bash
# Test fresh installation
brew tap jocham/mongo-migrate
brew install mongo-migrate

# Verify installation
mongo-migrate --version
mongo-migrate --help
mongo-migrate cert diagnose
```

### Test Go Installation
```bash
# Test from different directory
cd /tmp
go install github.com/jocham/mongo-migrate@latest
mongo-migrate --version
```

### Test Docker Image
```bash
# Test Docker image
docker run --rm ghcr.io/jocham/mongo-migrate:latest --version
docker run --rm ghcr.io/jocham/mongo-migrate:latest cert diagnose
```

### Verify Documentation
- [ ] https://pkg.go.dev/github.com/jocham/mongo-migrate shows documentation
- [ ] README.md renders correctly on GitHub
- [ ] All badges in README are working
- [ ] Examples in documentation work

## 📈 Post-Launch Tasks

### Community Building
- [ ] Submit to awesome-go: https://github.com/avelino/awesome-go
- [ ] Post on relevant subreddits (r/golang, r/MongoDB)
- [ ] Share on Twitter/LinkedIn with relevant hashtags
- [ ] Write blog post about the tool

### Monitoring and Maintenance
- [ ] Set up GitHub issue templates
- [ ] Monitor package download stats
- [ ] Respond to community feedback
- [ ] Plan next version features

### SEO and Discovery
- [ ] Add to MongoDB community tools list
- [ ] Submit to tool directories
- [ ] Update personal/company portfolio
- [ ] Create demo videos

## 🛠️ Development Commands Reference

### Local Development
```bash
# Build and test locally
make build
make test
make lint

# Or manually:
go build -o mongo-migrate .
go test ./...
golangci-lint run
```

### Release Preparation
```bash
# Update version
echo "v1.1.0" > VERSION

# Update changelog
vim CHANGELOG.md

# Tag and release
git add VERSION CHANGELOG.md
git commit -m "Bump version to v1.1.0"
git tag v1.1.0
git push origin main v1.1.0
```

### Formula Updates
```bash
# Update Homebrew formula after release
cd ../homebrew-mongo-migrate
./scripts/update-formula.sh v1.1.0
git add Formula/mongo-migrate.rb
git commit -m "Update mongo-migrate to v1.1.0"
git push origin main
```

## 🎉 Success Criteria

Publication is successful when:

- [ ] ✅ **Homebrew**: `brew install jocham/mongo-migrate/mongo-migrate` works
- [ ] ✅ **Go Install**: `go install github.com/jocham/mongo-migrate@latest` works  
- [ ] ✅ **Documentation**: https://pkg.go.dev/github.com/jocham/mongo-migrate is live
- [ ] ✅ **GitHub**: Repository has proper README, issues, and releases
- [ ] ✅ **Docker**: `docker run ghcr.io/jocham/mongo-migrate:latest` works
- [ ] ✅ **Cross-platform**: Binaries work on macOS, Linux, Windows

---

**Ready to share mongo-migrate with the world! 🚀**