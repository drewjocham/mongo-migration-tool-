# Release Process for mongo-essential

This document outlines the process for publishing mongo-essential to both go.dev (pkg.go.dev) and Homebrew.

## Overview

mongo-essential is published to two main distribution channels:

1. **go.dev (pkg.go.dev)**: Automatic publishing for Go library usage
2. **Homebrew**: Package manager for CLI tool installation
3. **GitHub Releases**: Binary downloads and Docker images

## Prerequisites

### GitHub Repository Setup

1. **Repository**: Must be public at `github.com/jocham/mongo-essential`
2. **GitHub Tokens**: Required secrets in repository settings:
   - `HOMEBREW_TAP_TOKEN`: Personal access token for updating Homebrew tap

### Homebrew Tap Repository

1. **Create tap repository**: `github.com/jocham/homebrew-mongo-essential`
2. **Repository structure**:
   ```
   homebrew-mongo-essential/
   ├── Formula/
   │   └── mongo-essential.rb
   └── README.md
   ```

## Release Process

### 1. Prepare Release

1. **Update version information** (if needed):
   - Update `CHANGELOG.md` with new features, fixes, and changes
   - Ensure `README.md` is up to date

2. **Ensure tests pass**:
   ```bash
   make ci-test
   make test-examples
   make docker-build
   ```

3. **Test local build**:
   ```bash
   go build -o mongo-essential .
   ./mongo-essential version
   ```

### 2. Create GitHub Release

1. **Create and push tag**:
   ```bash
   # Create a new tag (use semantic versioning)
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Automated release process**:
   - GitHub Actions will trigger the release workflow
   - GoReleaser builds binaries for all platforms
   - Docker images are built and pushed to GitHub Container Registry
   - Homebrew formula is automatically updated
   - pkg.go.dev is notified of the new release

### 3. Manual Steps (if needed)

If the automated process fails, you can manually:

#### Update Homebrew Formula

1. **Calculate SHA256 checksums**:
   ```bash
   # Download release artifacts and calculate checksums
   shasum -a 256 mongo-essential-darwin-amd64.tar.gz
   shasum -a 256 mongo-essential-darwin-arm64.tar.gz
   shasum -a 256 mongo-essential-linux-amd64.tar.gz
   shasum -a 256 mongo-essential-linux-arm64.tar.gz
   ```

2. **Update formula** in `homebrew-mongo-essential` repository:
   ```ruby
   class MongoEssential < Formula
     desc "MongoDB migration and AI-powered database analysis tool"
     homepage "https://github.com/jocham/mongo-essential"
     version "1.0.0"  # Update version
     license "MIT"
     
     # Update URLs and SHA256 checksums
     # ...
   end
   ```

#### Trigger pkg.go.dev Update

1. **Manual trigger**:
   ```bash
   # Trigger pkg.go.dev to fetch the new version
   curl -X POST "https://proxy.golang.org/github.com/jocham/mongo-essential/@v1.0.0"
   curl -X POST "https://proxy.golang.org/github.com/jocham/mongo-essential/@latest"
   ```

## Verification

### 1. Verify go.dev Publishing

1. **Check pkg.go.dev**: Visit https://pkg.go.dev/github.com/jocham/mongo-essential
2. **Test library installation**:
   ```bash
   # In a test directory
   go mod init test-mongo-essential
   go get github.com/jocham/mongo-essential@latest
   ```

### 2. Verify Homebrew Publishing

1. **Test Homebrew installation**:
   ```bash
   brew tap jocham/mongo-essential
   brew install mongo-essential
   mongo-essential version
   ```

2. **Verify formula**:
   ```bash
   brew info mongo-essential
   brew test mongo-essential
   ```

### 3. Verify GitHub Release

1. **Check release page**: https://github.com/jocham/mongo-essential/releases
2. **Test binary downloads**:
   - Download appropriate binary for your platform
   - Verify it runs: `./mongo-essential version`

### 4. Verify Docker Images

1. **Test Docker images**:
   ```bash
   docker run --rm ghcr.io/jocham/mongo-essential:latest version
   docker run --rm ghcr.io/jocham/mongo-essential:v1.0.0 --help
   ```

## Distribution Channels Summary

| Channel | URL | Installation Method |
|---------|-----|-------------------|
| **Go Library** | https://pkg.go.dev/github.com/jocham/mongo-essential | `go get github.com/jocham/mongo-essential` |
| **Homebrew** | https://github.com/jocham/homebrew-mongo-essential | `brew install jocham/mongo-essential/mongo-essential` |
| **GitHub Releases** | https://github.com/jocham/mongo-essential/releases | Download binary directly |
| **Docker** | https://ghcr.io/jocham/mongo-essential | `docker pull ghcr.io/jocham/mongo-essential` |

## Troubleshooting

### Common Issues

1. **Homebrew formula fails**:
   - Check SHA256 checksums match release artifacts
   - Ensure tap repository has correct permissions
   - Verify `HOMEBREW_TAP_TOKEN` secret is valid

2. **pkg.go.dev not updating**:
   - Wait up to 30 minutes for automatic indexing
   - Manually trigger with curl commands above
   - Ensure repository is public and has proper Go module structure

3. **GoReleaser fails**:
   - Check Go version compatibility
   - Verify all required secrets are set
   - Ensure `.goreleaser.yml` syntax is valid

### Debug Commands

```bash
# Test GoReleaser locally
goreleaser check
goreleaser build --snapshot --rm-dist

# Check module proxy
curl "https://proxy.golang.org/github.com/jocham/mongo-essential/@latest"

# Test Homebrew formula locally
brew install --build-from-source ./Formula/mongo-essential.rb
```

## Post-Release Tasks

1. **Update documentation**: Ensure all documentation reflects the new version
2. **Social media**: Announce the release (Twitter, LinkedIn, etc.)
3. **Community**: Update any community resources or examples
4. **Monitoring**: Monitor for any issues or bug reports

## Version Strategy

Use [Semantic Versioning (SemVer)](https://semver.org/):

- **MAJOR** (v2.0.0): Incompatible API changes
- **MINOR** (v1.1.0): New functionality, backwards compatible
- **PATCH** (v1.0.1): Bug fixes, backwards compatible

Examples:
- `v1.0.0`: Initial stable release
- `v1.1.0`: Added AI analysis features
- `v1.0.1`: Fixed migration rollback bug
- `v2.0.0`: Breaking changes to Migration interface