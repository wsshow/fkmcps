package fetch

import (
	"context"
	"fkmcps/constants"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

const (
	// MaxResponseSize Maximum response body size (5MB)
	MaxResponseSize = 5 * 1024 * 1024
	// MaxTimeout Maximum timeout duration (120 seconds)
	MaxTimeout = 120
	// DefaultTimeout Default timeout duration (30 seconds)
	DefaultTimeout = 30
)

// FetchRequest HTTP request parameters
type FetchRequest struct {
	URL     string `json:"url" jsonschema:"required,description:URL address to fetch content from (must start with http:// or https://)"`
	Format  string `json:"format,omitempty" jsonschema:"description:Format of returned content (text/markdown/html/json), default text. text format automatically extracts plain text from HTML, markdown format converts HTML to markdown, html format returns raw HTML, json format keeps JSON as-is"`
	Timeout int    `json:"timeout,omitempty" jsonschema:"description:Request timeout in seconds, default 30 seconds, maximum 120 seconds"`
}

// FetchResponse HTTP response
type FetchResponse struct {
	Content      string `json:"content" jsonschema:"description:Response content (content processed according to format parameter)"`
	StatusCode   int    `json:"status_code,omitempty" jsonschema:"description:HTTP status code"`
	ContentType  string `json:"content_type,omitempty" jsonschema:"description:Original content type"`
	IsTruncated  bool   `json:"is_truncated,omitempty" jsonschema:"description:Whether content is truncated"`
	ErrorMessage string `json:"error_message,omitempty" jsonschema:"description:Error message"`
}

// Fetch Send HTTP request to get web resources
func Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	// Parameter validation
	if req.URL == "" {
		return &FetchResponse{ErrorMessage: "URL is required"}, nil
	}

	// Validate URL protocol
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		return &FetchResponse{ErrorMessage: "URL must start with http:// or https://"}, nil
	}

	// Set default values and limits
	if req.Timeout == 0 {
		req.Timeout = DefaultTimeout
	} else if req.Timeout > MaxTimeout {
		req.Timeout = MaxTimeout
	}

	// Set default format
	format := strings.ToLower(req.Format)
	if format == "" {
		format = "text"
	}
	if format != "text" && format != "markdown" && format != "html" && format != "json" {
		return &FetchResponse{ErrorMessage: "format must be one of: text, markdown, html, json"}, nil
	}

	// Create HTTP client
	client := createHTTPClient(req.Timeout)

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", req.URL, nil)
	if err != nil {
		return &FetchResponse{ErrorMessage: fmt.Sprintf("failed to create request: %v", err)}, nil
	}

	// Set User-Agent
	httpReq.Header.Set("User-Agent", "FKTEAMS/1.0")

	// Send request
	resp, err := client.Do(httpReq)
	if err != nil {
		return &FetchResponse{ErrorMessage: fmt.Sprintf("failed to fetch URL: %v", err)}, nil
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return &FetchResponse{
			StatusCode:   resp.StatusCode,
			ErrorMessage: fmt.Sprintf("request failed with status code: %d", resp.StatusCode),
		}, nil
	}

	// Read response body (limit size)
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseSize))
	if err != nil {
		return &FetchResponse{
			StatusCode:   resp.StatusCode,
			ErrorMessage: fmt.Sprintf("failed to read response body: %v", err),
		}, nil
	}

	content := string(body)
	contentType := resp.Header.Get("Content-Type")

	// Validate UTF-8 encoding
	if !utf8.ValidString(content) {
		return &FetchResponse{
			StatusCode:   resp.StatusCode,
			ContentType:  contentType,
			ErrorMessage: "response content is not valid UTF-8",
		}, nil
	}

	// Process content according to format
	processedContent, err := processContent(content, contentType, format)
	if err != nil {
		return &FetchResponse{
			StatusCode:   resp.StatusCode,
			ContentType:  contentType,
			ErrorMessage: fmt.Sprintf("failed to process content: %v", err),
		}, nil
	}

	// Check if truncated
	isTruncated := int64(len(body)) >= MaxResponseSize

	return &FetchResponse{
		Content:     processedContent,
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		IsTruncated: isTruncated,
	}, nil
}

// processContent Process content according to format
func processContent(content, contentType, format string) (string, error) {
	isHTML := strings.Contains(contentType, "text/html")

	switch format {
	case "text":
		if isHTML {
			return extractTextFromHTML(content)
		}
		return content, nil

	case "markdown":
		if isHTML {
			return convertHTMLToMarkdown(content)
		}
		// Wrap non-HTML content in code block
		return "```\n" + content + "\n```", nil

	case "html":
		if isHTML {
			// Return only body part
			return extractHTMLBody(content)
		}
		return content, nil

	case "json":
		return content, nil

	default:
		return content, nil
	}
}

// extractTextFromHTML Extract plain text from HTML
func extractTextFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Remove script and style tags
	doc.Find("script, style").Remove()

	// Extract text content
	text := doc.Find("body").Text()

	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n"), nil
}

// convertHTMLToMarkdown Convert HTML to Markdown
func convertHTMLToMarkdown(html string) (string, error) {
	markdown, err := md.ConvertString(html)
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	return markdown, nil
}

// extractHTMLBody Extract HTML body part
func extractHTMLBody(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	body, err := doc.Find("body").Html()
	if err != nil {
		return "", fmt.Errorf("failed to extract body: %w", err)
	}

	if body == "" {
		return html, nil // If no body tag, return original content
	}

	return "<html>\n<body>\n" + body + "\n</body>\n</html>", nil
}

// createHTTPClient Create HTTP client
func createHTTPClient(timeoutSec int) *http.Client {
	proxyStr := os.Getenv(constants.MCP_PROXY_URL)
	var proxyFunc func(*http.Request) (*url.URL, error)

	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err == nil {
			proxyFunc = http.ProxyURL(proxyURL)
		}
	} else {
		proxyFunc = http.ProxyFromEnvironment
	}

	transport := &http.Transport{
		Proxy:                 proxyFunc,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSec) * time.Second,
	}
}
