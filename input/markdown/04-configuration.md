# Configuration

Complete guide to configuring DocTrainerGO for your needs.

## Configuration File

DocTrainerGO uses `config.yaml` for all settings.

### Basic Structure

```yaml
# Input source type
input_type: markdown  # or "pdf"

# PDF configuration
pdf:
  path: input/document.pdf
  extract_images: true

# Markdown configuration
markdown:
  directory: input/markdown
  auto_discover: true
  files: []

# Output configuration
output:
  directory: docs
  title: "Documentation"

# Server configuration
server:
  port: 8080
  host: localhost

# Ollama configuration
ollama:
  url: http://localhost:11434
  model: llama3.2
```

## Input Configuration

### PDF Mode

```yaml
input_type: pdf

pdf:
  path: input/user_guide.pdf
  extract_images: true
```

**Options:**
- `path` - Path to PDF file
- `extract_images` - Use pdfimages for extraction (requires poppler)

### Markdown Mode

```yaml
input_type: markdown

markdown:
  directory: input/markdown
  auto_discover: true
  files:
    - input/markdown/01-intro.md
    - input/markdown/02-guide.md
```

**Options:**
- `directory` - Directory containing .md files
- `auto_discover` - Automatically find all .md files
- `files` - Explicitly list files to process (when auto_discover is false)

## Output Configuration

### Output Directory

```yaml
output:
  directory: docs
  title: "My Documentation"
```

**Generated Structure:**
```
docs/
├── data/
│   ├── content.json
│   └── sections/
├── images/
├── index.html
└── search-index.json
```

### Title Configuration

The `title` appears in:
- Browser tab
- Page header
- Navigation sidebar
- Generated metadata

## Server Configuration

```yaml
server:
  port: 8080
  host: localhost
```

**Options:**
- `port` - Server port (default: 8080)
- `host` - Server host (default: localhost)

### Custom Port Example

```bash
# Via config
server:
  port: 3000

# Via command line (overrides config)
go run cmd/main.go -serve -port=3000
```

## Ollama Configuration

```yaml
ollama:
  url: http://localhost:11434
  model: llama3.2
```

**Options:**
- `url` - Ollama API endpoint
- `model` - LLM model name

### Supported Models

DocTrainerGO works with any Ollama model:

```bash
# List available models
ollama list

# Pull a model
ollama pull llama3.2
ollama pull mistral
ollama pull codellama
```

**Recommended Models:**
- `llama3.2` - Fast, good quality
- `llama3.1` - More powerful
- `mistral` - Alternative option
- `codellama` - For technical docs

## Environment Variables

Override config with environment variables:

```bash
# Port
export DOCTRAINER_PORT=3000

# Ollama URL
export OLLAMA_URL=http://localhost:11434

# Ollama Model
export OLLAMA_MODEL=llama3.2
```

## Advanced Configuration

### Custom Image Directory

```yaml
markdown:
  directory: input/markdown
  image_directory: input/markdown/images
```

### Processing Options

```yaml
processing:
  # Maximum section size (characters)
  max_section_size: 10000
  
  # Minimum heading level to create sections
  min_heading_level: 1
  
  # Maximum heading level to create sections
  max_heading_level: 6
```

### Search Configuration

```yaml
search:
  # Fuse.js threshold (0.0 = exact, 1.0 = match anything)
  threshold: 0.4
  
  # Minimum characters for search
  min_length: 2
  
  # Maximum results to show
  max_results: 10
```

## Configuration Examples

### Example 1: PDF Documentation

```yaml
input_type: pdf
pdf:
  path: input/manual.pdf
  extract_images: true
output:
  title: "User Manual"
ollama:
  model: llama3.2
```

### Example 2: Markdown Knowledge Base

```yaml
input_type: markdown
markdown:
  directory: input/kb
  auto_discover: true
output:
  title: "Knowledge Base"
  directory: build
server:
  port: 8000
```

### Example 3: API Documentation

```yaml
input_type: markdown
markdown:
  directory: input/api-docs
  files:
    - input/api-docs/overview.md
    - input/api-docs/authentication.md
    - input/api-docs/endpoints.md
    - input/api-docs/examples.md
output:
  title: "API Documentation"
ollama:
  model: codellama
```

## Configuration Validation

DocTrainerGO validates configuration on startup:

```bash
go run cmd/main.go -config=config.yaml -validate
```

**Checks:**
- Required fields present
- Valid input type
- Input files/directories exist
- Output directory writable
- Ollama accessible

## Configuration Best Practices

### 1. Use Version Control

```bash
# Track config.yaml
git add config.yaml
git commit -m "Add documentation config"
```

### 2. Environment-Specific Configs

```
config.yaml           # Development
config.prod.yaml      # Production
config.staging.yaml   # Staging
```

### 3. Sensitive Data

Don't commit sensitive information:

```yaml
# DON'T do this
ollama:
  api_key: "secret-key-123"  # Bad!

# Instead use environment variables
ollama:
  api_key: ${OLLAMA_API_KEY}  # Good!
```

### 4. Documentation

Comment your configuration:

```yaml
# AI Model Configuration
# Using llama3.2 for faster responses
# Switch to llama3.1 for better quality
ollama:
  model: llama3.2
```

## Troubleshooting

### Config Not Found

```bash
# Specify config path
go run cmd/main.go -config=./configs/prod.yaml
```

### Invalid YAML Syntax

```bash
# Validate YAML
yamllint config.yaml

# Or use online validator
# https://www.yamllint.com/
```

### Config Changes Not Applied

```bash
# Restart server after config changes
# Kill existing process
pkill -f "docTrainerGO"

# Start with new config
go run cmd/main.go -serve
```

## Next Steps

Explore [Advanced Usage](05-advanced.md) for power user features.
