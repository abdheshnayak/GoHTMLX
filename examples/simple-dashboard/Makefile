# Simple Dashboard Makefile

GOHTMLX := ../../../GoHTMLX/bin/gohtmlx
CONFIG := gohtmlx.config.yaml

.PHONY: build watch serve dev clean help

# Default target
all: build

# Build components
build:
	@echo "🔨 Building components..."
	$(GOHTMLX) build --config $(CONFIG)
	@echo "✅ Components built successfully!"

# Watch for changes
watch:
	@echo "👀 Watching for changes..."
	$(GOHTMLX) watch --config $(CONFIG)

# Run the server
serve:
	@echo "🚀 Starting server..."
	go run main.go

# Development mode (build then serve)
dev: build serve

# Clean generated files
clean:
	@echo "🧹 Cleaning..."
	rm -rf dist/
	@echo "✅ Clean complete!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	@echo "✅ Dependencies installed!"

# Show help
help:
	@echo "🚀 GoHTMLX Simple Dashboard"
	@echo ""
	@echo "Available commands:"
	@echo "  build  - Build components"
	@echo "  watch  - Watch for changes"
	@echo "  serve  - Run the server"
	@echo "  dev    - Build and serve"
	@echo "  clean  - Clean generated files"
	@echo "  deps   - Install dependencies"
	@echo "  help   - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "  make deps  # Install dependencies"
	@echo "  make dev   # Build and run"
