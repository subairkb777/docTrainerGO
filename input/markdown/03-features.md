# Features

DocTrainerGO offers a comprehensive set of features to create powerful documentation websites.

## Core Features

### 1. Dual Input Support

Process documentation from multiple sources:

#### PDF Processing
- Extract text and structure from PDF files
- Automatic image extraction using Poppler
- Section detection and organization
- Preserve formatting and hierarchy

![PDF Processing](images/pdf-process.png)

#### Markdown Processing
- Parse standard Markdown syntax
- Support for images, code blocks, tables
- Front matter metadata support
- Link resolution and validation

![Markdown Processing](images/md-process.png)

### 2. AI-Powered Chat Assistant

Interactive chat that understands your documentation:

- **Context-Aware Responses** - Answers based on your content
- **Local LLM** - Runs on your machine using Ollama
- **Formatted Responses** - Supports lists, code blocks, emphasis
- **Real-time Interaction** - Fast responses with streaming support

```javascript
// Chat uses your documentation as context
const context = loadDocumentationContext();
const answer = await ollamaClient.AskWithContext(prompt, context);
```

![Chat Interface](images/chat-demo.png)

#### Supported Chat Features

- ✅ Numbered and bulleted lists
- ✅ **Bold** and *italic* text
- ✅ `Inline code` and code blocks
- ✅ Contextual answers from your docs
- ✅ Error handling and graceful degradation

### 3. Advanced Search

Fast, intelligent search powered by Fuse.js:

- **Fuzzy Matching** - Find content even with typos
- **Weighted Results** - Prioritize titles over content
- **Real-time Results** - Instant search as you type
- **Keyboard Shortcuts** - `Ctrl/Cmd + K` to focus search

```javascript
// Search configuration
const options = {
  keys: ['heading', 'content'],
  threshold: 0.4,
  minMatchCharLength: 2
};
```

![Search Results](images/search-results.png)

### 4. Organized Data Structure

Content stored in maintainable JSON format:

```
docs/data/
├── content.json           # Master file (all content)
└── sections/              # Individual sections
    ├── section-1.json     # Introduction
    ├── section-2.json     # Getting Started
    └── ...
```

#### Benefits

- 📝 **Easy Editing** - Modify individual sections
- 🔄 **Version Control** - Track changes with git
- 🔌 **API-Ready** - Use data in other applications
- 🎯 **Single Source of Truth** - One data source for all outputs

### 5. Responsive UI

Beautiful interface that works everywhere:

- **Mobile-First Design** - Optimized for all screen sizes
- **Sidebar Navigation** - Collapsible on mobile
- **Floating Chat Widget** - Minimizable and accessible
- **Lazy Image Loading** - Fast page loads
- **Dark Mode Ready** - Easy to customize

![Responsive Design](images/responsive.png)

#### UI Components

| Component | Features |
|-----------|----------|
| Sidebar | Collapsible, searchable navigation |
| Content Area | Clean typography, code highlighting |
| Chat Widget | Floating, minimizable, keyboard accessible |
| Search | Dropdown results with highlighting |

### 6. Image Support

Comprehensive image handling:

- **Automatic Extraction** - From PDFs using pdfimages
- **Lazy Loading** - Images load as needed
- **Responsive Images** - Scale to container
- **Alt Text Support** - Accessibility built-in

### 7. Code Highlighting

Syntax highlighting for code blocks:

```go
// Example Go code
func main() {
    fmt.Println("Hello, DocTrainerGO!")
}
```

```python
# Example Python code
def process_markdown(file):
    with open(file, 'r') as f:
        return parse(f.read())
```

### 8. Configuration Management

Flexible configuration via YAML:

```yaml
input_type: markdown
markdown:
  directory: input/markdown
  auto_discover: true
output:
  title: "My Docs"
ollama:
  model: llama3.2
```

## Advanced Features

### Multi-Format Support

Switch between PDF and Markdown easily:

```bash
# Process PDF
go run cmd/main.go -pdf=input/guide.pdf

# Process Markdown
go run cmd/main.go -config=config.yaml
```

### Custom Styling

Easy to customize appearance:

- Modify `static/style.css` for styling
- Edit `templates/page.html` for structure
- CSS variables for theming

### API Integration

Use generated data in other applications:

```bash
# Access structured data
curl http://localhost:8080/docs/data/content.json

# Individual sections
curl http://localhost:8080/docs/data/sections/section-1.json
```

### Extensible Architecture

Clean, modular design for easy extension:

- Add new input formats
- Custom generators
- Alternative frontends
- Integration with other tools

## Performance

Optimized for speed:

- ⚡ **Fast Processing** - Efficient parsing algorithms
- 🚀 **Quick Startup** - Lightweight server
- 📦 **Small Footprint** - Minimal dependencies
- 🔄 **Incremental Updates** - Process only changed files

## Security

Built with security in mind:

- 🔒 **Local Processing** - No cloud dependencies
- 🛡️ **No External Calls** - Everything runs locally
- 🔐 **Privacy First** - Your data stays on your machine
- ✅ **Safe Defaults** - Secure configuration out of the box

## Next Steps

Learn how to configure DocTrainerGO in the [Configuration Guide](04-configuration.md).
