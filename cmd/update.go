package cmd

import (
	"context"
	"fkmcps/update"

	cli "github.com/urfave/cli/v3"
)

func newUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Check for updates and download the latest release if available",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "proxy",
				Usage: "Proxy URL to use for downloading updates (e.g. http://127.0.0.1:7890)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			proxyURL := cmd.String("proxy")
			return update.SelfUpdate("wsshow", "fkmcps", proxyURL)
		},
	}
}
