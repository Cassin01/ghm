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
				Usage: "Clone repository with optional instance management",
				Description: `Clone a repository to the ghm root directory. Multiple instances of the
same repository can be managed using instance numbers.

Examples:
  ghm get https://github.com/user/repo          # Clone to repo/
  ghm get https://github.com/user/repo -n 1     # Clone to repo_1/
  ghm get https://github.com/user/repo --auto   # Auto-assign next number`,
				ArgsUsage: "<repository-url>",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "number",
						Aliases: []string{"n"},
						Usage:   "Specify instance number (creates repo_N directory)",
					},
					&cli.BoolFlag{
						Name:  "auto",
						Usage: "Automatically assign next available instance number",
					},
				},
				Action: func(c *cli.Context) error {
					return getCommand(c, cfg)
				},
			},
			{
				Name:  "list",
				Usage: "List repositories",
				Description: "List all repositories managed by ghm, including all instances.",
				ArgsUsage: "[pattern]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "branch",
						Aliases: []string{"b"},
						Usage:   "Show current branch name for each repository",
					},
				},
				Action: func(c *cli.Context) error {
					return listCommand(c, cfg)
				},
			},
			{
				Name:  "root",
				Usage: "Show root directory",
				Description: "Display the ghm root directory path where repositories are stored.",
				Action: func(c *cli.Context) error {
					return rootCommand(c, cfg)
				},
			},
			{
				Name:  "remove",
				Usage: "Remove repository",
				Description: "Remove a repository instance from ghm management.",
				ArgsUsage: "<repository-path>",
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
