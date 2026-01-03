package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"docTrainerGO/internal/chat"
	"docTrainerGO/internal/generator"
	"docTrainerGO/internal/pdf"
	"docTrainerGO/internal/search"
)

const (
	defaultPort     = "8080"
	defaultDocsDir  = "docs"
	defaultInputDir = "input"
)

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Prompt string `json:"prompt"`
}

// ChatResponse represents the chat response
type ChatResponse struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

// SearchItem represents an item in the search index
type SearchItem struct {
	ID      string `json:"id"`
	Heading string `json:"heading"`
	Content string `json:"content"`
	Level   int    `json:"level"`
}

// SearchIndex represents the search index structure
type SearchIndex struct {
	Items []SearchItem `json:"items"`
}

var (
	ollamaClient *chat.OllamaClient
	docsDir      string
)

func main() {
	// Parse command line flags
	port := flag.String("port", defaultPort, "Port to serve on")
	pdfPath := flag.String("pdf", "", "Path to PDF file to process")
	serve := flag.Bool("serve", false, "Serve the documentation site")
	ollamaURL := flag.String("ollama", "http://localhost:11434", "Ollama API URL")
	model := flag.String("model", "llama3.2", "Ollama model name")
	flag.Parse()

	docsDir = defaultDocsDir

	// Initialize Ollama client
	ollamaClient = chat.NewOllamaClient(*ollamaURL, *model)

	// Process PDF if provided
	if *pdfPath != "" {
		if err := processPDF(*pdfPath); err != nil {
			log.Fatalf("Failed to process PDF: %v", err)
		}
		fmt.Println("\n✓ PDF processing complete!")
		fmt.Println("✓ Documentation generated in:", docsDir)
		fmt.Println("\nRun with -serve flag to start the web server:")
		fmt.Printf("  go run cmd/main.go -serve -port=%s\n", *port)
		return
	}

	// Serve documentation site
	if *serve {
		// Check if docs directory exists
		if _, err := os.Stat(docsDir); os.IsNotExist(err) {
			log.Fatalf("Documentation directory '%s' does not exist. Please process a PDF first.", docsDir)
		}

		// Check Ollama health
		if err := ollamaClient.HealthCheck(); err != nil {
			log.Printf("Warning: Ollama health check failed: %v", err)
			log.Println("Chat functionality may not work. Make sure Ollama is running:")
			log.Println("  ollama run", *model)
		} else {
			fmt.Println("✓ Connected to Ollama")
		}

		serveDocs(*port)
		return
	}

	// Show usage if no flags provided
	fmt.Println("DocTrainerGO - PDF to Documentation Website Generator")
	fmt.Println("\nUsage:")
	fmt.Println("  Process PDF:")
	fmt.Println("    go run cmd/main.go -pdf=input/document.pdf")
	fmt.Println("\n  Serve documentation:")
	fmt.Println("    go run cmd/main.go -serve")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

// processPDF converts a PDF file to documentation website
func processPDF(pdfPath string) error {
	fmt.Println("Processing PDF:", pdfPath)

	// Check if PDF exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file not found: %s", pdfPath)
	}

	// Parse PDF
	fmt.Println("→ Parsing PDF and extracting content...")
	parser := pdf.NewParser(docsDir)
	doc, err := parser.Parse(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to parse PDF: %w", err)
	}
	fmt.Printf("  Found %d sections\n", len(doc.Sections))

	// Generate HTML
	fmt.Println("→ Generating HTML pages...")
	gen := generator.NewGenerator("templates/page.html", docsDir)
	if err := gen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate search index
	fmt.Println("→ Creating search index...")
	indexGen := search.NewIndexGenerator(docsDir)
	if err := indexGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	return nil
}

// serveDocs starts the HTTP server to serve documentation
func serveDocs(port string) {
	// Serve static files
	staticDir := "static"
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve docs directory
	docsFS := http.FileServer(http.Dir(docsDir))
	http.Handle("/docs/", http.StripPrefix("/docs/", docsFS))

	// Serve main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(docsDir, "index.html"))
			return
		}
		http.NotFound(w, r)
	})

	// Chat API endpoint
	http.HandleFunc("/api/chat", handleChat)

	// Start server
	addr := ":" + port
	fmt.Printf("\n🚀 Server running at http://localhost:%s\n", port)
	fmt.Println("   Press Ctrl+C to stop")
	fmt.Println()

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// handleChat processes chat requests from the frontend
func handleChat(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		respondWithError(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	// Load documentation context from search index
	context, err := loadDocumentationContext()
	if err != nil {
		log.Printf("Warning: Could not load documentation context: %v", err)
		// Fall back to asking without context
		context = "Documentation not available."
	}

	// Query Ollama with documentation context
	answer, err := ollamaClient.AskWithContext(req.Prompt, context)
	if err != nil {
		log.Printf("Ollama error: %v", err)
		respondWithError(w, "Failed to get response from AI", http.StatusInternalServerError)
		return
	}

	// Send response
	resp := ChatResponse{
		Answer: answer,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// loadDocumentationContext loads the documentation content from search index
func loadDocumentationContext() (string, error) {
	// Load the search index JSON
	indexPath := filepath.Join(docsDir, "search-index.json")
	file, err := os.Open(indexPath)
	if err != nil {
		return "", fmt.Errorf("failed to open search index: %w", err)
	}
	defer file.Close()

	// Parse the search index
	var index SearchIndex
	if err := json.NewDecoder(file).Decode(&index); err != nil {
		return "", fmt.Errorf("failed to parse search index: %w", err)
	}

	// Build context from all sections
	var contextBuilder strings.Builder
	contextBuilder.WriteString("=== DOCUMENTATION CONTENT ===\n\n")

	for _, item := range index.Items {
		contextBuilder.WriteString(fmt.Sprintf("## %s\n", item.Heading))
		contextBuilder.WriteString(fmt.Sprintf("%s\n\n", item.Content))
	}

	context := contextBuilder.String()

	// Limit context size to avoid token limits (keep first 15000 chars)
	if len(context) > 15000 {
		context = context[:15000] + "\n\n... (documentation continues)"
	}

	return context, nil
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	resp := ChatResponse{
		Error: message,
	}
	json.NewEncoder(w).Encode(resp)
}
