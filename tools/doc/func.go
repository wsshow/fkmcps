package doc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wsshow/docreader"
)

// GetDocumentInfoRequest Document information request
type GetDocumentInfoRequest struct {
	FilePath string `json:"file_path" jsonschema:"required,description:Document file path (supports .docx, .pdf, .xlsx, .pptx, .txt, .csv, .md, .rtf)"`
}

// GetDocumentInfoResponse Document information response
type GetDocumentInfoResponse struct {
	FilePath      string            `json:"file_path" jsonschema:"description:File path"`
	FileType      string            `json:"file_type" jsonschema:"description:File type (e.g., PDF, DOCX, XLSX)"`
	FileSize      string            `json:"file_size" jsonschema:"description:File size (formatted string)"`
	TotalPages    int               `json:"total_pages,omitempty" jsonschema:"description:Total number of pages (for paginated documents like PDF, PPTX)"`
	TotalSheets   int               `json:"total_sheets,omitempty" jsonschema:"description:Total number of sheets (XLSX only)"`
	SheetNames    []string          `json:"sheet_names,omitempty" jsonschema:"description:Sheet name list (XLSX only)"`
	EstimatedSize string            `json:"estimated_size" jsonschema:"description:Estimated text size (for evaluating whether it fits in one read)"`
	Metadata      map[string]string `json:"metadata" jsonschema:"description:Document metadata (title, author, etc.)"`
	ErrorMessage  string            `json:"error_message,omitempty" jsonschema:"description:Error message"`
}

// ReadDocumentByPagesRequest Read document by pages request
type ReadDocumentByPagesRequest struct {
	FilePath  string `json:"file_path" jsonschema:"required,description:Document file path"`
	StartPage int    `json:"start_page,omitempty" jsonschema:"description:Starting page number (0-based, default 0)"`
	EndPage   int    `json:"end_page,omitempty" jsonschema:"description:Ending page number (inclusive, -1 means to end, default -1)"`
}

// ReadDocumentByPagesResponse Read document by pages response
type ReadDocumentByPagesResponse struct {
	Content      string              `json:"content" jsonschema:"description:Read document content"`
	Pages        []PageContentDetail `json:"pages" jsonschema:"description:Page content details"`
	TotalPages   int                 `json:"total_pages" jsonschema:"description:Total number of pages in document"`
	ReadPages    int                 `json:"read_pages" jsonschema:"description:Actual number of pages read"`
	Metadata     map[string]string   `json:"metadata" jsonschema:"description:Document metadata"`
	ErrorMessage string              `json:"error_message,omitempty" jsonschema:"description:Error message"`
}

// PageContentDetail Page content detail
type PageContentDetail struct {
	PageNumber int    `json:"page_number" jsonschema:"description:Page number (0-based)"`
	PageName   string `json:"page_name,omitempty" jsonschema:"description:Page name (e.g., sheet name)"`
	LineCount  int    `json:"line_count" jsonschema:"description:Number of lines in this page"`
}

// ReadDocumentByLinesRequest Read document by lines request
type ReadDocumentByLinesRequest struct {
	FilePath  string `json:"file_path" jsonschema:"required,description:Document file path"`
	StartLine int    `json:"start_line,omitempty" jsonschema:"description:Starting line number (0-based, default 0)"`
	EndLine   int    `json:"end_line,omitempty" jsonschema:"description:Ending line number (inclusive, -1 means to end, default -1)"`
	PageIndex int    `json:"page_index,omitempty" jsonschema:"description:Specified page index (0-based, only valid for multi-page documents, -1 means first page, default -1)"`
}

// ReadDocumentByLinesResponse Read document by lines response
type ReadDocumentByLinesResponse struct {
	Content      string            `json:"content" jsonschema:"description:Read document content"`
	TotalLines   int               `json:"total_lines" jsonschema:"description:Total number of lines in this page"`
	ReadLines    int               `json:"read_lines" jsonschema:"description:Actual number of lines read"`
	PageIndex    int               `json:"page_index" jsonschema:"description:Page index read"`
	Metadata     map[string]string `json:"metadata" jsonschema:"description:Document metadata"`
	ErrorMessage string            `json:"error_message,omitempty" jsonschema:"description:Error message"`
}

// ReadDocumentSmartRequest Smart document reading request
type ReadDocumentSmartRequest struct {
	FilePath     string `json:"file_path" jsonschema:"required,description:Document file path"`
	MaxChars     int    `json:"max_chars,omitempty" jsonschema:"description:Maximum character limit (default 50000, recommended between 10000-100000)"`
	SampleMode   bool   `json:"sample_mode,omitempty" jsonschema:"description:Sampling mode (true for uniform sampling throughout, false for reading from start, default false)"`
	CleanContent bool   `json:"clean_content,omitempty" jsonschema:"description:Whether to clean text (remove extra spaces, blank lines, etc., default true)"`
}

// ReadDocumentSmartResponse Smart document reading response
type ReadDocumentSmartResponse struct {
	Content      string            `json:"content" jsonschema:"description:Read document content"`
	IsTruncated  bool              `json:"is_truncated" jsonschema:"description:Whether content is truncated"`
	OriginalSize int               `json:"original_size" jsonschema:"description:Original text size (character count)"`
	ReturnedSize int               `json:"returned_size" jsonschema:"description:Returned text size (character count)"`
	Strategy     string            `json:"strategy" jsonschema:"description:Reading strategy used"`
	Metadata     map[string]string `json:"metadata" jsonschema:"description:Document metadata"`
	ErrorMessage string            `json:"error_message,omitempty" jsonschema:"description:Error message"`
	Suggestion   string            `json:"suggestion,omitempty" jsonschema:"description:Suggestion (how to better read this document)"`
}

// GetDocumentInfo Get document basic information
func GetDocumentInfo(ctx context.Context, req *GetDocumentInfoRequest) (*GetDocumentInfoResponse, error) {
	// Check if file exists
	fileInfo, err := os.Stat(req.FilePath)
	if err != nil {
		return &GetDocumentInfoResponse{
			ErrorMessage: fmt.Sprintf("File access failed: %v", err),
		}, nil
	}

	// Get file type
	ext := strings.ToLower(filepath.Ext(req.FilePath))
	fileType := strings.TrimPrefix(ext, ".")
	fileType = strings.ToUpper(fileType)

	// Format file size
	fileSize := formatFileSize(fileInfo.Size())

	// Get metadata
	doc, err := docreader.ReadDocument(req.FilePath)
	if err != nil {
		return &GetDocumentInfoResponse{
			FilePath:     req.FilePath,
			FileType:     fileType,
			FileSize:     fileSize,
			ErrorMessage: fmt.Sprintf("Failed to read document: %v", err),
		}, nil
	}

	response := &GetDocumentInfoResponse{
		FilePath:      req.FilePath,
		FileType:      fileType,
		FileSize:      fileSize,
		EstimatedSize: fmt.Sprintf("Approx %d characters", len(doc.Content)),
		Metadata:      doc.Metadata,
	}

	// Get additional information based on file type
	switch ext {
	case ".pdf":
		if pages, ok := doc.Metadata["pages"]; ok {
			fmt.Sscanf(pages, "%d", &response.TotalPages)
		}
	case ".pptx":
		if slides, ok := doc.Metadata["slide_count"]; ok {
			fmt.Sscanf(slides, "%d", &response.TotalPages)
		}
	case ".xlsx":
		if sheets, ok := doc.Metadata["sheets"]; ok {
			response.SheetNames = strings.Split(sheets, ",")
			response.TotalSheets = len(response.SheetNames)
		}
	}

	return response, nil
}

// ReadDocumentByPages Read document by page range
func ReadDocumentByPages(ctx context.Context, req *ReadDocumentByPagesRequest) (*ReadDocumentByPagesResponse, error) {
	// Set default values
	if req.EndPage < 0 {
		req.EndPage = 999999 // Set a large value to indicate reading to the end
	}

	// Create read configuration
	config := docreader.NewReadConfig().WithPageRange(req.StartPage, req.EndPage)

	// Read document
	result, err := docreader.ReadDocumentWithConfig(req.FilePath, config)
	if err != nil {
		return &ReadDocumentByPagesResponse{
			ErrorMessage: fmt.Sprintf("Failed to read document: %v", err),
		}, nil
	}

	// Build response
	response := &ReadDocumentByPagesResponse{
		Content:    result.Content,
		TotalPages: result.TotalPages,
		ReadPages:  len(result.Pages),
		Metadata:   result.Metadata,
	}

	// Fill in page details
	response.Pages = make([]PageContentDetail, len(result.Pages))
	for i, page := range result.Pages {
		response.Pages[i] = PageContentDetail{
			PageNumber: page.PageNumber,
			PageName:   page.PageName,
			LineCount:  page.TotalLines,
		}
	}

	return response, nil
}

// ReadDocumentByLines Read document by line range
func ReadDocumentByLines(ctx context.Context, req *ReadDocumentByLinesRequest) (*ReadDocumentByLinesResponse, error) {
	pageIndex := req.PageIndex
	if pageIndex < 0 {
		pageIndex = 0
	}

	// Set default values
	endLine := req.EndLine
	if endLine < 0 {
		endLine = 999999 // Set a large value to indicate reading to the end
	}

	// Create read configuration
	config := docreader.NewReadConfig().
		AddPageLineRange(pageIndex, req.StartLine, endLine)

	// Read document
	result, err := docreader.ReadDocumentWithConfig(req.FilePath, config)
	if err != nil {
		return &ReadDocumentByLinesResponse{
			ErrorMessage: fmt.Sprintf("Failed to read document: %v", err),
		}, nil
	}

	// Build response
	response := &ReadDocumentByLinesResponse{
		Content:   result.Content,
		PageIndex: pageIndex,
		Metadata:  result.Metadata,
	}

	// Get line information for this page
	if len(result.Pages) > 0 {
		page := result.Pages[0]
		response.TotalLines = page.TotalLines
		response.ReadLines = len(page.Lines)
	}

	return response, nil
}

// ReadDocumentSmart Smart document reading (automatically adapts to context limitations)
func ReadDocumentSmart(ctx context.Context, req *ReadDocumentSmartRequest) (*ReadDocumentSmartResponse, error) {
	// Set default values
	maxChars := req.MaxChars
	if maxChars <= 0 {
		maxChars = 50000 // Default 50k characters
	}

	cleanContent := req.CleanContent
	if !req.CleanContent && maxChars == 50000 { // If user hasn't set, enable cleaning by default
		cleanContent = true
	}

	// First read the complete document
	var doc *docreader.Document
	var err error

	if cleanContent {
		doc, err = docreader.ReadDocumentWithClean(req.FilePath)
	} else {
		doc, err = docreader.ReadDocument(req.FilePath)
	}

	if err != nil {
		return &ReadDocumentSmartResponse{
			ErrorMessage: fmt.Sprintf("Failed to read document: %v", err),
		}, nil
	}

	originalSize := len(doc.Content)

	response := &ReadDocumentSmartResponse{
		OriginalSize: originalSize,
		Metadata:     doc.Metadata,
	}

	// If document size is within limit, return all content directly
	if originalSize <= maxChars {
		response.Content = doc.Content
		response.IsTruncated = false
		response.ReturnedSize = originalSize
		response.Strategy = "Complete read"
		return response, nil
	}

	// Document is too large, needs truncation
	response.IsTruncated = true

	if req.SampleMode {
		// Sampling mode: uniformly sample from different positions in the document
		response.Content = sampleContent(doc.Content, maxChars)
		response.Strategy = "Uniform sampling"
		response.Suggestion = fmt.Sprintf("Document is large (%d characters), key parts have been sampled. Recommend using ReadDocumentByPages or ReadDocumentByLines to read specific parts as needed", originalSize)
	} else {
		// Default: read from beginning
		response.Content = doc.Content[:maxChars]
		response.Strategy = "Truncate from start"
		response.Suggestion = fmt.Sprintf("Document is large (%d characters), only returning first %d characters. Recommend using ReadDocumentByPages or ReadDocumentByLines to read as needed", originalSize, maxChars)
	}

	response.ReturnedSize = len(response.Content)

	return response, nil
}

// formatFileSize Format file size
func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// sampleContent Uniformly sample from content
func sampleContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}

	// Divide content into 3 parts: beginning, middle, end
	// Reserve some space for separators
	const separator1 = "\n\n... [middle section] ...\n\n"
	const separator2 = "\n\n... [later section] ...\n\n"
	separatorLen := len(separator1) + len(separator2)

	availableChars := maxChars - separatorLen
	if availableChars < 300 { // Need at least 300 characters
		return content[:maxChars]
	}

	partSize := availableChars / 3

	start := content[:partSize]
	middleStart := len(content)/2 - partSize/2
	middleEnd := len(content)/2 + partSize/2
	if middleEnd > len(content) {
		middleEnd = len(content)
	}
	middle := content[middleStart:middleEnd]

	endStart := len(content) - partSize
	if endStart < 0 {
		endStart = 0
	}
	end := content[endStart:]

	result := fmt.Sprintf("%s%s%s%s%s", start, separator1, middle, separator2, end)

	// Ensure not exceeding the limit
	if len(result) > maxChars {
		return result[:maxChars]
	}

	return result
}
