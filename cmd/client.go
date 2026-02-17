package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	cli "github.com/urfave/cli/v3"
)

func newClientCommand() *cli.Command {
	return &cli.Command{
		Name:  "client",
		Usage: "Connect to the MCP server as a client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "Host to connect to",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8000,
				Usage: "Port number to connect to",
			},
			&cli.StringFlag{
				Name:  "proto",
				Value: "http",
				Usage: "Protocol to use (http or https)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String("host")
			port := cmd.Int("port")
			proto := cmd.String("proto")
			url := fmt.Sprintf("%s://%s:%d", proto, host, port)
			return runClient(ctx, url)
		},
	}
}

func runClient(ctx context.Context, url string) error {
	log.Printf("Connecting to MCP server at %s", url)

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "feikong-mcp-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: url}, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer session.Close()

	log.Printf("Connected to server (session ID: %s)", session.ID())

	log.Println("Listing available tools...")
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	for _, tool := range toolsResult.Tools {
		log.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	fmt.Fprintln(os.Stdout, "Client completed successfully")
	return nil
}
