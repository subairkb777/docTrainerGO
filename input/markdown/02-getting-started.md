# Getting Started

This guide will help you install and configure DocTrainerGO on your system.

## Prerequisites

Before installing DocTrainerGO, ensure you have the following:

### Required

- **Go 1.21 or higher**
  ```bash
  go version
  # Should output: go version go1.21.0 or higher
  ```

- **Ollama** - For AI chat functionality
  ```bash
  # Install Ollama
  curl -fsSL https://ollama.com/install.sh | sh
  
  # Pull a model
  ollama pull llama3.2
  ```

### Optional

- **Poppler** - For better PDF image extraction
  ```bash
  # macOS
  brew install poppler
  
  # Ubuntu/Debian
  sudo apt install poppler-utils
  ```

![Installation Process](images/install-steps.png)

## Installation

### Step 1: Clone the Repository

```bash
git clone https://github.com/subairkb777/docTrainerGO-pdf.git
cd docTrainerGO-pdf
```

### Step 2: Install Dependencies

```bash
make setup
```

This command will:
- Download Go dependencies
- Download Fuse.js library
- Verify installation

### Step 3: Verify Installation

```bash
make test
```

## Configuration

DocTrainerGO uses a `config.yaml` file for configuration.

### Basic Configuration

```yaml
# Input type: pdf or markdown
input_type: markdown

# Markdown settings
markdown:
  directory: input/markdown
  auto_discover: true

# Output settings
output:
  directory: docs
  title: "My Documentation"

# Ollama settings
ollama:
  url: http://localhost:11434
  model: llama3.2
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `input_type` | Source format (pdf/markdown) | `pdf` |
| `markdown.directory` | Markdown files location | `input/markdown` |
| `output.directory` | Generated files location | `docs` |
| `ollama.model` | LLM model name | `llama3.2` |
| `server.port` | Web server port | `8080` |

![Configuration Example](images/config-example.png)

## Quick Start

### Process Markdown Files

1. **Add your Markdown files** to `input/markdown/`
2. **Configure** `config.yaml` to use markdown
3. **Process** the files:
   ```bash
   go run cmd/main.go -config=config.yaml
   ```

### Process PDF Files

1. **Add your PDF** to `input/`
2. **Process** the PDF:
   ```bash
   go run cmd/main.go -pdf=input/your-file.pdf
   ```

### Start the Server

```bash
make serve
# or
go run cmd/main.go -serve
```

Visit `http://localhost:8080` to view your documentation.

## Directory Structure

After installation, your project structure will look like:

```
docTrainerGO/
├── config.yaml          # Configuration file
├── input/               # Input files
│   ├── markdown/        # Markdown files
│   └── *.pdf            # PDF files
├── docs/                # Generated output
│   ├── data/
│   │   ├── content.json
│   │   └── sections/
│   ├── images/
│   └── index.html
├── static/              # Frontend assets
├── templates/           # HTML templates
└── cmd/                 # Main application
```

## Troubleshooting

### Ollama Not Running

**Error:** `Failed to connect to Ollama`

**Solution:**
```bash
# Start Ollama
ollama serve

# In another terminal
ollama run llama3.2
```

### Port Already in Use

**Error:** `Port 8080 already in use`

**Solution:**
```bash
# Use a different port
go run cmd/main.go -serve -port=3000
```

### Missing Dependencies

**Error:** `Command not found: pdfimages`

**Solution:**
```bash
# macOS
brew install poppler

# Ubuntu
sudo apt install poppler-utils
```

## Next Steps

Now that you have DocTrainerGO installed, explore the [Features](03-features.md) to learn what you can do.
