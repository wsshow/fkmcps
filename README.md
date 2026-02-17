# FeiKong MCP Server

A Model Context Protocol (MCP) server providing tools for document reading, web fetching, and web searching over HTTP.

## Features

- **Document Tools**: Read and extract content from various document formats (PDF, DOCX, XLSX, PPTX, TXT, CSV, MD, RTF)
- **Web Fetch**: Retrieve web content from URLs in multiple formats (markdown, html, text)
- **Web Search**: Search the web using DuckDuckGo search engine

## Installation

```bash
go install
```

Or build from source:

```bash
go build -o fkmcps
```

## Usage

### Start Server

```bash
# Start with all tools enabled (default: localhost:8000)
fkmcps server

# Specify host and port
fkmcps server --host 0.0.0.0 --port 8080

# Enable specific tools only
fkmcps server --tools doc,search

# Interactive tool selection
fkmcps server --interactive
```

### Update

```bash
# Check for updates and install
fkmcps update
```

## Available Tools

### Document Tools

- `get_document_info` - Get document metadata (type, size, pages, etc.)
- `read_document_smart` - Intelligently read document content with automatic chunking
- `read_document_by_page` - Read specific page ranges
- `read_document_by_line` - Read specific line ranges

### Web Fetch Tools

- `fetch` - Fetch web content from URL with customizable output format

### Web Search Tools

- `search` - Search the web using DuckDuckGo

## Configuration

Server flags:

- `--host` - Host to listen on (default: `localhost`)
- `--port` - Port number (default: `8000`)
- `--tools` - Comma-separated list of tools to enable
- `--interactive` / `-i` - Interactive tool selection

Client flags:

- `--host` - Host to connect to (default: `localhost`)
- `--port` - Port number (default: `8000`)
- `--proto` - Protocol (default: `http`)

## License

MIT
