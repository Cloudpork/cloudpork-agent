# CloudPork Agent Makefile

.PHONY: build test clean install dev fmt lint run help

# Variables
BINARY_NAME := cloudpork
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-s -w -X github.com/carsor007/cloudpork-agent/cmd.version=$(VERSION) -X github.com/carsor007/cloudpork-agent/cmd.commit=$(COMMIT) -X github.com/carsor007/cloudpork-agent/cmd.date=$(DATE)"

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building CloudPork Agent..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go
	@echo "✅ Build complete: bin/$(BINARY_NAME)"

# Build for all platforms
build-all:
	@echo "🔨 Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 main.go
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "✅ Multi-platform build complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ dist/ coverage.out coverage.html
	@go clean
	@echo "✅ Clean complete"

# Install to system
install: build
	@echo "📦 Installing CloudPork Agent..."
	@sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

# Development mode - build and run
dev: build
	@echo "🚀 Running in development mode..."
	@./bin/$(BINARY_NAME) --help

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run
	@echo "✅ Lint complete"

# Run the binary
run: build
	@./bin/$(BINARY_NAME) $(ARGS)

# Initialize development environment
init:
	@echo "🛠️  Initializing development environment..."
	@go mod tidy
	@go mod download
	@echo "✅ Development environment ready"

# Release using goreleaser
release:
	@echo "🚀 Creating release..."
	@goreleaser release --clean

# Snapshot release (for testing)
snapshot:
	@echo "📸 Creating snapshot release..."
	@goreleaser release --snapshot --clean

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t cloudpork/agent:latest .
	@echo "✅ Docker image built: cloudpork/agent:latest"

# Show help
help:
	@echo "CloudPork Agent Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build         Build the binary"
	@echo "  build-all     Build for all platforms"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean         Clean build artifacts"
	@echo "  install       Install to system (/usr/local/bin)"
	@echo "  dev           Build and show help (development mode)"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code"
	@echo "  run           Build and run with ARGS"
	@echo "  init          Initialize development environment"
	@echo "  release       Create release with goreleaser"
	@echo "  snapshot      Create snapshot release"
	@echo "  docker-build  Build Docker image"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run ARGS='analyze --help'"
	@echo "  make run ARGS='auth login'"
	@echo "  make run ARGS='analyze /path/to/project'"