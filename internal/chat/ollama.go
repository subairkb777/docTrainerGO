package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaClient handles communication with local Ollama LLM
type OllamaClient struct {
	baseURL string
	model   string
	timeout time.Duration
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.1"
	}

	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		timeout: 60 * time.Second,
	}
}

// ChatRequest represents a chat request to Ollama
type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ChatResponse represents the response from Ollama
type ChatResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// Ask sends a question to the local Ollama LLM and returns the answer
func (c *OllamaClient) Ask(prompt string) (string, error) {
	// Prepare request
	reqBody := ChatRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/generate", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request with timeout
	client := &http.Client{
		Timeout: c.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w (is Ollama running?)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Clean up response
	answer := strings.TrimSpace(chatResp.Response)
	if answer == "" {
		return "I'm sorry, I couldn't generate a response. Please try again.", nil
	}

	return answer, nil
}

// AskWithContext sends a question with document context to the LLM
func (c *OllamaClient) AskWithContext(question, context string) (string, error) {
	prompt := fmt.Sprintf(`You are a helpful documentation assistant. Use the following context from the documentation to answer the user's question. If the answer is not in the context, say so.

Context:
%s

Question: %s

Answer:`, context, question)

	return c.Ask(prompt)
}

// HealthCheck verifies if Ollama is accessible
func (c *OllamaClient) HealthCheck() error {
	url := fmt.Sprintf("%s/api/tags", c.baseURL)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("Ollama is not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	return nil
}
