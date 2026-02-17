package cmd

import (
	cli "github.com/urfave/cli/v3"
)

// NewApp creates the root CLI command for the MCP application.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:    "fkmcps",
		Usage:   "FeiKong MCP Server/Client over HTTP using the streamable transport",
		Version: "1.0.0",
		Commands: []*cli.Command{
			newServerCommand(),
			newClientCommand(),
			newUpdateCommand(),
		},
	}
}
