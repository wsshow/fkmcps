package doc

import (
	"fkmcps/structs"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GetTools(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name: "get_document_info",
		Description: `Get basic information about a document, including file type, size, page count, and metadata. Supported formats: .docx, .pdf, .xlsx, .pptx, .txt, .csv, .md, .rtf.
This is the first step before reading a document, helping you understand document structure and decide how to read it.`,
	}, structs.WarpToolFunc(GetDocumentInfo))

	mcp.AddTool(s, &mcp.Tool{
		Name: "read_document_smart",
		Description: `Intelligently read document content, automatically handling large documents. Features:
- Automatically adapts to context limits (default 50000 characters)
- Supports sampling mode (uniform sampling) or truncation from start
- Automatically cleans extra spaces and blank lines
- Provides suggestions when document is too large
Best for: First-time document reading, quick content overview`,
	}, structs.WarpToolFunc(ReadDocumentSmart))

	mcp.AddTool(s, &mcp.Tool{
		Name: "read_document_by_page",
		Description: `Read document content by page range. Supports multi-page documents like PDF, PPTX.
Parameters:
- start_page: Starting page (0-based)
- end_page: Ending page (-1 means to end)
Returns detailed information for each page (page number, line count).
Best for: Reading specific pages or sections`,
	}, structs.WarpToolFunc(ReadDocumentByPages))

	mcp.AddTool(s, &mcp.Tool{
		Name: "read_document_by_line",
		Description: `Read document content by line range. Supports specifying page and line range.
Parameters:
- start_line: Starting line (0-based)
- end_line: Ending line (-1 means to end)
- page_index: Page index (-1 means first page)
Best for: Reading specific lines or paragraphs`,
	}, structs.WarpToolFunc(ReadDocumentByLines))
}
