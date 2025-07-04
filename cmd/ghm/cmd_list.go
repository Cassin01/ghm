package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cassin01/ghm/internal/git"
	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func listCommand(c *cli.Context, cfg *config.Config) error {
	pattern := c.Args().Get(0)
	showBranch := c.Bool("branch")

	repositories, err := findRepositories(cfg.Root, pattern)
	if err != nil {
		return fmt.Errorf("failed to find repositories: %w", err)
	}

	for _, repo := range repositories {
		if showBranch {
			branch, err := git.GetCurrentBranch(filepath.Join(cfg.Root, repo))
			if err != nil {
				fmt.Printf("%s [N/A]\n", repo)
			} else {
				fmt.Printf("%s [%s]\n", repo, branch)
			}
		} else {
			fmt.Println(repo)
		}
	}

	return nil
}

func findRepositories(root, pattern string) ([]string, error) {
	var repositories []string

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return repositories, nil
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if git.IsGitRepository(path) {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			relPath = strings.ReplaceAll(relPath, "\\", "/")

			if pattern == "" || strings.Contains(relPath, pattern) {
				repositories = append(repositories, relPath)
			}
		}

		return nil
	})

	return repositories, err
}
