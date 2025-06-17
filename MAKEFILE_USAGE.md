# OnionCLI Makefile Usage Guide

This document provides a comprehensive guide to using the Makefile for OnionCLI development, building, and deployment.

## Quick Start

```bash
# Show all available commands
make help

# Quick development cycle (format + build)
make quick

# Full development setup
make dev

# Run the application
make run
```

## Available Commands

### 🚀 Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Full development setup (clean, deps, build) |
| `make run` | Run the application in development mode |
| `make quick` | Quick development cycle (format + build) |
| `make full` | Full development cycle (clean, deps, check, build) |

### 🎬 Demo Commands

| Command | Description |
|---------|-------------|
| `make run-demo` | Run the collections & environments demo |
| `make run-config-demo` | Run the configuration management demo |
| `make run-performance-demo` | Run the performance enhancements demo |

### 🔨 Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary for current platform |
| `make build-all` | Build binaries for all platforms (Linux, macOS, Windows) |
| `make build-linux` | Build binary for Linux |
| `make build-darwin` | Build binary for macOS |
| `make build-windows` | Build binary for Windows |

### 📦 Dependency Management

| Command | Description |
|---------|-------------|
| `make deps` | Download and verify dependencies |
| `make deps-update` | Update all dependencies to latest versions |
| `make deps-clean` | Clean Go module cache |

### 🧪 Testing Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report (generates HTML) |
| `make test-race` | Run tests with race condition detection |
| `make benchmark` | Run performance benchmarks |

### 🔍 Code Quality Commands

| Command | Description |
|---------|-------------|
| `make fmt` | Format code using `go fmt` |
| `make vet` | Run `go vet` static analysis |
| `make lint` | Run golangci-lint (requires installation) |
| `make check` | Run all code quality checks (fmt, vet, lint, test) |

### 📦 Installation Commands

| Command | Description |
|---------|-------------|
| `make install` | Install binary to `$GOPATH/bin` |
| `make uninstall` | Remove installed binary |

### 🚀 Release Commands

| Command | Description |
|---------|-------------|
| `make release` | Create a release (clean, check, build-all) |
| `make package` | Create release packages (tar.gz, zip) |

### 🧹 Cleanup Commands

| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts |
| `make clean-all` | Clean everything including module cache |

### 🐳 Docker Commands (Optional)

| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

### 📊 Utility Commands

| Command | Description |
|---------|-------------|
| `make version` | Show version information |
| `make size` | Show binary size |
| `make deps-graph` | Generate dependency graph (requires graphviz) |
| `make help` | Show help message with all commands |

## Common Workflows

### 🔄 Daily Development

```bash
# Start development session
make dev

# Make changes to code...

# Quick build and test
make quick

# Run the application
make run

# Test your changes
make test
```

### 🚀 Preparing a Release

```bash
# Full quality check and build
make check

# Create release builds for all platforms
make release

# Create distribution packages
make package
```

### 🧪 Testing Workflow

```bash
# Run basic tests
make test

# Run tests with coverage
make test-coverage

# Run race condition tests
make test-race

# Run benchmarks
make benchmark
```

### 🔍 Code Quality Workflow

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run linter (install golangci-lint first)
make lint

# Run all quality checks
make check
```

## Build Artifacts

### Directory Structure

```
OnionCLI/
├── build/              # Single platform builds
│   ├── onioncli        # Current platform binary
│   └── coverage.html   # Test coverage report
├── dist/               # Multi-platform builds
│   ├── onioncli_linux_amd64
│   ├── onioncli_darwin_amd64
│   ├── onioncli_windows_amd64.exe
│   └── *.tar.gz, *.zip # Release packages
└── Makefile
```

### Binary Information

The Makefile embeds version information into binaries:

- **Version**: Git tag or commit hash
- **Build Time**: UTC timestamp
- **Git Commit**: Short commit hash

View this information with:
```bash
make version
```

## Prerequisites

### Required Tools

- **Go 1.19+**: For building and running
- **Git**: For version information
- **Make**: For running Makefile commands

### Optional Tools

- **golangci-lint**: For advanced linting
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- **graphviz**: For dependency graphs
  ```bash
  # Ubuntu/Debian
  sudo apt-get install graphviz
  
  # macOS
  brew install graphviz
  ```

- **Docker**: For containerized builds
  ```bash
  # Install Docker from https://docker.com
  ```

## Configuration

### Environment Variables

The Makefile respects these environment variables:

- `GOOS`: Target operating system
- `GOARCH`: Target architecture
- `GOPATH`: Go workspace path

### Customization

You can override default values:

```bash
# Build for specific platform
GOOS=windows GOARCH=amd64 make build

# Use different binary name
BINARY_NAME=my-onioncli make build
```

## Troubleshooting

### Common Issues

1. **"make: command not found"**
   - Install make: `sudo apt-get install make` (Ubuntu) or `brew install make` (macOS)

2. **"golangci-lint: command not found"**
   - Install golangci-lint or skip with `make fmt vet test`

3. **Permission denied on binary**
   - Make executable: `chmod +x build/onioncli`

4. **Module download issues**
   - Clean and retry: `make clean-all && make deps`

### Getting Help

- Run `make help` for command overview
- Check individual command output for specific errors
- Ensure Go 1.19+ is installed: `go version`

## Examples

### Cross-Platform Build

```bash
# Build for all platforms
make build-all

# Check what was built
ls -la dist/
```

### Development with Testing

```bash
# Full development cycle with testing
make clean
make deps
make check
make build
make test-coverage
```

### Release Preparation

```bash
# Prepare a complete release
make clean-all
make deps
make check
make release
make package

# Verify release
ls -la dist/
make version
```

This Makefile provides a complete build system for efficient OnionCLI development and deployment! 🚀
