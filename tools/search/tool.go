package search

import (
	"context"
	"fkmcps/structs"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const searchToolDescription = `Search for information using DuckDuckGo search engine.

## When to Use
Use this tool when you need to:
- Search for the latest information on the internet
- Find materials on specific topics
- Get search results within a specified time range

## Usage Tips
- Provide clear, specific search keywords for better results
- You can use the time_range parameter to limit search results to a specific time period`

func GetTools(s *mcp.Server) {
	search, err := NewDuckDuckGoSearch(context.Background())
	if err != nil {
		log.Printf("failed to create search tool: %v", err)
		return
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search",
		Description: searchToolDescription,
	}, structs.WarpToolFunc(search.TextSearch))
}
