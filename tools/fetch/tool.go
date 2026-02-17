package fetch

import (
	"fkmcps/structs"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const toolDescription = `Fetch web resource content from URL and return content in specified format.

## When to Use
Use this tool when you need to:
- Get raw content from a web page
- Access API endpoints to get JSON data
- Download HTML/text/Markdown content
- Quickly obtain web resources without complex processing

Don't use this tool when you need to:
- Extract specific information from a web page (should use specialized extraction tools)
- Analyze or summarize web page content (should fetch first then analyze)

## Features
- Supports four output formats: text (plain text), markdown (Markdown format), html (HTML format), json (JSON format)
- Automatically handles HTTP redirects
- Automatically extracts plain text from HTML (text format)
- Automatically converts HTML to Markdown (markdown format)
- Sets reasonable timeout to prevent long waits
- Limits response size (maximum 5MB) to prevent memory overflow

## Usage Tips
- text format: Suitable for getting plain text content or extracting text from HTML
- markdown format: Suitable for content that needs formatted rendering
- html format: Suitable for scenarios requiring raw HTML structure
- json format: Suitable for JSON data returned by API endpoints
- Set appropriate timeout based on website speed (default 30 seconds, maximum 120 seconds)`

func GetTools(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "fetch",
		Description: toolDescription,
	}, structs.WarpToolFunc(Fetch))
}
