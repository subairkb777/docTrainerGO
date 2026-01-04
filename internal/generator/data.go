package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"docTrainerGO/internal/pdf"
)

// ContentData represents the structured content storage
type ContentData struct {
	Title    string           `json:"title"`
	Sections []SectionData    `json:"sections"`
	Metadata DocumentMetadata `json:"metadata"`
}

// SectionData represents a single section with all its data
type SectionData struct {
	ID      string   `json:"id"`
	Level   int      `json:"level"`
	Heading string   `json:"heading"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

// DocumentMetadata contains document-level information
type DocumentMetadata struct {
	TotalSections int `json:"total_sections"`
	TotalImages   int `json:"total_images"`
}

// DataGenerator handles structured data generation
type DataGenerator struct {
	outputDir string
}

// NewDataGenerator creates a new data generator
func NewDataGenerator(outputDir string) *DataGenerator {
	return &DataGenerator{
		outputDir: outputDir,
	}
}

// Generate creates structured JSON data files
func (dg *DataGenerator) Generate(doc *pdf.Document) error {
	// Create data directory
	dataDir := filepath.Join(dg.outputDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Convert sections to data format
	sections := make([]SectionData, len(doc.Sections))
	totalImages := 0
	for i, section := range doc.Sections {
		sections[i] = SectionData{
			ID:      section.ID,
			Level:   section.Level,
			Heading: section.Heading,
			Content: section.Content,
			Images:  section.Images,
		}
		totalImages += len(section.Images)
	}

	// Create content data
	contentData := ContentData{
		Title:    doc.Title,
		Sections: sections,
		Metadata: DocumentMetadata{
			TotalSections: len(sections),
			TotalImages:   totalImages,
		},
	}

	// Save main content.json
	contentPath := filepath.Join(dataDir, "content.json")
	if err := dg.saveJSON(contentPath, contentData); err != nil {
		return fmt.Errorf("failed to save content.json: %w", err)
	}

	fmt.Printf("Generated: %s\n", contentPath)

	// Optionally save individual section files for easier maintenance
	sectionsDir := filepath.Join(dataDir, "sections")
	if err := os.MkdirAll(sectionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sections directory: %w", err)
	}

	for _, section := range sections {
		sectionPath := filepath.Join(sectionsDir, fmt.Sprintf("%s.json", section.ID))
		if err := dg.saveJSON(sectionPath, section); err != nil {
			return fmt.Errorf("failed to save section %s: %w", section.ID, err)
		}
	}

	fmt.Printf("Generated: %d individual section files\n", len(sections))

	return nil
}

// saveJSON writes data to a JSON file
func (dg *DataGenerator) saveJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
