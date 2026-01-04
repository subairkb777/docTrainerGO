package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"docTrainerGO/internal/chat"
)

// Server represents the HTTP server
type Server struct {
	port         string
	docsDir      string
	ollamaClient *chat.OllamaClient
}

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

// New creates a new server instance
func New(port, docsDir string, ollamaClient *chat.OllamaClient) *Server {
	return &Server{
		port:         port,
		docsDir:      docsDir,
		ollamaClient: ollamaClient,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Serve static files
	staticDir := "static"
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve docs directory
	docsFS := http.FileServer(http.Dir(s.docsDir))
	http.Handle("/docs/", http.StripPrefix("/docs/", docsFS))

	// Serve main page
	http.HandleFunc("/", s.handleIndex)

	// Chat API endpoint
	http.HandleFunc("/api/chat", s.handleChat)

	// Start server
	addr := ":" + s.port
	fmt.Printf("\n🚀 Server running at http://localhost:%s\n", s.port)
	fmt.Println("   Press Ctrl+C to stop")
	fmt.Println()

	return http.ListenAndServe(addr, nil)
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, filepath.Join(s.docsDir, "index.html"))
		return
	}
	http.NotFound(w, r)
}

// handleChat processes chat requests from the frontend
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
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

	// Check if Ollama is available
	if s.ollamaClient == nil {
		s.respondWithError(w, "AI chat is disabled. Enable it in config.yaml (ollama.enabled: true)", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		s.respondWithError(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	// Load documentation context from search index
	context, err := s.loadDocumentationContext()
	if err != nil {
		log.Printf("Warning: Could not load documentation context: %v", err)
		context = "Documentation not available."
	}

	// Query Ollama with documentation context
	answer, err := s.ollamaClient.AskWithContext(req.Prompt, context)
	if err != nil {
		log.Printf("Ollama error: %v", err)
		s.respondWithError(w, "Failed to get response from AI", http.StatusInternalServerError)
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
func (s *Server) loadDocumentationContext() (string, error) {
	contentPath := filepath.Join(s.docsDir, "data", "content.json")
	file, err := os.Open(contentPath)
	if err != nil {
		return "", fmt.Errorf("failed to open content.json: %w", err)
	}
	defer file.Close()

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
func (s *Server) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	resp := ChatResponse{
		Error: message,
	}
	json.NewEncoder(w).Encode(resp)
}
