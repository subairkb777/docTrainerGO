package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"docTrainerGO/internal/pdf"
)

// Generator handles HTML page generation
type Generator struct {
	templatePath string
	outputDir    string
}

// NewGenerator creates a new HTML generator
func NewGenerator(templatePath, outputDir string) *Generator {
	return &Generator{
		templatePath: templatePath,
		outputDir:    outputDir,
	}
}

// PageData represents data passed to HTML templates (minimal - just metadata)
type PageData struct {
	Title string
}

// Generate creates a lightweight HTML shell that loads content dynamically
func (g *Generator) Generate(doc *pdf.Document) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare minimal page data (only title for initial load)
	pageData := PageData{
		Title: doc.Title,
	}

	// Parse template
	tmpl, err := template.ParseFiles(g.templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Generate lightweight index.html
	outputPath := filepath.Join(g.outputDir, "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, pageData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Generated: %s (lightweight - loads from content.json)\n", outputPath)
	return nil
}

// GenerateAll generates all necessary HTML files
func (g *Generator) GenerateAll(doc *pdf.Document) error {
	return g.Generate(doc)
}
