# DocTrainerGO

🚀 **Convert PDF/Markdown documents into fully searchable, interactive documentation websites with AI-powered chat assistance.**

DocTrainerGO is a modular Go-based solution that transforms PDFs or Markdown files into beautiful, responsive documentation websites featuring:
- ✅ **Dual Input Support**: Process PDFs or Markdown files
- ✅ **Sidebar Navigation**: Collapsible sections with smooth scrolling
- ✅ **Real-time Search**: Fuzzy search powered by Fuse.js
- ✅ **Image Support**: Auto-extract from PDFs or copy from Markdown
- ✅ **AI Chat Assistant**: Context-aware responses using local Ollama LLM
- ✅ **Responsive Design**: Mobile & desktop optimized
- ✅ **Privacy First**: Runs entirely on your laptop (no cloud dependencies)
- ✅ **Modular Architecture**: Clean, maintainable codebase (refactored from 396 → 79 lines in main.go!)

---

## 📋 Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Deployment](#deployment)

---

## ✨ Features

### Core Capabilities

#### 1. Dual Input Support
- **PDF Processing**: Extract text, structure, and images from PDF files using `github.com/ledongthuc/pdf` and Poppler
- **Markdown Processing**: Parse standard Markdown with front matter, code blocks, tables, and links
- **Auto-Discovery**: Automatically find all `.md` files in a directory
- **Smart Parsing**: Detects heading hierarchy (H1-H6) and document structure
- **Config-Driven**: Switch between PDF/Markdown via `config.yaml`

#### 2. AI-Powered Chat Assistant
- **Context-Aware Responses**: Answers based on your documentation content
- **Local LLM**: Runs on your machine using Ollama (llama3.2, mistral, codellama)
- **Formatted Responses**: Supports lists, code blocks, **bold**, *italic*, inline code
- **Real-time Interaction**: Fast responses with streaming support
- **Floating Widget**: Non-intrusive chat interface, minimizable and accessible
- **Optional**: Can be disabled via config for lightweight deployments
- **Privacy-First**: No data sent to cloud services

#### 3. Advanced Search
- **Fuzzy Matching**: Find content even with typos using Fuse.js 6.6.2
- **Weighted Results**: Prioritizes headings over content
- **Real-time Results**: Instant search as you type
- **Keyboard Shortcuts**: `Ctrl/Cmd + K` to focus search
- **Client-Side**: No server queries needed

#### 4. Organized Data Structure
Content stored in maintainable JSON format:
- **content.json**: Master file with all sections (single source of truth)
- **sections/*.json**: Individual section files for easy editing and version control
- **search-index.json**: Optimized for fast search queries
- **API-Ready**: Use generated data in other applications

#### 5. Responsive UI
- **Mobile-First Design**: Optimized for all screen sizes
- **Sidebar Navigation**: Collapsible with searchable sections
- **Dynamic Loading**: Lightweight HTML (3.7KB) loads content from JSON (96% size reduction!)
- **Lazy Image Loading**: Images load as needed for fast page loads
- **Dark Mode Ready**: Easy to customize with CSS variables
- **Smooth Navigation**: Scroll-to-heading with active section highlighting

#### 6. Image Support
- **Automatic Extraction**: From PDFs using pdfimages (Poppler)
- **Markdown Images**: Copy and reference images from markdown directories
- **Lazy Loading**: Performance optimized
- **Responsive Scaling**: Images adapt to container size
- **Alt Text Support**: Built-in accessibility

#### 7. Code Highlighting
- Syntax highlighting for multiple languages
- Inline code and code blocks
- Language detection for better formatting

#### 8. Configuration Management
- Flexible YAML configuration
- Enable/disable features (Ollama, image extraction, auto-discovery)
- Override via command-line flags
- Environment-specific configs

### Advanced Features
- **Multi-Format Support**: Easily switch between PDF and Markdown
- **Custom Styling**: Modify CSS and HTML templates
- **API Integration**: Access structured data via HTTP endpoints
- **Extensible Architecture**: Clean, modular design (refactored from 396 → 79 lines in main.go!)
- **Performance Optimized**: Fast processing, quick startup, minimal dependencies
- **Security First**: Local processing, no external calls, privacy-focused

### Generated Output
```
docs/
├── index.html              # Lightweight shell (loads content dynamically)
├── data/
│   ├── content.json        # Complete documentation (all sections)
│   └── sections/           # Individual section files
│       ├── section-1.json  # Easy to edit, version control
│       └── ...
├── search-index.json       # Optimized for Fuse.js search
├── images/                 # Extracted/copied images
└── static/                 # CSS, JavaScript assets

---

## 📦 Prerequisites

### 1. Go (version 1.21 or higher)
```bash
# Check if Go is installed
go version

# If not installed, download from: https://go.dev/dl/
```

### 2. Ollama (for AI chat functionality)
```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh

# Windows: Download from https://ollama.com/download

# Verify installation
ollama --version
```

### 3. Install Ollama Model
```bash
# Download llama3.2 model (recommended)
ollama pull llama3.2

# Or try other models
ollama pull mistral
ollama pull codellama

# Verify models
ollama list
```

### 4. Optional: pdfimages (for better PDF image extraction)
```bash
# macOS
brew install poppler

# Linux
sudo apt-get install poppler-utils  # Debian/Ubuntu
sudo yum install poppler-utils       # RHEL/CentOS

# Windows: Download from https://blog.alivate.com.au/poppler-windows/
```

---

## 🚀 Quick Start

### 1. Setup Project

```bash
# Navigate to project directory
cd docTrainerGO

# Download Go dependencies
go mod download

# Build the application
go build ./cmd/main.go
```

### 2. Choose Your Input Type

#### Option A: Process Markdown Files

```bash
# Edit config.yaml and set:
# input_type: markdown

# Place markdown files in input/markdown/
cp your-docs/*.md input/markdown/

# Process and start server
./main -serve
```

#### Option B: Process PDF File

```bash
# Place PDF in input directory
cp ~/Documents/manual.pdf input/

# Process PDF directly
./main -pdf input/manual.pdf -serve
```

### 3. Start Ollama (in separate terminal)

```bash
# Start Ollama service
ollama run llama3.2

# Keep this terminal open for chat functionality
```

### 4. Access Documentation

Open browser and navigate to:
```
http://localhost:8080
```

---

## 📖 Usage

### Command-Line Interface

```bash
# Show help
./main -help

# Process with config file (default: config.yaml)
./main -serve

# Process PDF and start server
./main -pdf input/document.pdf -serve

# Process markdown and exit (no server)
./main -process

# Process PDF only (no server)
./main -pdf input/document.pdf -process

# Use custom config file
./main -config custom-config.yaml -serve
```

### Configuration File

Create or edit `config.yaml`:

```yaml
# Input type: 'pdf' or 'markdown'
input_type: markdown

# PDF configuration
pdf:
  path: "input/user_guide.pdf"
  extract_images: true

# Markdown configuration
markdown:
  directory: "input/markdown"
  auto_discover: true              # Automatically find all .md files
  files:                           # Or specify files manually
    - "01-introduction.md"
    - "02-getting-started.md"

# Output configuration
output:
  directory: "docs"
  title: "My Documentation"

# Server configuration
server:
  port: "8080"
  host: "localhost"

# Ollama configuration
ollama:
  enabled: true                    # Enable/disable AI chat functionality
  url: "http://localhost:11434"
  model: "llama3.2"
```

> **Note**: Set `ollama.enabled: false` to run without AI chat. This is useful for:
> - Deployments where Ollama is not available
> - Reducing resource usage (no LLM required)
> - Documentation-only sites without chat features

### Using Makefile

```bash
# Clean generated files
make clean

# Process markdown files
make process

# Process specific PDF
make process PDF=input/manual.pdf

# Start server
make serve

# Start server on custom port
make serve PORT=3000

# Build binary
make build

# Run all tests
make test
```

---

## 📁 Project Structure

```
docTrainerGO/
├── cmd/
│   └── main.go                    # 79 lines - Application entry point
│
├── internal/
│   ├── cli/
│   │   └── cli.go                 # 99 lines - Command-line interface
│   ├── config/
│   │   └── config.go              # 67 lines - Configuration management
│   ├── processor/
│   │   └── processor.go           # 216 lines - Document processing
│   ├── server/
│   │   └── server.go              # 176 lines - HTTP server & chat API
│   ├── pdf/
│   │   └── parser.go              # PDF parsing & image extraction
│   ├── md/
│   │   └── parser.go              # Markdown parsing & processing
│   ├── generator/
│   │   ├── data.go                # JSON data generation
│   │   └── html.go                # HTML generation
│   ├── chat/
│   │   └── ollama.go              # Ollama LLM integration
│   └── search/
│       └── index.go               # Search index generation
│
├── templates/
│   └── page.html                  # HTML template (lightweight - 3.7KB)
│
├── static/
│   ├── style.css                  # 639 lines - Responsive styling
│   ├── script.js                  # 415 lines - Dynamic content loading
│   └── fuse.min.js                # Fuse.js library (download separately)
│
├── input/
│   ├── markdown/                  # Place markdown files here
│   │   ├── 01-introduction.md
│   │   └── images/                # Markdown images
│   └── user_guide.pdf             # Or place PDF here
│
├── docs/                          # Generated documentation (auto-created)
│   ├── index.html                 # Main page (3.7KB)
│   ├── data/
│   │   ├── content.json           # Master content file
│   │   └── sections/              # Individual section JSON files
│   ├── images/                    # Extracted/copied images
│   └── search-index.json          # Search index
│
├── config.yaml                    # Configuration file
├── Makefile                       # Build automation
├── go.mod                         # Go module definition
├── go.sum                         # Go dependencies checksum
└── README.md                      # This file
```

---

## 🏗️ Architecture

### Modular Design

The project follows clean architecture principles with clear separation of concerns:

#### 1. **CLI Layer** (`internal/cli/`)
- Parses command-line arguments
- Provides user-friendly help text
- Validates input parameters

#### 2. **Configuration Layer** (`internal/config/`)
- Loads and validates YAML configuration
- Provides default values
- Single source of truth for settings

#### 3. **Processing Layer** (`internal/processor/`)
- Orchestrates document processing
- Handles both PDF and Markdown inputs
- Coordinates image extraction and data generation

#### 4. **Server Layer** (`internal/server/`)
- HTTP server management
- Chat API endpoints
- Static file serving
- CORS handling

#### 5. **Parser Layer** (`internal/pdf/`, `internal/md/`)
- PDF text and image extraction
- Markdown file parsing
- Heading detection and hierarchy

#### 6. **Generator Layer** (`internal/generator/`)
- JSON data structure generation
- HTML page generation
- Search index creation

#### 7. **Chat Layer** (`internal/chat/`)
- Ollama LLM integration
- Context management
- Request/response handling

### Data Flow

```
Input (PDF/Markdown)
        ↓
    Processor
        ↓
   Parser (PDF/MD)
        ↓
    Generator
        ↓
   Output (docs/)
        ├── HTML (lightweight)
        ├── JSON (content + sections)
        └── Search Index
        ↓
     Server
        ├── Static Files
        ├── Documentation
        └── Chat API → Ollama
```

### Key Improvements (Recent Refactoring)

**Before:**
- 396 lines in single main.go file
- Mixed concerns and responsibilities
- Hard to test and maintain
- Global state and tight coupling

**After:**
- 79 lines in main.go (80% reduction!)
- 4 new modular packages (cli, config, processor, server)
- Single responsibility per package
- Easy to test, extend, and maintain
- No global state, clean interfaces

**Metrics:**

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| main.go lines | 396 | 79 | **-80%** |
| Functions in main.go | 10+ | 1 | **-90%** |
| Largest file | 396 lines | 216 lines | **-45%** |
| Packages | 7 | 11 | Better organization |
| Testability | Low | High | ✅ |
| Maintainability | Poor | Excellent | ✅ |

---

## 🔧 Development

### Building

```bash
# Build for current platform
go build ./cmd/main.go

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o main-linux ./cmd/main.go
GOOS=windows GOARCH=amd64 go build -o main.exe ./cmd/main.go
GOOS=darwin GOARCH=arm64 go build -o main-mac ./cmd/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Test specific package
go test ./internal/processor/
```

### Adding New Features

#### Add a New Input Type

1. Create parser in `internal/yourtype/parser.go`
2. Implement `Parse()` method returning `*pdf.Document`
3. Add case in `processor.Process()` switch statement
4. Update `config.yaml` schema

#### Add New Output Format

1. Create generator in `internal/generator/yourformat.go`
2. Implement `Generate(*pdf.Document)` method
3. Call from `processor.Process()`

#### Extend Chat Functionality

1. Add methods to `internal/chat/ollama.go`
2. Add endpoints in `internal/server/server.go`
3. Update frontend in `static/script.js`

---

## 🔍 Troubleshooting

### Issue: Build errors

```bash
# Solution: Ensure dependencies are up to date
go mod tidy
go mod download
go build ./cmd/main.go
```

### Issue: "Failed to load config"

```bash
# Solution: Check config.yaml exists and is valid YAML
./main -help  # See default config path

# Validate YAML syntax
cat config.yaml | grep -v "^#"
```

### Issue: "Ollama connection failed"

```bash
# Solution 1: Start Ollama service
ollama serve  # In one terminal
ollama run llama3.2  # In another terminal

# Solution 2: Check Ollama URL in config
curl http://localhost:11434/api/version

# Solution 3: Update config.yaml
ollama:
  url: "http://localhost:11434"
  model: "llama3.2"
```

### Issue: Images not displaying

```bash
# For PDFs: Install pdfimages
brew install poppler  # macOS
sudo apt-get install poppler-utils  # Linux

# For Markdown: Check image paths
ls input/markdown/images/

# Reprocess after fixing
make clean
./main -process
```

### Issue: Search not working

```bash
# Solution: Verify search index exists
ls docs/search-index.json

# Check search index content
cat docs/search-index.json | jq '.items | length'

# Reprocess if needed
make clean
./main -process
```

### Issue: Port already in use

```bash
# Solution 1: Use different port (edit config.yaml)
server:
  port: "3000"

# Solution 2: Kill process on port
lsof -ti:8080 | xargs kill -9

# Solution 3: Find what's using the port
lsof -i :8080
```

### Issue: Chat responses are generic

**Problem:** Chat doesn't seem to know about your documentation

**Solution:** Check if documentation context is loading:
1. Verify `docs/data/content.json` exists and has content
2. Check file size: `ls -lh docs/data/content.json`
3. Restart server to reload context
4. Try asking more specific questions about your docs

### Issue: PDF parsing fails

```bash
# Check PDF is valid
file input/your-file.pdf

# Try with a simple PDF first
# Some PDFs with complex formatting may not parse correctly

# Check PDF is not encrypted
pdfinfo input/your-file.pdf | grep Encrypted
```

---

## 🌐 Deployment

### Local Network Deployment

```bash
# Find your IP address
ifconfig | grep "inet " | grep -v 127.0.0.1

# Start server (accessible from network)
./main -serve

# Access from other devices: http://YOUR_IP:8080
```

### Production Deployment

```bash
# 1. Build optimized binary
go build -ldflags="-s -w" -o doctrainer ./cmd/main.go

# 2. Create systemd service (Linux)
sudo nano /etc/systemd/system/doctrainer.service

# Add:
# [Unit]
# Description=DocTrainer Service
# After=network.target
#
# [Service]
# Type=simple
# User=youruser
# WorkingDirectory=/path/to/docTrainerGO
# ExecStart=/path/to/docTrainerGO/doctrainer -serve
# Restart=always
#
# [Install]
# WantedBy=multi-user.target

# 3. Enable and start
sudo systemctl enable doctrainer
sudo systemctl start doctrainer
```

### Docker Deployment

#### Option 1: Documentation Only (No AI Chat)

Perfect for lightweight deployments where AI chat is not needed:

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates poppler-utils
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./docs
EXPOSE 8080
CMD ["./main", "-serve"]
```

Set `ollama.enabled: false` in config.yaml before building.

```bash
# Build and run
docker build -t doctrainer .
docker run -p 8080:8080 -v $(pwd)/docs:/root/docs doctrainer
```

#### Option 2: With AI Chat (Docker Compose)

For deployments requiring AI chat functionality, use separate containers:

```yaml
# docker-compose.yml
version: '3.8'

services:
  ollama:
    image: ollama/ollama:latest
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    environment:
      - OLLAMA_HOST=0.0.0.0
    command: serve

  doctrainer:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./docs:/root/docs
      - ./config.yaml:/root/config.yaml
    depends_on:
      - ollama
    environment:
      - OLLAMA_URL=http://ollama:11434

volumes:
  ollama_data:
```

Update config.yaml to use Docker service name:
```yaml
ollama:
  enabled: true
  url: "http://ollama:11434"  # Use service name
  model: "llama3.2"
```

```bash
# Start both services
docker-compose up -d

# Pull Ollama model (first time only)
docker exec -it <ollama-container-id> ollama pull llama3.2

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

**Resource Requirements:**
- Documentation only: ~50MB RAM, 100MB disk
- With Ollama: 4-6GB RAM, ~2GB disk (for llama3.2 model)

---

## 💡 Tips & Best Practices

### Performance Optimization
- Process large PDFs in chunks if memory is limited
- Use `auto_discover: true` for Markdown to automatically find files
- Search index is loaded entirely in browser - keep it reasonable (<5MB)
- Images are lazy-loaded for better initial page performance

### Content Organization
- Use clear heading hierarchy (H1 > H2 > H3)
- Keep section sizes reasonable (500-2000 words)
- Use descriptive heading text for better search
- Include images to break up text

### Chat Effectiveness
- Ask specific questions about your documentation
- Reference section headings in questions
- Provide context in multi-turn conversations
- Use code/technical terms from your docs

### Security Considerations
- This tool is designed for **local/internal use**
- For production: add authentication layer
- Rate-limit chat API if exposed publicly
- Validate and sanitize all file inputs
- Keep Ollama and dependencies updated

---

## 🤝 Contributing

Contributions welcome! Areas for improvement:
- Add unit tests for all packages
- Implement table of contents auto-generation
- Add support for DOCX/ODT input formats
- Create VS Code extension
- Add export to multiple formats
- Implement version comparison
- Add collaborative editing features

---

## 📄 License

This project is provided as-is for educational and personal use.

---

## 🙏 Acknowledgments

- **[ledongthuc/pdf](https://github.com/ledongthuc/pdf)** - PDF parsing library
- **[Fuse.js](https://www.fusejs.io/)** - Lightweight fuzzy-search library
- **[Ollama](https://ollama.com/)** - Run LLMs locally
- **[poppler](https://poppler.freedesktop.org/)** - PDF rendering library

---

## 📞 Support

**Need help?**
1. Check [Troubleshooting](#troubleshooting) section above
2. Verify [Prerequisites](#prerequisites) are correctly installed
3. Review browser console (F12) for JavaScript errors
4. Check server logs for error messages
5. Try with sample markdown files first

**Quick Diagnostics:**
```bash
# Check Go version
go version

# Check Ollama status
ollama list

# Verify project structure
ls -la cmd/ internal/ static/ templates/

# Test build
go build ./cmd/main.go
./main -help
```

---

**Built with ❤️ using Go, Ollama, and modern web technologies**

*Last Updated: January 2026*
