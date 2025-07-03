package main

import (
	"fmt"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func rootCommand(c *cli.Context, cfg *config.Config) error {
	fmt.Println(cfg.Root)
	return nil
}
