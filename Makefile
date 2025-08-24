# GoHTMLX Makefile
# Build configuration
BINARY_NAME := gohtmlx
BUILD_DIR := ./bin
CMD_DIR := ./cmd/gohtmlx
EXAMPLE_SRC := example/src
EXAMPLE_DIST := example/dist

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)"

# Go commands
GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOMOD := $(GO) mod
GOFMT := $(GO) fmt

.PHONY: all build install test test-coverage lint fmt clean docs deps release help \
	dev dev-full dev-server serve watch \
	build-example build-example-new \
	css-build css-watch init version check

# Default target
all: clean fmt lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

# Install globally
install:
	@echo "Installing $(BINARY_NAME) globally..."
	$(GO) install $(LDFLAGS) $(CMD_DIR)
	@echo "$(BINARY_NAME) installed successfully!"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(EXAMPLE_DIST)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@mkdir -p docs
	$(GO) doc -all ./pkg/element > docs/element.md
	@echo "Documentation generated in docs/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			ext=""; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-$$os-$$arch$$ext $(CMD_DIR); \
		done; \
	done
	@echo "Release binaries built in $(BUILD_DIR)/release/"

# Build example components
build-example: build
	@echo "Building example components..."
	@$(BUILD_DIR)/$(BINARY_NAME) build --src=$(EXAMPLE_SRC) --dist=$(EXAMPLE_DIST)
	@echo "Example components built successfully!"

# Build example components with verbose output
build-example-new: build
	@echo "Building example components with verbose output..."
	@$(BUILD_DIR)/$(BINARY_NAME) build --src=$(EXAMPLE_SRC) --dist=$(EXAMPLE_DIST) --verbose
	@echo "Example components built successfully!"

# Watch for changes and rebuild example
watch: build
	@echo "Starting file watcher..."
	@$(BUILD_DIR)/$(BINARY_NAME) watch --src=$(EXAMPLE_SRC) --dist=$(EXAMPLE_DIST) --verbose

# Run example server
serve:
	@echo "Starting example server..."
	@cd example && go run .

# Development workflow (watch + serve)
dev:
	@echo "Starting development workflow..."
	@echo "This will start the watcher in the background and then the server"
	@$(MAKE) watch &
	@sleep 2
	@$(MAKE) serve

# Full development workflow
dev-full: build-example-new
	@echo "Starting full development environment..."
	@$(MAKE) serve

# Development server with auto-restart for Go files
dev-server:
	@echo "Starting development server with auto-restart..."
	@cd example && go run ../cmd/dev-server/main.go . go run .

# Build CSS using Tailwind
css-build:
	@echo "Building CSS..."
	@cd example && tailwindcss -i ./src/input.css -o ./dist/static/main.css

# Watch CSS changes using Tailwind
css-watch:
	@echo "Watching CSS changes..."
	@cd example && tailwindcss -i ./src/input.css -o ./dist/static/main.css --watch

# Show version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

# Initialize a new GoHTMLX project
init:
	@echo "Initializing new GoHTMLX project..."
	@mkdir -p src components
	@if [ -f gohtmlx.config.yaml ]; then \
		cp gohtmlx.config.yaml ./gohtmlx.config.yaml; \
	else \
		echo "# GoHTMLX Config" > gohtmlx.config.yaml; \
	fi
	@echo "GoHTMLX project initialized!"

# Check code quality
check: fmt lint test
	@echo "All quality checks passed!"

# Show help
help:
	@echo "GoHTMLX Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build Targets:"
	@echo "  all              - Run clean, fmt, lint, test, and build"
	@echo "  build            - Build the binary"
	@echo "  install          - Install globally"
	@echo "  release          - Build for multiple platforms"
	@echo "  clean            - Clean build artifacts"
	@echo ""
	@echo "Development Targets:"
	@echo "  dev              - Run development workflow (watch + serve)"
	@echo "  dev-full         - Build components and start server"
	@echo "  dev-server       - Start development server with auto-restart"
	@echo "  watch            - Watch for changes and rebuild example"
	@echo "  serve            - Run example server"
	@echo ""
	@echo "Example Targets:"
	@echo "  build-example    - Build example components"
	@echo "  build-example-new - Build example components with verbose output"
	@echo "  css-build        - Build CSS using Tailwind"
	@echo "  css-watch        - Watch CSS changes using Tailwind"
	@echo ""
	@echo "Quality Targets:"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  check            - Run all quality checks"
	@echo ""
	@echo "Utility Targets:"
	@echo "  docs             - Generate documentation"
	@echo "  deps             - Download dependencies"
	@echo "  init             - Initialize a new GoHTMLX project"
	@echo "  version          - Show version information"
	@echo "  help             - Show this help"
