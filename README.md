# DocTrainerGO

🚀 **Convert PDF documents into fully searchable, interactive documentation websites with AI-powered chat assistance.**

DocTrainerGO is a complete Go-based solution that transforms PDFs (with text and images) into beautiful, responsive documentation websites featuring:
- ✅ Sidebar navigation with collapsible sections
- ✅ Real-time search with Fuse.js
- ✅ Embedded images from PDFs
- ✅ Floating AI chat assistant powered by local Ollama LLM
- ✅ Fully responsive design (mobile & desktop)
- ✅ Runs entirely on your laptop (no cloud dependencies)

---

## 📋 Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Deployment](#deployment)
- [License](#license)

---

## ✨ Features

### PDF Processing
- Extract text, headings, and images from PDF files
- Intelligent section detection and hierarchy
- Automatic image extraction and embedding

### Documentation Website
- Clean, modern UI with sidebar navigation
- Client-side search using Fuse.js
- Responsive design for all devices
- Image galleries with lazy loading
- Smooth scrolling and navigation

### AI Chat Assistant
- Floating chat widget
- Powered by local Ollama LLM (llama3.1, mistral, etc.)
- Context-aware responses
- No data leaves your machine

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

### 1. Go (version 1.21 or higher)
```bash
# Check if Go is installed
go version

# If not installed, download from: https://go.dev/dl/
```

### 2. Ollama (for AI chat functionality)
```bash
# macOS/Linux
curl -fsSL https://ollama.com/install.sh | sh

# Windows: Download from https://ollama.com/download

# Verify installation
ollama --version
```

### 3. Install Ollama Model
```bash
# Download and install llama3.1 model (or any other model)
ollama pull llama3.1

# Verify the model is available
ollama list
```

---

## 🚀 Installation

### Step 1: Clone or Download the Project

If you received this as files, ensure all files are in the correct structure. Otherwise:

```bash
cd ~/Desktop/workspace/LLM/docTrainerGO
```

### Step 2: Install Go Dependencies

```bash
# Download required Go modules
go mod download

# Verify dependencies
go mod verify
```

### Step 3: Download Fuse.js

Download the Fuse.js library for search functionality:

```bash
# Download Fuse.js v6.6.2
curl -L https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js

# Or download manually from: https://www.fusejs.io/
# and place it in the static/ directory
```

---

## 🏃 Quick Start

### 1. Prepare Your PDF

Place your PDF file in the `input/` directory:

```bash
mkdir -p input
cp /path/to/your/document.pdf input/
```

### 2. Process the PDF

```bash
# Process the PDF and generate documentation
go run cmd/main.go -pdf=input/document.pdf
```

This will:
- Extract text and images from the PDF
- Generate HTML pages in the `docs/` directory
- Create a search index JSON file
- Save images to `docs/images/`

### 3. Start Ollama

In a separate terminal, start the Ollama service with your chosen model:

```bash
# Start Ollama with llama3.1
ollama run llama3.1

# Keep this terminal open while using the chat feature
```

### 4. Start the Web Server

```bash
# Serve the documentation site
go run cmd/main.go -serve
```

### 5. Access the Documentation

Open your browser and navigate to:

```
http://localhost:8080
```

You should see your PDF converted to a beautiful documentation website with a working AI chat assistant!

---

## 📖 Usage

### Command-Line Options

```bash
# Process a PDF file
go run cmd/main.go -pdf=input/your-document.pdf

# Serve the documentation (default port: 8080)
go run cmd/main.go -serve

# Serve on a custom port
go run cmd/main.go -serve -port=3000

# Use a different Ollama URL
go run cmd/main.go -serve -ollama=http://localhost:11434

# Use a different Ollama model
go run cmd/main.go -serve -model=mistral

# Show help
go run cmd/main.go -h
```

### Complete Workflow

```bash
# 1. Process your PDF
go run cmd/main.go -pdf=input/user-manual.pdf

# 2. Start Ollama (in another terminal)
ollama run llama3.1

# 3. Start the web server
go run cmd/main.go -serve -port=8080

# 4. Open browser to http://localhost:8080
```

### Using Different Ollama Models

```bash
# Install different models
ollama pull mistral
ollama pull codellama
ollama pull phi

# Use a specific model
go run cmd/main.go -serve -model=mistral
```

---

## 📁 Project Structure

```
docTrainerGO/
├── cmd/
│   └── main.go                 # Main application entry point
├── internal/
│   ├── pdf/
│   │   └── parser.go           # PDF parsing and image extraction
│   ├── generator/
│   │   └── html.go             # HTML page generation
│   ├── search/
│   │   └── index.go            # Search index generation
│   └── chat/
│       └── ollama.go           # Ollama LLM integration
├── templates/
│   └── page.html               # HTML template for documentation pages
├── static/
│   ├── style.css               # CSS styling
│   ├── script.js               # Frontend JavaScript
│   └── fuse.min.js             # Fuse.js library (download separately)
├── input/
│   └── (place your PDFs here)
├── docs/                       # Generated documentation (created automatically)
│   ├── index.html
│   ├── search-index.json
│   └── images/                 # Extracted images
├── go.mod                      # Go module file
├── go.sum                      # Go dependencies checksum
└── README.md                   # This file
```

---

## ⚙️ Configuration

### Ollama Configuration

By default, the application connects to Ollama at `http://localhost:11434`. To customize:

```bash
# Use a remote Ollama instance
go run cmd/main.go -serve -ollama=http://192.168.1.100:11434

# Use a different model
go run cmd/main.go -serve -model=codellama
```

### Customizing the UI

Edit the following files to customize appearance:

- **Colors & Styling**: `static/style.css`
- **Layout & Structure**: `templates/page.html`
- **JavaScript Behavior**: `static/script.js`

### PDF Parser Settings

To adjust PDF parsing behavior, edit `internal/pdf/parser.go`:

- Modify heading detection patterns
- Adjust section hierarchy logic
- Customize image extraction

---

## 🔧 Troubleshooting

### Issue: "Failed to open PDF"

**Solution:**
- Ensure the PDF file exists at the specified path
- Check file permissions
- Verify the PDF is not corrupted

```bash
# Check if file exists
ls -la input/your-file.pdf

# Check file type
file input/your-file.pdf
```

### Issue: "Ollama is not accessible"

**Solution:**
- Make sure Ollama is running:

```bash
# Check if Ollama is running
ps aux | grep ollama

# Start Ollama
ollama serve

# In another terminal, run the model
ollama run llama3.1
```

### Issue: "Search not working"

**Solution:**
- Ensure `fuse.min.js` is in the `static/` directory
- Check browser console for JavaScript errors
- Verify `search-index.json` was generated in `docs/`

```bash
# Download Fuse.js if missing
curl -L https://cdn.jsdelivr.net/npm/fuse.js@6.6.2/dist/fuse.min.js -o static/fuse.min.js
```

### Issue: "Images not displaying"

**Solution:**
- Check if images were extracted to `docs/images/`
- Verify image file permissions
- Check browser console for 404 errors

```bash
# Check extracted images
ls -la docs/images/
```

### Issue: Port already in use

**Solution:**
```bash
# Use a different port
go run cmd/main.go -serve -port=3000

# Or find and kill the process using the port
lsof -ti:8080 | xargs kill
```

---

## 🌐 Deployment

### Option 1: Local Network Access

Make the documentation accessible to other devices on your network:

```bash
# Start server and allow external connections
go run cmd/main.go -serve -port=8080

# Find your local IP
ifconfig | grep "inet "

# Access from other devices: http://YOUR_IP:8080
```

### Option 2: Build Binary

Create a standalone executable:

```bash
# Build for your platform
go build -o doctrainer cmd/main.go

# Run the binary
./doctrainer -serve

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o doctrainer-linux cmd/main.go
GOOS=windows GOARCH=amd64 go build -o doctrainer.exe cmd/main.go
GOOS=darwin GOARCH=arm64 go build -o doctrainer-mac cmd/main.go
```

### Option 3: Static Export

Export as static files for hosting on any web server:

```bash
# After processing PDF and generating docs/
# Copy these directories to your web server:
# - docs/
# - static/

# Note: Chat functionality requires the Go backend to be running
```

### Option 4: Docker (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/main.go
EXPOSE 8080
CMD ["./main", "-serve"]
```

Build and run:

```bash
docker build -t doctrainer .
docker run -p 8080:8080 doctrainer
```

---

## 💡 Tips & Best Practices

### Performance

- PDFs with many images may take longer to process
- Large PDFs (>100 pages) may require more memory
- Search index is loaded entirely in browser memory

### Security

- This application is designed for local use
- For production, add authentication
- Validate all PDF inputs
- Rate-limit chat API requests

### Customization Ideas

1. **Add syntax highlighting** for code blocks
2. **Implement markdown support** in content
3. **Add PDF bookmarks** to navigation
4. **Create table of contents** generator
5. **Add export to PDF** functionality
6. **Implement dark mode**
7. **Add full-text search** with highlights

---

## 🤝 Contributing

Feel free to:
- Report issues
- Suggest features
- Submit pull requests
- Improve documentation

---

## 📄 License

This project is provided as-is for educational and personal use.

---

## 🙏 Acknowledgments

- **[ledongthuc/pdf](https://github.com/ledongthuc/pdf)** - PDF parsing library
- **[Fuse.js](https://www.fusejs.io/)** - Fuzzy search library
- **[Ollama](https://ollama.com/)** - Local LLM runtime

---

## 📞 Support

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section
2. Verify all [Prerequisites](#prerequisites) are installed
3. Check the browser console for errors
4. Review server logs for error messages

---

**Happy Documentation Building! 📚✨**
