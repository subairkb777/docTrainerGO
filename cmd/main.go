package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

// ContentData represents the structured content from data/content.json
type ContentData struct {
	Title    string        `json:"title"`
	Sections []SectionData `json:"sections"`
	Metadata struct {
		TotalSections int `json:"total_sections"`
		TotalImages   int `json:"total_images"`
	} `json:"metadata"`
}

// SectionData represents a section in content.json
type SectionData struct {
	ID      string   `json:"id"`
	Level   int      `json:"level"`
	Heading string   `json:"heading"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
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

	// Extract images using pdfimages (more reliable than Go libraries)
	fmt.Println("→ Extracting images from PDF...")
	if err := extractImagesWithPdfimages(pdfPath, docsDir); err != nil {
		fmt.Printf("  Warning: Image extraction failed: %v\n", err)
		fmt.Println("  Continuing without images...")
	}

	// Parse PDF
	fmt.Println("→ Parsing PDF and extracting content...")
	parser := pdf.NewParser(docsDir)
	doc, err := parser.Parse(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to parse PDF: %w", err)
	}
	fmt.Printf("  Found %d sections\n", len(doc.Sections))

	// Generate structured data files
	fmt.Println("→ Generating structured data...")
	dataGen := generator.NewDataGenerator(docsDir)
	if err := dataGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate data files: %w", err)
	}

	// Generate HTML
	fmt.Println("→ Generating HTML pages...")
	gen := generator.NewGenerator("templates/page.html", docsDir)
	if err := gen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate search index (keeping for backward compatibility)
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

// loadDocumentationContext loads the documentation content from data/content.json
func loadDocumentationContext() (string, error) {
	// Load the content.json file (new structured format)
	contentPath := filepath.Join(docsDir, "data", "content.json")
	file, err := os.Open(contentPath)
	if err != nil {
		return "", fmt.Errorf("failed to open content.json: %w", err)
	}
	defer file.Close()

	// Parse the content data
	var content ContentData
	if err := json.NewDecoder(file).Decode(&content); err != nil {
		return "", fmt.Errorf("failed to parse content.json: %w", err)
	}

	// Build context from all sections
	var contextBuilder strings.Builder
	contextBuilder.WriteString(fmt.Sprintf("=== %s ===\n\n", content.Title))
	contextBuilder.WriteString(fmt.Sprintf("Total Sections: %d | Total Images: %d\n\n",
		content.Metadata.TotalSections, content.Metadata.TotalImages))

	for _, section := range content.Sections {
		contextBuilder.WriteString(fmt.Sprintf("## %s\n", section.Heading))
		contextBuilder.WriteString(fmt.Sprintf("%s\n\n", section.Content))
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

// extractImagesWithPdfimages uses the pdfimages command-line tool to extract images
func extractImagesWithPdfimages(pdfPath, outputDir string) error {
	// Check if pdfimages is available
	if _, err := exec.LookPath("pdfimages"); err != nil {
		return fmt.Errorf("pdfimages not found (install with: brew install poppler)")
	}

	// Create images directory
	imageDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}

	// Extract images as PNG
	outputPrefix := filepath.Join(imageDir, "img")
	cmd := exec.Command("pdfimages", "-png", pdfPath, outputPrefix)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pdfimages failed: %w\nOutput: %s", err, string(output))
	}

	// Count extracted images
	entries, err := os.ReadDir(imageDir)
	if err != nil {
		return fmt.Errorf("failed to read image directory: %w", err)
	}

	imageCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".png") || strings.HasSuffix(entry.Name(), ".jpg")) {
			imageCount++
		}
	}

	fmt.Printf("  Extracted %d images\n", imageCount)
	return nil
}
