package cli

import (
	"flag"
	"fmt"
	"os"
)

// CLI represents command-line interface handler
type CLI struct {
	pdfPath        *string
	configPath     *string
	serve          *bool
	processAndExit *bool
	help           *bool
}

// New creates a new CLI instance
func New() *CLI {
	return &CLI{
		pdfPath:        flag.String("pdf", "", "Path to PDF file to process"),
		configPath:     flag.String("config", "config.yaml", "Path to configuration file"),
		serve:          flag.Bool("serve", false, "Start web server after processing"),
		processAndExit: flag.Bool("process", false, "Process document and exit (don't start server)"),
		help:           flag.Bool("help", false, "Show help message"),
	}
}

// Parse parses command-line flags
func (c *CLI) Parse() {
	flag.Parse()

	if *c.help {
		c.showHelp()
		os.Exit(0)
	}
}

// GetPDFPath returns the PDF path if provided
func (c *CLI) GetPDFPath() string {
	return *c.pdfPath
}

// GetConfigPath returns the configuration file path
func (c *CLI) GetConfigPath() string {
	return *c.configPath
}

// ShouldServe returns whether server should be started
func (c *CLI) ShouldServe() bool {
	return *c.serve
}

// ShouldProcessAndExit returns whether to process and exit
func (c *CLI) ShouldProcessAndExit() bool {
	return *c.processAndExit
}

// HasPDFPath returns whether a PDF path was provided
func (c *CLI) HasPDFPath() bool {
	return *c.pdfPath != ""
}

// showHelp displays usage information
func (c *CLI) showHelp() {
	fmt.Println("docTrainerGO - Generate searchable documentation with AI chat from PDF or Markdown files")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  docTrainerGO [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -config string")
	fmt.Println("        Path to configuration file (default: config.yaml)")
	fmt.Println("  -pdf string")
	fmt.Println("        Path to PDF file to process (overrides config)")
	fmt.Println("  -process")
	fmt.Println("        Process document and exit without starting server")
	fmt.Println("  -serve")
	fmt.Println("        Start web server after processing")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Process PDF and start server")
	fmt.Println("  docTrainerGO -pdf input/document.pdf -serve")
	fmt.Println()
	fmt.Println("  # Use markdown files from config and start server")
	fmt.Println("  docTrainerGO -serve")
	fmt.Println()
	fmt.Println("  # Process only without starting server")
	fmt.Println("  docTrainerGO -pdf input/document.pdf -process")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Edit config.yaml to configure:")
	fmt.Println("  - Input type (pdf or markdown)")
	fmt.Println("  - Markdown directory")
	fmt.Println("  - Output directory")
	fmt.Println("  - Ollama settings")
	fmt.Println("  - Server port")
	fmt.Println()
}
