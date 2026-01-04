package main

import (
	"fmt"
	"log"
	"os"

	"docTrainerGO/internal/chat"
	"docTrainerGO/internal/cli"
	"docTrainerGO/internal/config"
	"docTrainerGO/internal/processor"
	"docTrainerGO/internal/server"
)

func main() {
	// Parse command line arguments
	commandLine := cli.New()
	commandLine.Parse()

	// Load configuration
	cfg, err := config.Load(commandLine.GetConfigPath())
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override config with CLI PDF path if provided
	if commandLine.HasPDFPath() {
		cfg.InputType = "pdf"
		cfg.PDF.Path = commandLine.GetPDFPath()
		fmt.Println("Processing PDF from command line...")
	}

	// Initialize Ollama client
	ollamaClient := chat.NewOllamaClient(cfg.Ollama.URL, cfg.Ollama.Model)

	// Process document
	proc := processor.New(cfg)
	if err := proc.Process(); err != nil {
		log.Fatalf("Failed to process document: %v", err)
	}

	fmt.Println("\n✓ Processing complete!")
	fmt.Println("✓ Documentation generated in:", cfg.Output.Directory)

	// Exit if process-only mode
	if commandLine.ShouldProcessAndExit() {
		fmt.Println("\nProcess-only mode: exiting without starting server")
		return
	}

	// Start server if requested
	if commandLine.ShouldServe() {
		// Check if docs directory exists
		if _, err := os.Stat(cfg.Output.Directory); os.IsNotExist(err) {
			log.Fatalf("Documentation directory '%s' does not exist", cfg.Output.Directory)
		}

		// Check Ollama health
		if err := ollamaClient.HealthCheck(); err != nil {
			log.Printf("Warning: Ollama health check failed: %v", err)
			log.Println("Chat functionality may not work. Make sure Ollama is running:")
			log.Printf("  ollama run %s\n", cfg.Ollama.Model)
		} else {
			fmt.Println("✓ Connected to Ollama")
		}

		// Start server
		srv := server.New(cfg.Server.Port, cfg.Output.Directory, ollamaClient)
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		return
	}

	// Show next steps
	fmt.Println("\nRun with -serve flag to start the web server:")
	fmt.Printf("  go run cmd/main.go -serve\n")
}
