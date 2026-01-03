# Place Your PDF Files Here

This directory is for input PDF files that you want to convert to documentation.

## Usage

1. Copy your PDF file to this directory:
   ```bash
   cp /path/to/your/document.pdf input/
   ```

2. Process the PDF:
   ```bash
   go run cmd/main.go -pdf=input/document.pdf
   ```

## Supported PDF Features

- Text content (all pages)
- Headings and subheadings
- Paragraphs and sections
- Images (extracted and embedded)
- Tables (as text)

## Tips

- Use descriptive filenames for your PDFs
- Ensure PDFs are not password-protected
- Large PDFs (>50MB) may take longer to process
- PDFs with many images will create larger documentation sites
