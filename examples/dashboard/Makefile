# Dashboard Example Makefile

# Configuration
GOHTMLX := ../../bin/gohtmlx
CONFIG := gohtmlx.config.yaml
SRC_DIR := src
DIST_DIR := dist

.PHONY: build watch serve dev clean help

# Default target
all: build

# Build components
build:
	@echo "🔨 Building dashboard components..."
	$(GOHTMLX) build --config $(CONFIG)
	@echo "✅ Components built successfully!"

# Watch for changes and rebuild
watch:
	@echo "👀 Watching for changes..."
	$(GOHTMLX) watch --config $(CONFIG)

# Run the server
serve:
	@echo "🚀 Starting dashboard server..."
	go run main.go

# Development workflow (watch + serve in parallel)
dev:
	@echo "🔥 Starting development environment..."
	@echo "📦 Building components first..."
	@$(MAKE) build
	@echo "🚀 Starting server..."
	@$(MAKE) serve

# Clean generated files
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf $(DIST_DIR)
	@echo "✅ Clean complete!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed!"

# Run with hot reload (requires separate terminals)
dev-watch:
	@echo "👀 Starting file watcher (run 'make serve' in another terminal)..."
	@$(MAKE) watch

# Show help
help:
	@echo "🚀 GoHTMLX Dashboard Example"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        - Build dashboard components"
	@echo "  watch        - Watch for changes and rebuild"
	@echo "  serve        - Run the dashboard server"
	@echo "  dev          - Build and run server"
	@echo "  dev-watch    - Watch for changes (run in separate terminal)"
	@echo "  clean        - Clean generated files"
	@echo "  deps         - Install Go dependencies"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Development workflow:"
	@echo "  Terminal 1: make dev-watch"
	@echo "  Terminal 2: make serve"
	@echo ""
	@echo "Or simply: make dev"
