package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Cassin01/ghm/internal/git"
	"github.com/Cassin01/ghm/pkg/config"
	"github.com/Cassin01/ghm/pkg/instance"
	"github.com/Cassin01/ghm/pkg/repository"
	"github.com/urfave/cli/v2"
)

func getCommand(c *cli.Context, cfg *config.Config) error {
	if c.NArg() < 1 {
		return fmt.Errorf("repository URL is required")
	}

	repoURL := c.Args().Get(0)

	repo, err := repository.ParseURL(repoURL)
	if err != nil {
		return fmt.Errorf("failed to parse repository URL: %w", err)
	}

	instanceNumber := c.Int("number")
	autoAssign := c.Bool("auto")

	if autoAssign {
		nextInstance, err := instance.FindNextInstance(cfg.Root, repo.Host, repo.Owner, repo.Name)
		if err != nil {
			return fmt.Errorf("failed to find next instance: %w", err)
		}
		instanceNumber = nextInstance
	}

	repo.Instance = instanceNumber
	repoPath := repo.FullPath(cfg.Root)

	if _, err := os.Stat(repoPath); err == nil {
		return fmt.Errorf("repository already exists: %s", repoPath)
	}

	fmt.Printf("Cloning %s to %s\n", repo.URL, repoPath)

	if err := git.Clone(repo.URL, repoPath); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	info := &instance.InstanceInfo{
		URL:         repo.URL,
		Instance:    repo.Instance,
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	}

	if err := instance.SaveInstanceInfo(repoPath, info); err != nil {
		return fmt.Errorf("failed to save instance info: %w", err)
	}

	fmt.Printf("Successfully cloned to %s\n", repoPath)
	return nil
}
