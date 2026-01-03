# 🚀 Quick Start Guide

Get DocTrainerGO up and running in 5 minutes!

## Prerequisites

- **Go 1.21+** → [Download](https://go.dev/dl/)
- **Ollama** → [Download](https://ollama.com/) (optional, for chat)

## Installation

### Option 1: Automated Setup (Recommended)

```bash
# Navigate to project directory
cd docTrainerGO

# Run setup script
./setup.sh

# Or use Makefile
make setup
```

### Option 2: Manual Setup

```bash
# 1. Download Go dependencies
go mod download

# 2. Download Fuse.js
curl -L https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js

# 3. Verify installation
./verify.sh
```

## Usage

### Step 1: Prepare PDF

```bash
# Copy your PDF to input directory
cp ~/Documents/manual.pdf input/
```

### Step 2: Process PDF

```bash
# Option A: Using Makefile (recommended)
make process PDF=input/manual.pdf

# Option B: Using go run
go run cmd/main.go -pdf=input/manual.pdf
```

**Output:**
```
Processing PDF: input/manual.pdf
→ Parsing PDF and extracting content...
  Found 12 sections
→ Generating HTML pages...
Generated: docs/index.html
→ Creating search index...
✓ PDF processing complete!
```

### Step 3: Start Ollama (Optional)

Open a **new terminal** and run:

```bash
# Download model (first time only)
ollama pull llama3.1

# Start Ollama
ollama run llama3.1
```

Keep this terminal open while using the chat feature.

### Step 4: Start Web Server

Back in your main terminal:

```bash
# Option A: Using Makefile
make serve

# Option B: Using go run
go run cmd/main.go -serve

# Option C: Custom port
make serve PORT=3000
```

**Output:**
```
✓ Connected to Ollama
🚀 Server running at http://localhost:8080
   Press Ctrl+C to stop
```

### Step 5: Open Browser

Visit: **http://localhost:8080**

You should see:
- ✅ Your PDF content as a beautiful website
- ✅ Searchable documentation
- ✅ Sidebar navigation
- ✅ AI chat assistant (bottom-right)

## Common Commands

```bash
# Verify project structure
./verify.sh

# Process PDF
make process PDF=input/document.pdf

# Serve with custom port
make serve PORT=3000

# Build binary
make build
./doctrainer -serve

# Clean generated files
make clean

# View all commands
make help
```

## Troubleshooting

### "Go not found"
```bash
# Install Go from: https://go.dev/dl/
# Then verify:
go version
```

### "Ollama not accessible"
```bash
# Check if Ollama is running:
ps aux | grep ollama

# Start Ollama:
ollama serve

# In another terminal:
ollama run llama3.1
```

### "Search not working"
```bash
# Download Fuse.js:
make download-fuse

# Or manually:
curl -L https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js
```

### "Port already in use"
```bash
# Use different port:
make serve PORT=3000

# Or kill process on port 8080:
lsof -ti:8080 | xargs kill
```

## Next Steps

- 📖 Read the full [README.md](README.md)
- 💡 Check [EXAMPLES.md](EXAMPLES.md) for advanced usage
- 🎨 Customize `static/style.css` for your branding
- 🔧 Modify `templates/page.html` for layout changes

## Architecture Overview

```
PDF File (input/)
    ↓
Go Parser (internal/pdf/)
    ↓ extracts
Text + Images
    ↓
HTML Generator (internal/generator/)
    ↓ creates
Static Website (docs/)
    ↓
Web Server (cmd/main.go)
    ↓ serves
Browser ← → Chat API → Ollama LLM
```

## File Structure

```
docTrainerGO/
├── input/           # Place PDFs here
├── docs/            # Generated website (auto-created)
├── static/          # CSS, JS, libraries
├── templates/       # HTML templates
├── internal/        # Go packages
│   ├── pdf/        # PDF parsing
│   ├── generator/  # HTML generation
│   ├── search/     # Search index
│   └── chat/       # Ollama integration
└── cmd/            # Main application
```

## Tips

1. **Large PDFs**: Processing may take 1-2 minutes for files >50MB
2. **Images**: Extracted images are saved in `docs/images/`
3. **Search**: Press `Ctrl+K` (or `Cmd+K` on Mac) to focus search
4. **Chat**: Works best with llama3.1, mistral, or codellama models
5. **Offline**: Everything runs locally, no internet required after setup

## Support

- 🐛 Found a bug? Check browser console and server logs
- 💬 Need help? Open an issue on GitHub
- 📧 Questions? Read the full README.md

---

**Ready to go!** 🎉

Run `./setup.sh` to begin, or `make help` for all commands.
