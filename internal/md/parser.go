package md

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"docTrainerGO/internal/pdf"
)

// Parser handles Markdown file parsing
type Parser struct {
	outputDir string
	imageDir  string
	sectionID int
}

// NewParser creates a new Markdown parser
func NewParser(outputDir string) *Parser {
	return &Parser{
		outputDir: outputDir,
		imageDir:  filepath.Join(outputDir, "images"),
		sectionID: 0,
	}
}

// ParseFiles processes multiple markdown files
func (p *Parser) ParseFiles(files []string) (*pdf.Document, error) {
	// Ensure image directory exists
	if err := os.MkdirAll(p.imageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %w", err)
	}

	doc := &pdf.Document{
		Title:    "Documentation",
		Sections: make([]pdf.Section, 0),
	}

	// Process each markdown file
	for _, file := range files {
		sections, err := p.parseFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}
		doc.Sections = append(doc.Sections, sections...)
	}

	return doc, nil
}

// parseFile parses a single markdown file
func (p *Parser) parseFile(filePath string) ([]pdf.Section, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sections := make([]pdf.Section, 0)
	scanner := bufio.NewScanner(file)

	var currentSection *pdf.Section
	var contentBuilder strings.Builder
	var inCodeBlock bool
	var inFrontMatter bool
	lineNum := 0

	// Regular expressions - using raw strings
	headingRegex := regexp.MustCompile("^(#{1,6})\\s+(.+)$")
	imageRegex := regexp.MustCompile("!\\[([^\\]]*)\\]\\(([^)]+)\\)")
	codeBlockRegex := regexp.MustCompile("^```")
	frontMatterRegex := regexp.MustCompile("^---$")

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Handle front matter
		if lineNum == 1 && frontMatterRegex.MatchString(line) {
			inFrontMatter = true
			continue
		}
		if inFrontMatter {
			if frontMatterRegex.MatchString(line) {
				inFrontMatter = false
			}
			continue
		}

		// Handle code blocks
		if codeBlockRegex.MatchString(line) {
			inCodeBlock = !inCodeBlock
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
			continue
		}

		if inCodeBlock {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
			continue
		}

		// Check for heading
		if matches := headingRegex.FindStringSubmatch(line); matches != nil {
			// Save previous section
			if currentSection != nil {
				currentSection.Content = strings.TrimSpace(contentBuilder.String())
				sections = append(sections, *currentSection)
				contentBuilder.Reset()
			}

			// Create new section
			p.sectionID++
			level := len(matches[1])
			heading := strings.TrimSpace(matches[2])

			currentSection = &pdf.Section{
				ID:      fmt.Sprintf("section-%d", p.sectionID),
				Level:   level,
				Heading: heading,
				Images:  make([]string, 0),
			}
			continue
		}

		// Extract images from line
		if imageMatches := imageRegex.FindAllStringSubmatch(line, -1); imageMatches != nil {
			for _, match := range imageMatches {
				imagePath := match[2]
				if currentSection != nil {
					// Copy image to output directory
					if err := p.copyImage(imagePath, filePath); err == nil {
						imageName := filepath.Base(imagePath)
						currentSection.Images = append(currentSection.Images, imageName)
					}
				}
			}
		}

		// Add line to content
		if currentSection != nil {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString(" ")
		}
	}

	// Save last section
	if currentSection != nil {
		currentSection.Content = strings.TrimSpace(contentBuilder.String())
		sections = append(sections, *currentSection)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// copyImage copies an image from source to output directory
func (p *Parser) copyImage(imagePath string, markdownFile string) error {
	// Resolve relative paths
	baseDir := filepath.Dir(markdownFile)
	sourcePath := filepath.Join(baseDir, imagePath)

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("image not found: %s", sourcePath)
	}

	// Read source
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Write to destination
	imageName := filepath.Base(imagePath)
	destPath := filepath.Join(p.imageDir, imageName)

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return err
	}

	return nil
}

// ParseDirectory processes all markdown files in a directory
func (p *Parser) ParseDirectory(dir string) (*pdf.Document, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && info.Name() != "README.md" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no markdown files found in %s", dir)
	}

	fmt.Printf("Found %d markdown files\n", len(files))
	return p.ParseFiles(files)
}
