package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cassin01/ghm/internal/git"
	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func removeCommand(c *cli.Context, cfg *config.Config) error {
	if c.NArg() < 1 {
		return fmt.Errorf("repository path is required")
	}

	repoPath := c.Args().Get(0)
	fullPath := filepath.Join(cfg.Root, repoPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("repository does not exist: %s", repoPath)
	}

	if !git.IsGitRepository(fullPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	fmt.Printf("Removing repository: %s\n", repoPath)

	if err := os.RemoveAll(fullPath); err != nil {
		return fmt.Errorf("failed to remove repository: %w", err)
	}

	fmt.Printf("Successfully removed: %s\n", repoPath)
	return nil
}
