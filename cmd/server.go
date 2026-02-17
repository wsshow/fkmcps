package cmd

import (
	"context"
	"errors"
	"fkmcps/middlewares"
	"fkmcps/tools/doc"
	"fkmcps/tools/fetch"
	"fkmcps/tools/search"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	cli "github.com/urfave/cli/v3"
)

// toolInfo describes an available tool group.
type toolInfo struct {
	Name        string
	Description string
	Register    func(s *mcp.Server)
}

// availableTools is the registry of all tool groups.
var availableTools = []toolInfo{
	{Name: "doc", Description: "Document Tools (get_document_info, read_document_smart, read_document_by_page, read_document_by_line)", Register: doc.GetTools},
	{Name: "fetch", Description: "Web Fetch Tools (fetch)", Register: fetch.GetTools},
	{Name: "search", Description: "Web Search Tools (search)", Register: search.GetTools},
}

// allToolNames returns a slice of all available tool names.
func allToolNames() []string {
	names := make([]string, len(availableTools))
	for i, t := range availableTools {
		names[i] = t.Name
	}
	return names
}

func newServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the MCP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "Host to listen on",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 8000,
				Usage: "Port number to listen on",
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Value:   false,
				Usage:   "Interactively select which tools to enable",
			},
			&cli.StringSliceFlag{
				Name:  "tools",
				Usage: "Comma-separated list of tools to enable (e.g. doc,fetch,search). Defaults to all.",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String("host")
			port := cmd.Int("port")
			addr := fmt.Sprintf("%s:%d", host, port)

			var selectedTools []string

			if cmd.Bool("interactive") {
				selected, err := selectToolsInteractively()
				if err != nil {
					return err
				}
				selectedTools = selected
			} else if tools := cmd.StringSlice("tools"); len(tools) > 0 {
				selectedTools = tools
			} else {
				selectedTools = allToolNames()
			}

			return runServer(addr, selectedTools)
		},
	}
}

// selectToolsInteractively displays a multi-select TUI for choosing tools.
func selectToolsInteractively() ([]string, error) {
	var selected []string

	options := make([]huh.Option[string], len(availableTools))
	for i, t := range availableTools {
		options[i] = huh.NewOption(fmt.Sprintf("%s - %s", t.Name, t.Description), t.Name).Selected(true)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Tools").
				Description("Use space to toggle selection, enter to confirm").
				Options(options...).
				Validate(func(selected []string) error {
					if len(selected) == 0 {
						return errors.New("At least one tool must be selected")
					}
					return nil
				}).
				Value(&selected),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, fmt.Errorf("User aborted the selection")
		}
		return nil, fmt.Errorf("Interactive selection failed: %w", err)
	}

	return selected, nil
}

func runServer(addr string, enabledTools []string) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "feikong-mcp-server",
		Version: "1.0.0",
	}, nil)

	server.AddReceivingMiddleware(middlewares.Logger())

	enabled := make(map[string]bool, len(enabledTools))
	for _, name := range enabledTools {
		enabled[name] = true
	}

	var registered []string
	for _, t := range availableTools {
		if enabled[t.Name] {
			t.Register(server)
			registered = append(registered, t.Name)
		}
	}

	log.Printf("Enabled tools: [%s]", strings.Join(registered, ", "))

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	log.Printf("MCP server listening on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}
