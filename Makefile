# Makefile for DocTrainerGO

.PHONY: help setup deps download-fuse process serve clean build test

# Default target
help:
	@echo "DocTrainerGO - Makefile Commands"
	@echo "================================"
	@echo ""
	@echo "Setup & Installation:"
	@echo "  make setup          - Run initial setup (deps + fuse.js)"
	@echo "  make deps           - Download Go dependencies"
	@echo "  make download-fuse  - Download Fuse.js library"
	@echo ""
	@echo "Usage:"
	@echo "  make process PDF=input/doc.pdf  - Process a PDF file"
	@echo "  make serve          - Start the web server"
	@echo "  make serve PORT=3000 - Start server on custom port"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build the binary"
	@echo "  make clean          - Clean generated files"
	@echo "  make test           - Run tests"
	@echo ""
	@echo "Examples:"
	@echo "  make setup"
	@echo "  make process PDF=input/manual.pdf"
	@echo "  make serve"

# Initial setup
setup: deps download-fuse
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Place PDF in input/ directory"
	@echo "  2. make process PDF=input/your-file.pdf"
	@echo "  3. ollama run llama3.1  (in another terminal)"
	@echo "  4. make serve"

# Download Go dependencies
deps:
	@echo "→ Downloading Go dependencies..."
	@go mod download
	@go mod verify
	@echo "✓ Dependencies installed"

# Download Fuse.js
download-fuse:
	@echo "→ Downloading Fuse.js..."
	@mkdir -p static
	@curl -sL https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js
	@echo "✓ Fuse.js downloaded"

# Process PDF
process:
	@if [ -z "$(PDF)" ]; then \
		echo "Error: PDF parameter required"; \
		echo "Usage: make process PDF=input/document.pdf"; \
		exit 1; \
	fi
	@echo "→ Processing PDF: $(PDF)"
	@go run cmd/main.go -pdf=$(PDF)

# Start web server
serve:
	@PORT=$${PORT:-8080}; \
	echo "→ Starting server on port $$PORT..."; \
	go run cmd/main.go -serve -port=$$PORT

# Build binary
build:
	@echo "→ Building binary..."
	@go build -o doctrainer cmd/main.go
	@echo "✓ Binary created: ./doctrainer"

# Clean generated files
clean:
	@echo "🧹 Cleaning generated files..."
	@echo ""
	@if [ -d "docs/" ]; then \
		echo "  Removing docs/ directory..."; \
		rm -rf docs/; \
		echo "  ✓ Removed docs/ (HTML, data, images, search index)"; \
	else \
		echo "  ⚠️  docs/ directory not found (already clean)"; \
	fi
	@if [ -f "main" ]; then \
		rm -f main; \
		echo "  ✓ Removed binary: main"; \
	fi
	@if [ -f "doctrainer" ]; then \
		rm -f doctrainer; \
		echo "  ✓ Removed binary: doctrainer"; \
	fi
	@echo ""
	@echo "✅ Cleanup complete!"
	@echo ""
	@echo "To regenerate:"
	@echo "  make process PDF=input/your-file.pdf"

# Run tests
test:
	@echo "→ Running tests..."
	@go test ./...

# Development server with auto-reload (requires air)
dev:
	@if ! command -v air > /dev/null; then \
		echo "Installing air for auto-reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

# Check Ollama status
check-ollama:
	@echo "→ Checking Ollama..."
	@if command -v ollama > /dev/null; then \
		echo "✓ Ollama installed"; \
		ollama list 2>/dev/null || echo "⚠️  No models installed. Run: ollama pull llama3.1"; \
	else \
		echo "❌ Ollama not installed. Visit: https://ollama.com"; \
	fi

# Full workflow example
example: setup
	@echo ""
	@echo "→ Running example workflow..."
	@echo "  Please ensure you have a PDF file in input/"
	@echo "  Then run: make process PDF=input/your-file.pdf"
	@echo "  Then run: make serve"
