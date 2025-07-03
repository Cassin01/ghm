package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func TestRemoveCommand(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	cfg := &config.Config{
		Root:            tempDir,
		DefaultProtocol: "https",
	}

	// Create a test repository
	repoPath := filepath.Join(tempDir, "github.com", "user", "repo")
	gitPath := filepath.Join(repoPath, ".git")
	err = os.MkdirAll(gitPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test repo: %v", err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "remove",
				Action: func(c *cli.Context) error {
					return removeCommand(c, cfg)
				},
			},
		},
	}

	t.Run("Remove existing repository", func(t *testing.T) {
		// Verify repository exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			t.Fatal("Test repository should exist before removal")
		}

		err := app.Run([]string{"ghm", "remove", "github.com/user/repo"})
		if err != nil {
			t.Errorf("removeCommand() error = %v", err)
		}

		// Verify repository is removed
		if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
			t.Error("Repository should be removed")
		}
	})

	t.Run("Remove non-existent repository", func(t *testing.T) {
		err := app.Run([]string{"ghm", "remove", "github.com/user/nonexistent"})
		if err == nil {
			t.Error("Expected error for non-existent repository")
		}
	})

	t.Run("Remove without repository path", func(t *testing.T) {
		err := app.Run([]string{"ghm", "remove"})
		if err == nil {
			t.Error("Expected error when no repository path provided")
		}
	})

	t.Run("Remove non-git directory", func(t *testing.T) {
		// Create a non-git directory
		nonGitPath := filepath.Join(tempDir, "not-a-repo")
		err := os.MkdirAll(nonGitPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create non-git dir: %v", err)
		}

		err = app.Run([]string{"ghm", "remove", "not-a-repo"})
		if err == nil {
			t.Error("Expected error for non-git directory")
		}
	})
}
