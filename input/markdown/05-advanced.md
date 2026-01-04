# Advanced Usage

Advanced techniques and tips for power users.

## Command Line Interface

### Processing Commands

```bash
# Process PDF
go run cmd/main.go -pdf=input/doc.pdf

# Process with config
go run cmd/main.go -config=config.yaml

# Process specific markdown files
go run cmd/main.go -md=input/markdown/*.md

# Process and serve immediately
go run cmd/main.go -pdf=input/doc.pdf -serve
```

### Server Commands

```bash
# Start server
go run cmd/main.go -serve

# Custom port
go run cmd/main.go -serve -port=3000

# Custom host and port
go run cmd/main.go -serve -port=8080
```

### Utility Commands

```bash
# Validate configuration
go run cmd/main.go -config=config.yaml -validate

# Check Ollama connection
go run cmd/main.go -check-ollama

# Version information
go run cmd/main.go -version

# Help
go run cmd/main.go -help
```

## Makefile Automation

### Custom Make Targets

Add to your `Makefile`:

```makefile
# Process markdown
process-md:
    go run cmd/main.go -config=config.yaml

# Build for production
build-prod:
    CGO_ENABLED=0 go build -ldflags="-s -w" -o doctrainer cmd/main.go

# Run with hot reload
dev:
    air -c .air.toml
```

### Batch Processing

```bash
# Process multiple PDFs
for pdf in input/*.pdf; do
    go run cmd/main.go -pdf="$pdf"
done

# Process multiple configs
for config in configs/*.yaml; do
    go run cmd/main.go -config="$config"
done
```

## Custom Markdown Extensions

### Front Matter Support

Add metadata to your markdown files:

```markdown
---
title: "Getting Started"
author: "John Doe"
date: 2026-01-04
tags: [tutorial, beginner]
---

# Getting Started

Your content here...
```

### Custom Components

Extend markdown with custom syntax:

```markdown
::: warning
This is a warning box
:::

::: tip
This is a helpful tip
:::

::: info
Additional information
:::
```

## Data Manipulation

### Accessing Generated Data

```bash
# Read content programmatically
cat docs/data/content.json | jq '.sections[0]'

# Count sections
jq '.metadata.total_sections' docs/data/content.json

# Extract all headings
jq '.sections[].heading' docs/data/content.json
```

### Modifying Content

```bash
# Update a section
jq '.sections[0].content = "New content"' \
  docs/data/sections/section-1.json > temp.json
mv temp.json docs/data/sections/section-1.json

# Rebuild master file
go run scripts/rebuild-content.go
```

### Exporting Data

```bash
# Export to CSV
jq -r '.sections[] | [.id, .heading] | @csv' \
  docs/data/content.json > sections.csv

# Export to Markdown
jq -r '.sections[] | "# \(.heading)\n\n\(.content)\n"' \
  docs/data/content.json > output.md
```

## API Development

### REST API Endpoints

Create custom endpoints:

```go
// Add to cmd/main.go

// Get all sections
http.HandleFunc("/api/sections", func(w http.ResponseWriter, r *http.Request) {
    // Load and return sections
})

// Get specific section
http.HandleFunc("/api/sections/{id}", func(w http.ResponseWriter, r *http.Request) {
    // Load and return specific section
})

// Search API
http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
    // Perform search and return results
})
```

### Webhook Integration

Trigger processing via webhooks:

```go
http.HandleFunc("/api/webhook/process", func(w http.ResponseWriter, r *http.Request) {
    // Validate webhook
    // Process new content
    // Return status
})
```

## Performance Optimization

### Caching Strategies

```go
// Cache parsed content
var contentCache *ContentData
var cacheTime time.Time

func getCachedContent() (*ContentData, error) {
    if time.Since(cacheTime) > 5*time.Minute {
        // Reload content
        contentCache = loadContent()
        cacheTime = time.Now()
    }
    return contentCache, nil
}
```

### Lazy Loading

```javascript
// Load sections on demand
async function loadSection(sectionId) {
    const response = await fetch(`/docs/data/sections/${sectionId}.json`);
    return await response.json();
}

// Load visible sections only
const observer = new IntersectionObserver(entries => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            loadSection(entry.target.dataset.sectionId);
        }
    });
});
```

## Custom Styling

### Theme Customization

```css
/* Create custom theme in static/theme.css */
:root {
    --primary-color: #6366f1;
    --primary-hover: #4f46e5;
    --background: #0f172a;
    --text-primary: #f1f5f9;
}
```

### Component Styling

```css
/* Custom chat widget */
.chat-widget {
    border-radius: 20px;
    box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

/* Custom search results */
.search-results {
    backdrop-filter: blur(10px);
    background: rgba(255,255,255,0.9);
}
```

## Integration Examples

### CI/CD Integration

```yaml
# .github/workflows/docs.yml
name: Generate Docs
on:
  push:
    paths:
      - 'input/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: make process-md
      - run: make deploy
```

### Docker Integration

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o doctrainer cmd/main.go

EXPOSE 8080
CMD ["./doctrainer", "-serve"]
```

### Static Site Generation

```bash
# Generate static site
go run cmd/main.go -pdf=input/doc.pdf

# Optimize for static hosting
# Remove server-side features
# Inline all assets

# Deploy to Netlify/Vercel
netlify deploy --dir=docs
```

## Monitoring and Analytics

### Logging

```go
// Add structured logging
import "log/slog"

slog.Info("Processing document",
    "type", "pdf",
    "path", pdfPath,
    "sections", len(sections))
```

### Metrics

```go
// Track metrics
var (
    processingTime = prometheus.NewHistogram(...)
    chatRequests   = prometheus.NewCounter(...)
)
```

## Security Best Practices

### Input Validation

```go
// Validate file paths
func validatePath(path string) error {
    if !strings.HasPrefix(filepath.Clean(path), "input/") {
        return errors.New("invalid path")
    }
    return nil
}
```

### Rate Limiting

```go
// Limit chat requests
limiter := rate.NewLimiter(10, 20) // 10 req/s, burst 20

http.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "Rate limit exceeded", 429)
        return
    }
    handleChat(w, r)
})
```

## Troubleshooting

### Debug Mode

```bash
# Enable verbose logging
go run cmd/main.go -debug -pdf=input/doc.pdf

# Log to file
go run cmd/main.go -serve 2>&1 | tee server.log
```

### Performance Profiling

```bash
# CPU profiling
go run cmd/main.go -cpuprofile=cpu.prof -pdf=input/doc.pdf

# Memory profiling
go run cmd/main.go -memprofile=mem.prof -pdf=input/doc.pdf

# Analyze profiles
go tool pprof cpu.prof
```

## Contributing

### Development Setup

```bash
# Clone repository
git clone https://github.com/subairkb777/docTrainerGO-pdf.git

# Create feature branch
git checkout -b feature/my-feature

# Make changes and test
make test

# Commit and push
git commit -m "Add new feature"
git push origin feature/my-feature
```

### Code Style

Follow Go conventions:

```go
// Good
func processMarkdown(file string) error {
    // Implementation
}

// Bad
func Process_markdown(FILE string) error {
    // Implementation
}
```

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Ollama Documentation](https://ollama.com/docs)
- [Markdown Guide](https://www.markdownguide.org/)

### Community
- [GitHub Issues](https://github.com/subairkb777/docTrainerGO-pdf/issues)
- [Discussions](https://github.com/subairkb777/docTrainerGO-pdf/discussions)

### Tools
- [Poppler](https://poppler.freedesktop.org/)
- [Fuse.js](https://fusejs.io/)
- [Air (Hot Reload)](https://github.com/cosmtrek/air)
