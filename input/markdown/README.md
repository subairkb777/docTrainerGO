# Sample Markdown Documentation

This directory contains sample Markdown files for testing DocTrainerGO's Markdown processing feature.

## Files

1. **01-introduction.md** - Introduction and overview
2. **02-getting-started.md** - Installation and setup guide
3. **03-features.md** - Feature documentation
4. **04-configuration.md** - Configuration guide
5. **05-advanced.md** - Advanced usage and tips

## Image References

The markdown files reference images in the `images/` directory. To test with actual images:

1. Add PNG images to `input/markdown/images/`
2. Use these filenames (referenced in the docs):
   - logo.png
   - search-demo.png
   - installation.png
   - install-steps.png
   - config-example.png
   - pdf-process.png
   - md-process.png
   - chat-demo.png
   - search-results.png
   - responsive.png

## Processing

To process these markdown files:

```bash
# Using config
go run cmd/main.go -config=config.yaml

# Direct command (once implemented)
go run cmd/main.go -md=input/markdown/*.md
```

## Structure

Each file follows a clear structure:
- Headings (# ## ###)
- Paragraphs
- Code blocks
- Lists (ordered and unordered)
- Images with ![alt](path) syntax
- Links
- Tables
- Emphasis (*italic*, **bold**, `code`)

## Testing

These files are designed to test:
- Multiple document processing
- Image path resolution
- Code syntax highlighting
- Table rendering
- List formatting
- Link handling
- Section hierarchy
