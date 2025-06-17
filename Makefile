# OnionCLI Makefile
# A comprehensive build system for the OnionCLI project

# Project configuration
PROJECT_NAME := onioncli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go configuration
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOVERSION := $(shell go version | awk '{print $$3}')

# Build configuration
BUILD_DIR := build
DIST_DIR := dist
CMD_DIR := cmd/onioncli
MAIN_FILE := $(CMD_DIR)/main.go

# Binary names
BINARY_NAME := onioncli
BINARY_UNIX := $(BINARY_NAME)_unix_amd64
BINARY_LINUX := $(BINARY_NAME)_linux_amd64
BINARY_DARWIN := $(BINARY_NAME)_darwin_amd64
BINARY_WINDOWS := $(BINARY_NAME)_windows_amd64.exe

# Ldflags for version information
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)OnionCLI Build System$(NC)"
	@echo "$(CYAN)=====================$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Project Info:$(NC)"
	@echo "  Name:       $(PROJECT_NAME)"
	@echo "  Version:    $(VERSION)"
	@echo "  Go Version: $(GOVERSION)"
	@echo "  OS/Arch:    $(GOOS)/$(GOARCH)"

# Development targets
.PHONY: dev
dev: clean deps build ## Full development setup (clean, deps, build)
	@echo "$(GREEN)✅ Development environment ready!$(NC)"

.PHONY: run
run: ## Run the application in development mode
	@echo "$(BLUE)🚀 Running OnionCLI...$(NC)"
	$(GO) run $(MAIN_FILE)

.PHONY: run-demo
run-demo: ## Run the collections demo
	@echo "$(BLUE)🎬 Running collections demo...$(NC)"
	$(GO) run examples/collections_demo.go

.PHONY: run-config-demo
run-config-demo: ## Run the configuration demo
	@echo "$(BLUE)⚙️ Running configuration demo...$(NC)"
	$(GO) run examples/config_demo.go

.PHONY: run-performance-demo
run-performance-demo: ## Run the performance demo
	@echo "$(BLUE)⚡ Running performance demo...$(NC)"
	$(GO) run examples/performance_demo.go

# Build targets
.PHONY: build
build: ## Build the binary for current platform
	@echo "$(BLUE)🔨 Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)✅ Built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: build-all
build-all: build-linux build-darwin build-windows ## Build binaries for all platforms
	@echo "$(GREEN)✅ All platform builds completed!$(NC)"

.PHONY: build-linux
build-linux: ## Build binary for Linux
	@echo "$(BLUE)🐧 Building for Linux...$(NC)"
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_LINUX) $(MAIN_FILE)
	@echo "$(GREEN)✅ Built: $(DIST_DIR)/$(BINARY_LINUX)$(NC)"

.PHONY: build-darwin
build-darwin: ## Build binary for macOS
	@echo "$(BLUE)🍎 Building for macOS...$(NC)"
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_DARWIN) $(MAIN_FILE)
	@echo "$(GREEN)✅ Built: $(DIST_DIR)/$(BINARY_DARWIN)$(NC)"

.PHONY: build-windows
build-windows: ## Build binary for Windows
	@echo "$(BLUE)🪟 Building for Windows...$(NC)"
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_WINDOWS) $(MAIN_FILE)
	@echo "$(GREEN)✅ Built: $(DIST_DIR)/$(BINARY_WINDOWS)$(NC)"

# Dependency management
.PHONY: deps
deps: ## Download and verify dependencies
	@echo "$(BLUE)📦 Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod verify
	@echo "$(GREEN)✅ Dependencies ready!$(NC)"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "$(BLUE)🔄 Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✅ Dependencies updated!$(NC)"

.PHONY: deps-clean
deps-clean: ## Clean module cache
	@echo "$(BLUE)🧹 Cleaning module cache...$(NC)"
	$(GO) clean -modcache
	@echo "$(GREEN)✅ Module cache cleaned!$(NC)"

# Testing targets
.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	$(GO) test -v ./...
	@echo "$(GREEN)✅ All tests passed!$(NC)"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)📊 Running tests with coverage...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) test -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "$(GREEN)✅ Coverage report: $(BUILD_DIR)/coverage.html$(NC)"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "$(BLUE)🏃 Running tests with race detection...$(NC)"
	$(GO) test -race -v ./...
	@echo "$(GREEN)✅ Race tests passed!$(NC)"

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(BLUE)⚡ Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./...

# Code quality targets
.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "$(BLUE)🔍 Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Linting completed!$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)🎨 Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✅ Code formatted!$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)🔍 Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✅ Vet checks passed!$(NC)"

.PHONY: check
check: fmt vet lint test ## Run all code quality checks
	@echo "$(GREEN)✅ All quality checks passed!$(NC)"

# Installation targets
.PHONY: install
install: build ## Install binary to GOPATH/bin
	@echo "$(BLUE)📦 Installing $(BINARY_NAME)...$(NC)"
	$(GO) install $(LDFLAGS) $(MAIN_FILE)
	@echo "$(GREEN)✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

.PHONY: uninstall
uninstall: ## Remove installed binary
	@echo "$(BLUE)🗑️  Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ Uninstalled!$(NC)"

# Release targets
.PHONY: release
release: clean check build-all ## Create a release (clean, check, build-all)
	@echo "$(BLUE)🚀 Creating release...$(NC)"
	@mkdir -p $(DIST_DIR)
	@echo "$(VERSION)" > $(DIST_DIR)/VERSION
	@echo "$(GREEN)✅ Release $(VERSION) ready in $(DIST_DIR)/$(NC)"

.PHONY: package
package: release ## Create release packages
	@echo "$(BLUE)📦 Creating release packages...$(NC)"
	@cd $(DIST_DIR) && \
	tar -czf $(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_LINUX) && \
	tar -czf $(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz $(BINARY_DARWIN) && \
	zip $(BINARY_NAME)_$(VERSION)_windows_amd64.zip $(BINARY_WINDOWS)
	@echo "$(GREEN)✅ Release packages created in $(DIST_DIR)/$(NC)"

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f $(BINARY_NAME)
	@echo "$(GREEN)✅ Cleaned!$(NC)"

.PHONY: clean-all
clean-all: clean deps-clean ## Clean everything including module cache
	@echo "$(GREEN)✅ Everything cleaned!$(NC)"

# Docker targets (optional)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)🐳 Building Docker image...$(NC)"
	docker build -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest
	@echo "$(GREEN)✅ Docker image built: $(PROJECT_NAME):$(VERSION)$(NC)"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(BLUE)🐳 Running Docker container...$(NC)"
	docker run -it --rm $(PROJECT_NAME):latest

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "$(CYAN)OnionCLI Version Information$(NC)"
	@echo "$(CYAN)=============================$(NC)"
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GOVERSION)"
	@echo "OS/Arch:    $(GOOS)/$(GOARCH)"

.PHONY: size
size: build ## Show binary size
	@echo "$(BLUE)📏 Binary size:$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print $$5 " " $$9}'

.PHONY: deps-graph
deps-graph: ## Show dependency graph (requires graphviz)
	@echo "$(BLUE)📊 Generating dependency graph...$(NC)"
	@if command -v dot >/dev/null 2>&1; then \
		$(GO) mod graph | modgraphviz | dot -Tpng -o $(BUILD_DIR)/deps.png; \
		echo "$(GREEN)✅ Dependency graph: $(BUILD_DIR)/deps.png$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  graphviz not found. Install with your package manager.$(NC)"; \
	fi

# Development workflow shortcuts
.PHONY: quick
quick: fmt build ## Quick development cycle (format + build)
	@echo "$(GREEN)✅ Quick build completed!$(NC)"

.PHONY: full
full: clean deps check build ## Full development cycle
	@echo "$(GREEN)✅ Full development cycle completed!$(NC)"

# Make sure intermediate files are not deleted
.PRECIOUS: $(BUILD_DIR)/% $(DIST_DIR)/%
