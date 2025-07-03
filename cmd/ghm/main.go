package main

import (
	"fmt"
	"os"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func main() {
	cfg := config.New()

	app := &cli.App{
		Name:  "ghm",
		Usage: "GitHub Manager - manage multiple instances of the same repository",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Clone repository",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "number",
						Aliases: []string{"n"},
						Usage:   "Instance number",
					},
					&cli.BoolFlag{
						Name:  "auto",
						Usage: "Auto-assign next available instance number",
					},
				},
				Action: func(c *cli.Context) error {
					return getCommand(c, cfg)
				},
			},
			{
				Name:  "list",
				Usage: "List repositories",
				Action: func(c *cli.Context) error {
					return listCommand(c, cfg)
				},
			},
			{
				Name:  "root",
				Usage: "Show root directory",
				Action: func(c *cli.Context) error {
					return rootCommand(c, cfg)
				},
			},
			{
				Name:  "remove",
				Usage: "Remove repository",
				Action: func(c *cli.Context) error {
					return removeCommand(c, cfg)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
