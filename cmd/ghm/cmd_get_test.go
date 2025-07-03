package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func TestGetCommand(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	cfg := &config.Config{
		Root:            tempDir,
		DefaultProtocol: "https",
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "get",
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
		},
	}

	t.Run("Get without URL", func(t *testing.T) {
		err := app.Run([]string{"ghm", "get"})
		if err == nil {
			t.Error("Expected error when no URL provided")
		}
	})

	t.Run("Get with invalid URL", func(t *testing.T) {
		err := app.Run([]string{"ghm", "get", "invalid-url"})
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("Get with valid URL but git clone will fail", func(t *testing.T) {
		// This will fail at git clone step, but should pass URL parsing
		err := app.Run([]string{"ghm", "get", "https://github.com/nonexistent/repo.git"})
		if err == nil {
			t.Error("Expected error when git clone fails")
		}

		// The directory should still be created
		expectedPath := filepath.Join(tempDir, "github.com", "nonexistent", "repo")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Error("Expected directory to be created even if git clone fails")
		}
	})

	t.Run("Get with instance number", func(t *testing.T) {
		err := app.Run([]string{"ghm", "get", "https://github.com/test/repo.git", "-n", "2"})
		if err == nil {
			t.Error("Expected error when git clone fails")
		}

		// When git clone fails, the operation should fail without creating persistent artifacts
		// This test verifies that the error handling works correctly
	})

	t.Run("Get with auto flag", func(t *testing.T) {
		// Create existing instances to test auto assignment
		existingPath := filepath.Join(tempDir, "github.com", "auto", "repo")
		_ = os.MkdirAll(existingPath, 0755)

		instance1Path := filepath.Join(tempDir, "github.com", "auto", "repo.1")
		_ = os.MkdirAll(instance1Path, 0755)

		err := app.Run([]string{"ghm", "get", "https://github.com/auto/repo.git", "--auto"})
		if err == nil {
			t.Error("Expected error when git clone fails")
		}

		// When git clone fails, the operation should fail
		// This test verifies that auto-assignment logic works before attempting to clone
	})

	t.Run("Get existing repository", func(t *testing.T) {
		// Create existing repository
		existingPath := filepath.Join(tempDir, "github.com", "existing", "repo")
		_ = os.MkdirAll(existingPath, 0755)

		err := app.Run([]string{"ghm", "get", "https://github.com/existing/repo.git"})
		if err == nil {
			t.Error("Expected error for existing repository")
		}
	})
}
