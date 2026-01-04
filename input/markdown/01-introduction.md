# Introduction

Welcome to **DocTrainerGO** - a powerful tool that converts your documentation into beautiful, searchable websites with AI-powered chat assistance.

![DocTrainerGO Logo](images/logo.png)

## What is DocTrainerGO?

DocTrainerGO is an innovative documentation platform that:

- Converts PDF or Markdown files into interactive websites
- Provides intelligent search functionality
- Includes AI-powered chat assistance using local LLMs
- Generates organized, maintainable data structures
- Offers beautiful, responsive UI

## Key Features

### 1. Multiple Input Formats
Support for both PDF and Markdown files, allowing you to work with your existing documentation.

### 2. AI-Powered Chat
Built-in chat assistant that answers questions based on your documentation content using Ollama.

### 3. Smart Search
Fast, fuzzy search powered by Fuse.js to help users find information quickly.

![Search Feature](images/search-demo.png)

### 4. Organized Data Structure
All content is stored in structured JSON files for easy maintenance and updates.

## Architecture

The system follows a clean, modular architecture:

```
docTrainerGO/
├── cmd/           # Main application
├── internal/      # Internal packages
│   ├── pdf/       # PDF processing
│   ├── md/        # Markdown processing
│   ├── generator/ # HTML/Data generation
│   ├── search/    # Search indexing
│   └── chat/      # AI chat integration
├── templates/     # HTML templates
├── static/        # CSS, JavaScript
└── docs/          # Generated output
```

## Getting Started

To get started with DocTrainerGO, you'll need:

1. **Go 1.21+** - Programming language runtime
2. **Ollama** - Local LLM runtime
3. **Poppler** - For PDF image extraction (optional)

![Getting Started](images/installation.png)

## Use Cases

DocTrainerGO is perfect for:

- **Product Documentation** - Create beautiful docs for your products
- **Internal Knowledge Base** - Organize company knowledge with AI search
- **User Guides** - Convert manuals into interactive websites
- **API Documentation** - Present API docs with intelligent assistance
- **Training Materials** - Build interactive training resources

## Why DocTrainerGO?

Unlike other documentation tools, DocTrainerGO:

- ✅ **Runs Locally** - No cloud dependencies, complete privacy
- ✅ **AI-Enhanced** - Built-in chat using your own LLM
- ✅ **Flexible Input** - Works with PDF or Markdown
- ✅ **Lightweight** - Minimal dependencies, fast performance
- ✅ **Easy Maintenance** - Organized data structure
- ✅ **Beautiful UI** - Modern, responsive design

## Next Steps

Continue to the [Getting Started](02-getting-started.md) guide to install and configure DocTrainerGO.
