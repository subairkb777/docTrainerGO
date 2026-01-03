#!/bin/bash

# Quick diagnostic script for DocTrainerGO issues

echo "🔍 DocTrainerGO Diagnostic Tool"
echo "================================"
echo ""

# Check 1: Ollama
echo "1️⃣  Checking Ollama..."
if command -v ollama &> /dev/null; then
    echo "   ✓ Ollama installed"
    
    # Check if Ollama is running
    if curl -s http://localhost:11434/api/tags &> /dev/null; then
        echo "   ✓ Ollama is running on port 11434"
        
        # List models
        echo "   📦 Available models:"
        curl -s http://localhost:11434/api/tags | grep -o '"name":"[^"]*"' | cut -d'"' -f4 | sed 's/^/      - /'
    else
        echo "   ❌ Ollama is NOT running"
        echo "      Solution: Run 'ollama serve' in a separate terminal"
        echo "                Then run 'ollama run llama3.1' in another terminal"
    fi
else
    echo "   ❌ Ollama not installed"
    echo "      Install from: https://ollama.com"
fi
echo ""

# Check 2: pdfimages for image extraction
echo "2️⃣  Checking image extraction tools..."
if command -v pdfimages &> /dev/null; then
    echo "   ✓ pdfimages installed"
    echo "      Can extract images from PDFs"
else
    echo "   ⚠️  pdfimages not found (optional)"
    echo "      Install with: brew install poppler"
    echo "      This enables actual image extraction from PDFs"
fi
echo ""

# Check 3: Server status
echo "3️⃣  Checking server..."
if curl -s http://localhost:8080 &> /dev/null; then
    echo "   ✓ Server is running on port 8080"
    
    # Test chat endpoint
    CHAT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/chat \
        -H "Content-Type: application/json" \
        -d '{"prompt":"test"}')
    
    if echo "$CHAT_RESPONSE" | grep -q "error"; then
        echo "   ❌ Chat endpoint has errors"
        echo "      Response: $CHAT_RESPONSE"
    else
        echo "   ✓ Chat endpoint is working"
    fi
else
    echo "   ⚠️  Server is not running"
    echo "      Start with: make serve"
fi
echo ""

# Check 4: PDF and generated files
echo "4️⃣  Checking files..."
if [ -d "input" ]; then
    PDF_COUNT=$(ls -1 input/*.pdf 2>/dev/null | wc -l)
    echo "   📄 PDFs in input/: $PDF_COUNT"
    if [ "$PDF_COUNT" -gt 0 ]; then
        ls -1 input/*.pdf | sed 's/^/      - /'
    fi
else
    echo "   ⚠️  input/ directory not found"
fi

if [ -d "docs" ]; then
    echo "   ✓ docs/ directory exists"
    if [ -f "docs/index.html" ]; then
        echo "      ✓ docs/index.html generated"
    else
        echo "      ⚠️  docs/index.html not found (run: make process PDF=input/yourfile.pdf)"
    fi
    
    if [ -d "docs/images" ]; then
        IMG_COUNT=$(ls -1 docs/images/*.{png,jpg,jpeg} 2>/dev/null | wc -l)
        echo "      📸 Images in docs/images/: $IMG_COUNT"
    else
        echo "      ⚠️  docs/images/ directory not found"
    fi
else
    echo "   ⚠️  docs/ directory not found"
    echo "      Run: make process PDF=input/yourfile.pdf"
fi
echo ""

# Check 5: Dependencies
echo "5️⃣  Checking dependencies..."
if [ -f "static/fuse.min.js" ]; then
    SIZE=$(wc -c < static/fuse.min.js)
    if [ "$SIZE" -gt 1000 ]; then
        echo "   ✓ Fuse.js downloaded ($SIZE bytes)"
    else
        echo "   ⚠️  Fuse.js seems incomplete"
        echo "      Download: make download-fuse"
    fi
else
    echo "   ❌ Fuse.js not found"
    echo "      Download: make download-fuse"
fi
echo ""

# Summary and recommendations
echo "================================"
echo "📋 Summary & Next Steps:"
echo ""

# Determine what needs to be done
NEEDS_OLLAMA=false
NEEDS_SERVER=false
NEEDS_PDF=false

if ! curl -s http://localhost:11434/api/tags &> /dev/null; then
    NEEDS_OLLAMA=true
fi

if ! curl -s http://localhost:8080 &> /dev/null; then
    NEEDS_SERVER=true
fi

if [ ! -f "docs/index.html" ]; then
    NEEDS_PDF=true
fi

if [ "$NEEDS_OLLAMA" = true ]; then
    echo "🔴 Start Ollama:"
    echo "   Terminal 1: ollama serve"
    echo "   Terminal 2: ollama run llama3.1"
    echo ""
fi

if [ "$NEEDS_PDF" = true ]; then
    echo "🔴 Process your PDF:"
    echo "   make process PDF=input/your-file.pdf"
    echo ""
fi

if [ "$NEEDS_SERVER" = true ]; then
    echo "🔴 Start the server:"
    echo "   make serve"
    echo ""
fi

if [ "$NEEDS_OLLAMA" = false ] && [ "$NEEDS_SERVER" = false ] && [ "$NEEDS_PDF" = false ]; then
    echo "✅ Everything looks good!"
    echo "   Visit: http://localhost:8080"
else
    echo "Follow the steps above to fix any issues."
fi

echo ""
echo "For more help, see: TROUBLESHOOTING.md"
