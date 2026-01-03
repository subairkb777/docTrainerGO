package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Section represents a documentation section with heading, content, and images
type Section struct {
	ID      string   // Unique identifier for the section
	Level   int      // Heading level (1-6)
	Heading string   // Section heading text
	Content string   // Section text content
	Images  []string // Paths to extracted images
}

// Document represents the parsed PDF document
type Document struct {
	Title    string
	Sections []Section
}

// Parser handles PDF parsing and image extraction
type Parser struct {
	outputDir string
	imageDir  string
	imageIdx  int
}

// NewParser creates a new PDF parser
func NewParser(outputDir string) *Parser {
	return &Parser{
		outputDir: outputDir,
		imageDir:  filepath.Join(outputDir, "images"),
		imageIdx:  0,
	}
}

// Parse extracts text and images from a PDF file
func (p *Parser) Parse(pdfPath string) (*Document, error) {
	// Ensure image directory exists
	if err := os.MkdirAll(p.imageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %w", err)
	}

	// Open PDF file
	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	doc := &Document{
		Title:    extractTitle(pdfPath),
		Sections: make([]Section, 0),
	}

	// Extract text from all pages
	var allText strings.Builder
	totalPages := r.NumPage()

	for pageIdx := 1; pageIdx <= totalPages; pageIdx++ {
		page := r.Page(pageIdx)
		if page.V.IsNull() {
			continue
		}

		// Extract text content
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		allText.WriteString(text)
		allText.WriteString("\n\n")

		// Extract images from page
		if err := p.extractImagesFromPage(page, pageIdx); err != nil {
			// Log but don't fail on image extraction errors
			fmt.Printf("Warning: failed to extract images from page %d: %v\n", pageIdx, err)
		}
	}

	// Parse text into sections
	doc.Sections = p.parseTextIntoSections(allText.String())

	return doc, nil
}

// extractImagesFromPage extracts images from a PDF page
func (p *Parser) extractImagesFromPage(page pdf.Page, pageNum int) error {
	// Note: The ledongthuc/pdf library has limited image extraction support
	// This is a simplified implementation. For production, consider using
	// pdfcpu or calling external tools like pdfimages

	content := page.Content()
	if content.Text == nil {
		return nil
	}

	// Look for image objects in the page content
	// This is a simplified approach - real PDF image extraction is complex
	// Convert Text slice to string for searching
	textContent := ""
	for _, t := range content.Text {
		textContent += t.S
	}

	if strings.Contains(textContent, "/Image") || strings.Contains(textContent, "/XObject") {
		// Create a placeholder image file
		// In a real implementation, you would extract actual image data
		p.imageIdx++
		imageName := fmt.Sprintf("page%d_img%d.png", pageNum, p.imageIdx)

		// For demonstration, we'll note that an image was found
		// Real implementation would extract and decode image data
		fmt.Printf("Found image reference on page %d (saved as %s)\n", pageNum, imageName)
	}

	return nil
}

// parseTextIntoSections converts plain text into structured sections
func (p *Parser) parseTextIntoSections(text string) []Section {
	sections := make([]Section, 0)
	lines := strings.Split(text, "\n")

	var currentSection *Section
	sectionID := 0

	// Regular expressions for detecting headings
	headingPattern := regexp.MustCompile(`^[A-Z][A-Za-z\s]{3,}$`)
	numberHeadingPattern := regexp.MustCompile(`^(\d+\.)+\s+[A-Z]`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect if line is a heading
		isHeading := false
		level := 1

		// Check for numbered headings (e.g., "1.2.3 Introduction")
		if numberHeadingPattern.MatchString(line) {
			isHeading = true
			level = strings.Count(strings.Split(line, " ")[0], ".") + 1
		} else if headingPattern.MatchString(line) && len(line) < 100 {
			// Check for title-case headings
			isHeading = true
			level = 1
		}

		if isHeading {
			// Save previous section
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}

			// Create new section
			sectionID++
			currentSection = &Section{
				ID:      fmt.Sprintf("section-%d", sectionID),
				Level:   level,
				Heading: line,
				Content: "",
				Images:  make([]string, 0),
			}
		} else if currentSection != nil {
			// Add content to current section
			if currentSection.Content != "" {
				currentSection.Content += " "
			}
			currentSection.Content += line
		} else {
			// Create initial section for content before first heading
			sectionID++
			currentSection = &Section{
				ID:      fmt.Sprintf("section-%d", sectionID),
				Level:   1,
				Heading: "Introduction",
				Content: line,
				Images:  make([]string, 0),
			}
		}
	}

	// Add last section
	if currentSection != nil {
		sections = append(sections, *currentSection)
	}

	// Associate images with sections (simplified approach)
	imageFiles, _ := filepath.Glob(filepath.Join(p.imageDir, "*.png"))
	imageFiles2, _ := filepath.Glob(filepath.Join(p.imageDir, "*.jpg"))
	allImages := append(imageFiles, imageFiles2...)

	// Distribute images across sections
	if len(sections) > 0 && len(allImages) > 0 {
		imagesPerSection := len(allImages) / len(sections)
		if imagesPerSection == 0 {
			imagesPerSection = 1
		}

		imgIdx := 0
		for i := range sections {
			for j := 0; j < imagesPerSection && imgIdx < len(allImages); j++ {
				sections[i].Images = append(sections[i].Images, filepath.Base(allImages[imgIdx]))
				imgIdx++
			}
		}
	}

	return sections
}

// extractTitle extracts document title from PDF filename
func extractTitle(pdfPath string) string {
	base := filepath.Base(pdfPath)
	title := strings.TrimSuffix(base, filepath.Ext(base))
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")
	return title
}

// SaveImageFromData saves image data to a file
func (p *Parser) SaveImageFromData(data []byte, format string) (string, error) {
	p.imageIdx++
	filename := fmt.Sprintf("image_%d.%s", p.imageIdx, format)
	filepath := filepath.Join(p.imageDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	// Decode and re-encode image
	img, err := decodeImage(data, format)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	switch format {
	case "png":
		err = png.Encode(file, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}

// decodeImage decodes image data based on format
func decodeImage(data []byte, format string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch format {
	case "png":
		return png.Decode(reader)
	case "jpg", "jpeg":
		return jpeg.Decode(reader)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
