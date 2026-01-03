#!/bin/bash

# DocTrainerGO Quick Start Script
# This script helps you get started with DocTrainerGO

set -e

echo "🚀 DocTrainerGO Setup"
echo "===================="
echo ""

# Check Go installation
echo "→ Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher from https://go.dev/dl/"
    exit 1
fi
echo "✓ Go found: $(go version)"
echo ""

# Check Ollama installation
echo "→ Checking Ollama installation..."
if ! command -v ollama &> /dev/null; then
    echo "⚠️  Ollama not found. Install from: https://ollama.com"
    echo "   Chat functionality will not work without Ollama."
    read -p "   Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "✓ Ollama found: $(ollama --version 2>&1 | head -n 1)"
fi
echo ""

# Download dependencies
echo "→ Downloading Go dependencies..."
go mod download
echo "✓ Dependencies downloaded"
echo ""

# Download Fuse.js
echo "→ Downloading Fuse.js for search functionality..."
if [ -f "static/fuse.min.js" ] && grep -q "Fuse.js" static/fuse.min.js 2>/dev/null; then
    echo "✓ Fuse.js already downloaded"
else
    curl -L https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js
    echo "✓ Fuse.js downloaded"
fi
echo ""

# Create input directory
mkdir -p input
echo "✓ Created input directory"
echo ""

# Check for PDF files
PDF_COUNT=$(ls -1 input/*.pdf 2>/dev/null | wc -l)
if [ "$PDF_COUNT" -eq 0 ]; then
    echo "⚠️  No PDF files found in input/ directory"
    echo "   Please place a PDF file in the input/ directory:"
    echo "   cp /path/to/your/document.pdf input/"
    echo ""
fi

echo "✅ Setup complete!"
echo ""
echo "📖 Next steps:"
echo ""
echo "1. Place your PDF in the input/ directory:"
echo "   cp /path/to/your/document.pdf input/"
echo ""
echo "2. Process the PDF:"
echo "   go run cmd/main.go -pdf=input/document.pdf"
echo ""
echo "3. Start Ollama (in a separate terminal):"
echo "   ollama run llama3.1"
echo ""
echo "4. Start the web server:"
echo "   go run cmd/main.go -serve"
echo ""
echo "5. Open http://localhost:8080 in your browser"
echo ""
