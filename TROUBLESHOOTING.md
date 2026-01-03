# 🔧 Troubleshooting Guide

## Issue 1: Chat Error - "Failed to get response from AI"

### **Quick Fix:**

```bash
# Terminal 1: Start Ollama server
ollama serve

# Terminal 2: Run the model
ollama run llama3.1

# Terminal 3: Start your app
cd /Users/subairkb/Desktop/workspace/LLM/docTrainerGO
make serve
```

### **Detailed Diagnosis:**

**Step 1: Check if Ollama is installed**
```bash
ollama --version
# If not found, install from: https://ollama.com
```

**Step 2: Check if model is downloaded**
```bash
ollama list
# Should show llama3.1 or other models

# If not, download it:
ollama pull llama3.1
```

**Step 3: Test Ollama directly**
```bash
# Test the API
curl http://localhost:11434/api/tags

# Should return JSON like:
# {"models":[{"name":"llama3.1:latest",...}]}
```

**Step 4: Test chat API**
```bash
# With server running, test the chat endpoint
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Hello, can you help me?"}'

# Should return: {"answer": "..."}
```

**Step 5: Check server logs**
```bash
# When you start the server, you should see:
# ✓ Connected to Ollama

# If you see:
# ⚠️  Ollama health check failed
# Then Ollama is not running properly
```

### **Common Causes:**

1. **Ollama not running**
   ```bash
   # Solution: Start Ollama
   ollama serve &
   ollama run llama3.1
   ```

2. **Wrong port**
   ```bash
   # Check if Ollama is on different port
   lsof -i :11434
   
   # Start server with correct Ollama URL
   go run cmd/main.go -serve -ollama=http://localhost:11434
   ```

3. **Model not loaded**
   ```bash
   # Load the model
   ollama run llama3.1
   # Keep this terminal open
   ```

4. **Firewall blocking**
   ```bash
   # Test connection
   telnet localhost 11434
   ```

---

## Issue 2: Images Not Loading from PDF

### **Root Cause:**

The `ledongthuc/pdf` library has **very limited image extraction capabilities**. It can detect image references but cannot easily extract the actual image data.

### **Current Status:**

The parser finds image references but doesn't extract actual images. You'll see messages like:
```
Found image reference on page 2 (saved as page2_img1.png)
```

But the files aren't actually created in `docs/images/`.

### **Solution Options:**

#### **Option 1: Use External Tool (Recommended for Production)**

Use `pdfimages` command-line tool to extract images:

```bash
# Install pdfimages (part of poppler-utils)
# macOS:
brew install poppler

# Linux:
sudo apt-get install poppler-utils

# Extract images from your PDF
pdfimages -png input/user_guide.pdf docs/images/img

# This creates: docs/images/img-000.png, img-001.png, etc.

# Then process the PDF
make process PDF=input/user_guide.pdf
```

#### **Option 2: Manual Image Placement**

If you know which images should be in the documentation:

```bash
# 1. Create the images directory
mkdir -p docs/images

# 2. Manually extract images from PDF using Preview/Adobe/etc.
# Save them as: image_1.png, image_2.png, etc.

# 3. Place them in docs/images/
cp ~/Desktop/extracted_image1.png docs/images/image_1.png
cp ~/Desktop/extracted_image2.png docs/images/image_2.png

# 4. Process the PDF
make process PDF=input/user_guide.pdf
```

#### **Option 3: Add placeholder images**

Create placeholder images to test the layout:

```bash
# Install ImageMagick
brew install imagemagick

# Create test images
for i in {1..5}; do
  convert -size 800x600 -gravity center \
    -background lightgray -fill darkgray \
    -pointsize 72 label:"Image $i" \
    docs/images/image_$i.png
done

# Refresh the page
```

#### **Option 4: Use a Different PDF Library (Advanced)**

For production use with actual image extraction, consider using:

1. **pdfcpu** (Go library with better image support)
2. **External tool integration** (call pdfimages from Go)
3. **Python pdfplumber** (via API)

### **Quick Workaround Script:**

Create a helper script `extract_images.sh`:

```bash
#!/bin/bash
# Extract images from PDF using pdfimages

PDF_FILE="$1"
OUTPUT_DIR="docs/images"

if [ -z "$PDF_FILE" ]; then
    echo "Usage: ./extract_images.sh input/your_file.pdf"
    exit 1
fi

# Check if pdfimages is installed
if ! command -v pdfimages &> /dev/null; then
    echo "❌ pdfimages not found"
    echo "Install with: brew install poppler"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Extract images
echo "→ Extracting images from PDF..."
pdfimages -png "$PDF_FILE" "$OUTPUT_DIR/img"

# Count extracted images
COUNT=$(ls -1 "$OUTPUT_DIR"/img-*.png 2>/dev/null | wc -l)
echo "✓ Extracted $COUNT images to $OUTPUT_DIR/"

# List extracted images
ls -lh "$OUTPUT_DIR"/img-*.png 2>/dev/null
```

**Usage:**

```bash
chmod +x extract_images.sh
./extract_images.sh input/user_guide.pdf
make process PDF=input/user_guide.pdf
```

---

## Issue 3: Images in HTML but Not Displaying

If images are referenced but showing as broken:

### **Check 1: Verify image paths**

```bash
# Check what's in the docs/images directory
ls -la docs/images/

# Check HTML references
grep -n "img src" docs/index.html
```

### **Check 2: Verify server is serving images**

```bash
# With server running, test image URL
curl -I http://localhost:8080/docs/images/image_1.png

# Should return: HTTP/1.1 200 OK
```

### **Check 3: Check browser console**

1. Open http://localhost:8080
2. Press F12 (Developer Tools)
3. Go to Console tab
4. Look for 404 errors on image files

### **Fix: Update image paths**

The images should be accessible at `/docs/images/filename.png`

Check [templates/page.html](templates/page.html) to ensure correct path:

```html
<!-- Should be: -->
<img src="/docs/images/{{.}}" alt="{{.}}" loading="lazy">
```

---

## Complete Workflow for PDF with Images

### **Method 1: Using pdfimages (Recommended)**

```bash
# 1. Install pdfimages
brew install poppler

# 2. Extract images first
pdfimages -png input/user_guide.pdf docs/images/img

# 3. Process the PDF
make process PDF=input/user_guide.pdf

# 4. Start Ollama (new terminal)
ollama run llama3.1

# 5. Start server
make serve

# 6. Open browser
open http://localhost:8080
```

### **Method 2: Without Image Extraction**

```bash
# Just process text content
make process PDF=input/user_guide.pdf

# Images will be noted but not displayed
# Documentation will still work for text content
```

---

## Testing Both Fixes Together

```bash
# Terminal 1: Start Ollama
ollama serve
# Keep running

# Terminal 2: Load model
ollama run llama3.1
# Keep running

# Terminal 3: Extract images (if pdfimages available)
pdfimages -png input/user_guide.pdf docs/images/img

# Terminal 3: Process PDF
make process PDF=input/user_guide.pdf

# Terminal 3: Start server
make serve

# Should see:
# ✓ Connected to Ollama
# 🚀 Server running at http://localhost:8080
```

**Then test:**
1. Navigate to http://localhost:8080
2. Verify text content is there
3. Check if images display
4. Click chat widget (bottom-right)
5. Type a question and send
6. Should get AI response

---

## Quick Diagnostic Commands

```bash
# Check Ollama status
curl http://localhost:11434/api/tags

# Check if server is running
curl http://localhost:8080

# Check if images exist
ls -la docs/images/

# Check server logs
# Look at the terminal where you ran 'make serve'

# Test chat endpoint
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "test"}'

# Check for pdfimages
which pdfimages
```

---

## Summary: Your Current Issues

### **Chat Error Fix:**
```bash
# Run these in order in separate terminals:
# Terminal 1:
ollama serve

# Terminal 2:
ollama run llama3.1

# Terminal 3:
make serve
```

### **Image Loading Fix:**
```bash
# Option A: Install pdfimages and extract
brew install poppler
pdfimages -png input/user_guide.pdf docs/images/img
make process PDF=input/user_guide.pdf

# Option B: Accept text-only for now
# The PDF library we're using has limited image support
# Text content will still work perfectly
```

---

Need more help? Check:
- [README.md](README.md) - Full documentation
- [EXAMPLES.md](EXAMPLES.md) - More examples
- Server logs in your terminal
