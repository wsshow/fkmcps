package structs

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolFunc[I any, O any] func(ctx context.Context, input I) (output O, err error)

func WarpToolFunc[I any, O any](toolFunc ToolFunc[I, O]) mcp.ToolHandlerFor[I, O] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input I) (_ *mcp.CallToolResult, output O, _ error) {
		result, err := toolFunc(ctx, input)
		return nil, result, err
	}
}
