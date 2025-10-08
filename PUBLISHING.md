# Quick Publishing Guide

This guide provides the essential steps to publish mongo-essential to go.dev and Homebrew.

## 🚀 Quick Release Process

### 1. Prerequisites (One-time setup)

1. **Create Homebrew Tap Repository**:
   ```bash
   # Create a new repository on GitHub: homebrew-mongo-essential
   # Initialize it with:
   mkdir homebrew-mongo-essential
   cd homebrew-mongo-essential
   mkdir Formula
   cp ../mongo-migration-tool/Formula/mongo-essential.rb Formula/
   git init
   git add .
   git commit -m "Initial Homebrew tap"
   git remote add origin git@github.com:jocham/homebrew-mongo-essential.git
   git push -u origin main
   ```

2. **Set GitHub Secrets**:
   - Go to your repository Settings → Secrets and Variables → Actions
   - Add `HOMEBREW_TAP_TOKEN`: Create a Personal Access Token with `repo` permissions

### 2. Release (Every time)

1. **Prepare and tag release**:
   ```bash
   # Make sure everything is committed
   git add .
   git commit -m "Prepare v1.0.0 release"
   git push

   # Create and push tag
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Automated publishing happens**:
   - ✅ GitHub Actions builds binaries
   - ✅ Updates Homebrew formula automatically
   - ✅ Pushes Docker images to GitHub Container Registry
   - ✅ Triggers pkg.go.dev update

### 3. Verify (5 minutes after release)

1. **Check go.dev**: https://pkg.go.dev/github.com/jocham/mongo-essential
2. **Test Homebrew**:
   ```bash
   brew tap jocham/mongo-essential
   brew install mongo-essential
   mongo-essential version
   ```
3. **Test Go library**:
   ```bash
   go get github.com/jocham/mongo-essential@latest
   ```

## 📦 Distribution Summary

| What | Where | How Users Install |
|------|-------|------------------|
| **Go Library** | pkg.go.dev | `go get github.com/jocham/mongo-essential` |
| **CLI via Homebrew** | Homebrew tap | `brew install jocham/mongo-essential/mongo-essential` |
| **CLI binary** | GitHub Releases | Download from releases page |
| **Docker** | GitHub Container Registry | `docker pull ghcr.io/jocham/mongo-essential` |

## 🎯 That's It!

The entire process is automated. Just push a tag and everything else happens automatically!

For detailed information, troubleshooting, and manual steps, see [RELEASE.md](./RELEASE.md).