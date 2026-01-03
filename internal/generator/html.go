package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

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

// PageData represents data passed to HTML templates
type PageData struct {
	Title           string
	Sections        []pdf.Section
	NavigationItems []NavItem
}

// NavItem represents a navigation menu item
type NavItem struct {
	ID      string
	Heading string
	Level   int
}

// Generate creates HTML pages from the parsed document
func (g *Generator) Generate(doc *pdf.Document) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build navigation items from sections
	navItems := make([]NavItem, len(doc.Sections))
	for i, section := range doc.Sections {
		navItems[i] = NavItem{
			ID:      section.ID,
			Heading: section.Heading,
			Level:   section.Level,
		}
	}

	// Prepare page data
	pageData := PageData{
		Title:           doc.Title,
		Sections:        doc.Sections,
		NavigationItems: navItems,
	}

	// Parse template
	tmpl, err := template.New("page.html").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"formatContent": func(s string) template.HTML {
			// Convert line breaks to paragraphs
			paragraphs := strings.Split(s, ". ")
			var formatted strings.Builder
			for i, p := range paragraphs {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				if i < len(paragraphs)-1 && !strings.HasSuffix(p, ".") {
					p += "."
				}
				formatted.WriteString("<p>")
				formatted.WriteString(template.HTMLEscapeString(p))
				formatted.WriteString("</p>\n")
			}
			return template.HTML(formatted.String())
		},
	}).ParseFiles(g.templatePath)

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Generate index.html
	outputPath := filepath.Join(g.outputDir, "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, pageData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Generated: %s\n", outputPath)
	return nil
}

// GenerateAll generates all necessary HTML files
func (g *Generator) GenerateAll(doc *pdf.Document) error {
	return g.Generate(doc)
}
