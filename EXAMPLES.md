# Examples & Usage Guide

This guide provides detailed examples of using DocTrainerGO.

## Table of Contents

1. [Basic Usage](#basic-usage)
2. [Advanced Configuration](#advanced-configuration)
3. [Customization Examples](#customization-examples)
4. [API Examples](#api-examples)
5. [Troubleshooting Examples](#troubleshooting-examples)

---

## Basic Usage

### Example 1: Process a Simple PDF

```bash
# 1. Place your PDF in the input directory
cp ~/Downloads/user-guide.pdf input/

# 2. Process the PDF
go run cmd/main.go -pdf=input/user-guide.pdf

# Output:
# Processing PDF: input/user-guide.pdf
# → Parsing PDF and extracting content...
#   Found 15 sections
# → Generating HTML pages...
# Generated: docs/index.html
# → Creating search index...
# Generated search index: docs/search-index.json
# ✓ PDF processing complete!
```

### Example 2: Serve the Documentation

```bash
# Start the server
go run cmd/main.go -serve

# Output:
# ✓ Connected to Ollama
# 🚀 Server running at http://localhost:8080
#    Press Ctrl+C to stop
```

### Example 3: Complete Workflow with Makefile

```bash
# One-time setup
make setup

# Process PDF
make process PDF=input/documentation.pdf

# Start Ollama (separate terminal)
ollama run llama3.1

# Serve the site
make serve
```

---

## Advanced Configuration

### Example 4: Custom Port and Model

```bash
# Use custom port and Ollama model
go run cmd/main.go -serve -port=3000 -model=mistral

# Access at: http://localhost:3000
```

### Example 5: Remote Ollama Instance

```bash
# Connect to Ollama running on another machine
go run cmd/main.go -serve -ollama=http://192.168.1.100:11434
```

### Example 6: Build and Deploy Binary

```bash
# Build binary
go build -o doctrainer cmd/main.go

# Process PDF with binary
./doctrainer -pdf=input/manual.pdf

# Serve with binary
./doctrainer -serve -port=8080

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o doctrainer-linux cmd/main.go
GOOS=windows GOARCH=amd64 go build -o doctrainer.exe cmd/main.go
```

---

## Customization Examples

### Example 7: Customize CSS Styling

Edit `static/style.css` to change colors:

```css
/* Change primary color to green */
:root {
    --primary-color: #10b981;
    --primary-hover: #059669;
}

/* Change sidebar background */
.sidebar {
    background: linear-gradient(180deg, #f8fafc 0%, #e0e7ff 100%);
}
```

### Example 8: Modify HTML Template

Edit `templates/page.html` to add a footer:

```html
<!-- Add before closing </main> tag -->
<footer class="page-footer">
    <p>&copy; 2026 Your Company. All rights reserved.</p>
</footer>
```

### Example 9: Custom Search Options

Edit `static/script.js` to adjust search sensitivity:

```javascript
// Change Fuse.js options
const options = {
    keys: ['heading', 'content'],
    threshold: 0.3,  // More strict matching (default: 0.4)
    minMatchCharLength: 3  // Require 3 chars minimum
};
```

---

## API Examples

### Example 10: Chat API Usage with cURL

```bash
# Send a question to the AI assistant
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "What is this documentation about?"}'

# Response:
# {
#   "answer": "Based on the documentation, this appears to be..."
# }
```

### Example 11: Chat API with JavaScript

```javascript
// Frontend JavaScript example
async function askQuestion(question) {
    const response = await fetch('/api/chat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ prompt: question })
    });
    
    const data = await response.json();
    
    if (data.error) {
        console.error('Error:', data.error);
    } else {
        console.log('Answer:', data.answer);
    }
}

// Usage
askQuestion('How do I install this software?');
```

### Example 12: Chat API with Python

```python
import requests

def ask_ai(question):
    url = 'http://localhost:8080/api/chat'
    payload = {'prompt': question}
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    if 'error' in data:
        print(f"Error: {data['error']}")
    else:
        print(f"Answer: {data['answer']}")

# Usage
ask_ai('What are the key features?')
```

---

## Troubleshooting Examples

### Example 13: Debug PDF Processing

```bash
# Add verbose logging
go run cmd/main.go -pdf=input/document.pdf 2>&1 | tee process.log

# Check what was extracted
ls -la docs/images/
cat docs/search-index.json | jq '.items[] | .heading'
```

### Example 14: Test Ollama Connection

```bash
# Test Ollama directly
curl http://localhost:11434/api/tags

# Test with a simple prompt
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.1",
    "prompt": "Say hello",
    "stream": false
  }'
```

### Example 15: Check Server Logs

```bash
# Run server with verbose output
go run cmd/main.go -serve 2>&1 | tee server.log

# In another terminal, test the chat endpoint
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "test"}' -v
```

---

## Production Examples

### Example 16: Systemd Service (Linux)

Create `/etc/systemd/system/doctrainer.service`:

```ini
[Unit]
Description=DocTrainerGO Documentation Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/doctrainer
ExecStart=/opt/doctrainer/doctrainer -serve -port=8080
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable doctrainer
sudo systemctl start doctrainer
sudo systemctl status doctrainer
```

### Example 17: Nginx Reverse Proxy

Create `/etc/nginx/sites-available/doctrainer`:

```nginx
server {
    listen 80;
    server_name docs.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/chat {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_read_timeout 90s;
    }
}
```

### Example 18: Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  doctrainer:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./docs:/app/docs
      - ./input:/app/input
    environment:
      - OLLAMA_URL=http://ollama:11434
    depends_on:
      - ollama

  ollama:
    image: ollama/ollama
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama

volumes:
  ollama_data:
```

Run:

```bash
docker-compose up -d
```

---

## Performance Examples

### Example 19: Batch Processing

```bash
# Process multiple PDFs
for pdf in input/*.pdf; do
    echo "Processing $pdf..."
    go run cmd/main.go -pdf="$pdf"
done
```

### Example 20: Monitoring

```bash
# Check memory usage
ps aux | grep "go run"

# Monitor server with htop
htop -p $(pgrep -f "go run")

# Check port usage
lsof -i :8080
```

---

## Integration Examples

### Example 21: CI/CD Pipeline (.github/workflows/deploy.yml)

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Process PDF
        run: |
          go run cmd/main.go -pdf=input/docs.pdf
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

### Example 22: Webhook Integration

```go
// Add to cmd/main.go for webhook notifications
func notifyWebhook(message string) {
    payload := map[string]string{"text": message}
    data, _ := json.Marshal(payload)
    http.Post("https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
        "application/json", bytes.NewBuffer(data))
}
```

---

For more examples, visit the [GitHub repository](https://github.com/yourusername/docTrainerGO) or [open an issue](https://github.com/yourusername/docTrainerGO/issues).
