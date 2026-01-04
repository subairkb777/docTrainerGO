package processor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"docTrainerGO/internal/config"
	"docTrainerGO/internal/generator"
	"docTrainerGO/internal/md"
	"docTrainerGO/internal/pdf"
	"docTrainerGO/internal/search"
)

// Processor handles document processing
type Processor struct {
	config *config.Config
}

// New creates a new processor
func New(cfg *config.Config) *Processor {
	return &Processor{config: cfg}
}

// Process processes documents based on configuration
func (p *Processor) Process() error {
	outputDir := p.config.Output.Directory

	var doc *pdf.Document
	var err error

	switch p.config.InputType {
	case "markdown":
		fmt.Println("Processing Markdown files...")
		doc, err = p.processMarkdown(outputDir)
	case "pdf":
		fmt.Println("Processing PDF file...")
		doc, err = p.processPDF(outputDir)
	default:
		return fmt.Errorf("invalid input_type: %s (must be 'pdf' or 'markdown')", p.config.InputType)
	}

	if err != nil {
		return err
	}

	fmt.Printf("  Found %d sections\n", len(doc.Sections))

	// Generate structured data
	fmt.Println("→ Generating structured data...")
	dataGen := generator.NewDataGenerator(outputDir)
	if err := dataGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate data files: %w", err)
	}

	// Generate HTML
	fmt.Println("→ Generating HTML pages...")
	gen := generator.NewGenerator("templates/page.html", outputDir)
	if err := gen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate search index
	fmt.Println("→ Creating search index...")
	indexGen := search.NewIndexGenerator(outputDir)
	if err := indexGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	return nil
}

// processMarkdown processes markdown files
func (p *Processor) processMarkdown(outputDir string) (*pdf.Document, error) {
	parser := md.NewParser(outputDir)

	var doc *pdf.Document
	var err error

	if p.config.Markdown.AutoDiscover {
		fmt.Printf("→ Auto-discovering files in: %s\n", p.config.Markdown.Directory)
		doc, err = parser.ParseDirectory(p.config.Markdown.Directory)
	} else {
		fmt.Printf("→ Processing %d specified files\n", len(p.config.Markdown.Files))
		doc, err = parser.ParseFiles(p.config.Markdown.Files)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Set title from config
	if p.config.Output.Title != "" {
		doc.Title = p.config.Output.Title
	}

	return doc, nil
}

// processPDF processes a PDF file
func (p *Processor) processPDF(outputDir string) (*pdf.Document, error) {
	pdfPath := p.config.PDF.Path
	if pdfPath == "" {
		return nil, fmt.Errorf("PDF path not specified in config")
	}

	fmt.Println("Processing PDF:", pdfPath)

	// Check if PDF exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("PDF file not found: %s", pdfPath)
	}

	// Extract images if enabled
	if p.config.PDF.ExtractImages {
		fmt.Println("→ Extracting images from PDF...")
		if err := extractImagesWithPdfimages(pdfPath, outputDir); err != nil {
			fmt.Printf("  Warning: Image extraction failed: %v\n", err)
			fmt.Println("  Continuing without images...")
		}
	}

	// Parse PDF
	fmt.Println("→ Parsing PDF and extracting content...")
	parser := pdf.NewParser(outputDir)
	doc, err := parser.Parse(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PDF: %w", err)
	}

	return doc, nil
}

// ProcessPDFDirect processes a PDF file directly (for CLI usage)
func ProcessPDFDirect(pdfPath, outputDir string) error {
	// Check if PDF exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file not found: %s", pdfPath)
	}

	// Extract images using pdfimages
	fmt.Println("→ Extracting images from PDF...")
	if err := extractImagesWithPdfimages(pdfPath, outputDir); err != nil {
		fmt.Printf("  Warning: Image extraction failed: %v\n", err)
		fmt.Println("  Continuing without images...")
	}

	// Parse PDF
	fmt.Println("→ Parsing PDF and extracting content...")
	parser := pdf.NewParser(outputDir)
	doc, err := parser.Parse(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to parse PDF: %w", err)
	}
	fmt.Printf("  Found %d sections\n", len(doc.Sections))

	// Generate structured data files
	fmt.Println("→ Generating structured data...")
	dataGen := generator.NewDataGenerator(outputDir)
	if err := dataGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate data files: %w", err)
	}

	// Generate HTML
	fmt.Println("→ Generating HTML pages...")
	gen := generator.NewGenerator("templates/page.html", outputDir)
	if err := gen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate search index
	fmt.Println("→ Creating search index...")
	indexGen := search.NewIndexGenerator(outputDir)
	if err := indexGen.Generate(doc); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	return nil
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
