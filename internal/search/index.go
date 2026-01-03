package search

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"docTrainerGO/internal/pdf"
)

// SearchIndex represents the search index structure for Fuse.js
type SearchIndex struct {
	Items []SearchItem `json:"items"`
}

// SearchItem represents a searchable item
type SearchItem struct {
	ID      string `json:"id"`
	Heading string `json:"heading"`
	Content string `json:"content"`
	Level   int    `json:"level"`
}

// IndexGenerator generates search indexes
type IndexGenerator struct {
	outputDir string
}

// NewIndexGenerator creates a new index generator
func NewIndexGenerator(outputDir string) *IndexGenerator {
	return &IndexGenerator{
		outputDir: outputDir,
	}
}

// Generate creates a JSON search index from the document
func (ig *IndexGenerator) Generate(doc *pdf.Document) error {
	// Build search items from document sections
	items := make([]SearchItem, 0, len(doc.Sections))

	for _, section := range doc.Sections {
		// Truncate content for search preview (first 200 chars)
		content := section.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}

		items = append(items, SearchItem{
			ID:      section.ID,
			Heading: section.Heading,
			Content: content,
			Level:   section.Level,
		})
	}

	index := SearchIndex{
		Items: items,
	}

	// Write JSON file
	indexPath := filepath.Join(ig.outputDir, "search-index.json")
	file, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to create search index file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(index); err != nil {
		return fmt.Errorf("failed to encode search index: %w", err)
	}

	fmt.Printf("Generated search index: %s\n", indexPath)
	return nil
}
