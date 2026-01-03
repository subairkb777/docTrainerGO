#!/bin/bash

# Project verification script
# Checks if all required files are in place

echo "🔍 DocTrainerGO Project Verification"
echo "===================================="
echo ""

ERRORS=0
WARNINGS=0

# Function to check file exists
check_file() {
    if [ -f "$1" ]; then
        echo "✓ $1"
    else
        echo "❌ MISSING: $1"
        ((ERRORS++))
    fi
}

# Function to check directory exists
check_dir() {
    if [ -d "$1" ]; then
        echo "✓ $1/"
    else
        echo "❌ MISSING: $1/"
        ((ERRORS++))
    fi
}

echo "Checking core files..."
check_file "go.mod"
check_file "README.md"
check_file "Makefile"
check_file "setup.sh"
echo ""

echo "Checking source code..."
check_file "cmd/main.go"
check_file "internal/pdf/parser.go"
check_file "internal/generator/html.go"
check_file "internal/search/index.go"
check_file "internal/chat/ollama.go"
echo ""

echo "Checking templates..."
check_file "templates/page.html"
echo ""

echo "Checking static files..."
check_file "static/style.css"
check_file "static/script.js"
check_file "static/fuse.min.js"
echo ""

echo "Checking directories..."
check_dir "input"
check_dir "static"
check_dir "templates"
check_dir "cmd"
check_dir "internal"
echo ""

# Check Go installation
echo "Checking dependencies..."
if command -v go &> /dev/null; then
    echo "✓ Go installed: $(go version | awk '{print $3}')"
else
    echo "❌ Go not installed"
    ((ERRORS++))
fi

# Check Ollama installation
if command -v ollama &> /dev/null; then
    echo "✓ Ollama installed"
else
    echo "⚠️  Ollama not installed (optional for chat)"
    ((WARNINGS++))
fi

# Check for fuse.js content
if grep -q "Fuse.js not loaded" static/fuse.min.js 2>/dev/null; then
    echo "⚠️  Fuse.js needs to be downloaded"
    echo "   Run: make download-fuse"
    ((WARNINGS++))
else
    echo "✓ Fuse.js appears to be installed"
fi

echo ""
echo "===================================="
echo "Summary:"
echo "  Errors: $ERRORS"
echo "  Warnings: $WARNINGS"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo "✅ Project structure is complete!"
    echo ""
    echo "Next steps:"
    echo "  1. Download Fuse.js: make download-fuse"
    echo "  2. Install Go deps: make deps"
    echo "  3. Place PDF in input/"
    echo "  4. Process: make process PDF=input/your-file.pdf"
    echo "  5. Serve: make serve"
    exit 0
else
    echo "❌ Project has missing files"
    exit 1
fi
