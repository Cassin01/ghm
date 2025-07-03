package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Cassin01/ghm/pkg/config"
	"github.com/urfave/cli/v2"
)

func TestListCommand(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	cfg := &config.Config{
		Root:            tempDir,
		DefaultProtocol: "https",
	}

	// Create test repositories
	repos := []string{
		"github.com/user/repo1",
		"github.com/user/repo2",
		"github.com/user/repo2.1",
		"gitlab.com/user/project",
	}

	for _, repo := range repos {
		repoPath := filepath.Join(tempDir, repo)
		gitPath := filepath.Join(repoPath, ".git")
		err := os.MkdirAll(gitPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test repo %s: %v", repo, err)
		}
	}

	t.Run("List all repositories", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		app := &cli.App{
			Commands: []*cli.Command{
				{
					Name: "list",
					Action: func(c *cli.Context) error {
						return listCommand(c, cfg)
					},
				},
			},
		}

		err := app.Run([]string{"ghm", "list"})

		_ = w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)

		if err != nil {
			t.Errorf("listCommand() error = %v", err)
		}

		output := buf.String()
		for _, repo := range repos {
			if !bytes.Contains([]byte(output), []byte(repo)) {
				t.Errorf("Expected output to contain %s, got: %s", repo, output)
			}
		}
	})

	t.Run("List repositories with pattern", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		app := &cli.App{
			Commands: []*cli.Command{
				{
					Name: "list",
					Action: func(c *cli.Context) error {
						return listCommand(c, cfg)
					},
				},
			},
		}

		err := app.Run([]string{"ghm", "list", "repo2"})

		_ = w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)

		if err != nil {
			t.Errorf("listCommand() error = %v", err)
		}

		output := buf.String()

		// Should contain repo2 and repo2.1
		if !bytes.Contains([]byte(output), []byte("github.com/user/repo2")) {
			t.Errorf("Expected output to contain repo2, got: %s", output)
		}
		if !bytes.Contains([]byte(output), []byte("github.com/user/repo2.1")) {
			t.Errorf("Expected output to contain repo2.1, got: %s", output)
		}

		// Should not contain repo1 or gitlab project
		if bytes.Contains([]byte(output), []byte("github.com/user/repo1")) {
			t.Errorf("Expected output to not contain repo1, got: %s", output)
		}
		if bytes.Contains([]byte(output), []byte("gitlab.com/user/project")) {
			t.Errorf("Expected output to not contain gitlab project, got: %s", output)
		}
	})
}

func TestFindRepositories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test repositories
	repos := []string{
		"github.com/user/repo1",
		"github.com/user/repo2",
		"gitlab.com/user/project",
	}

	for _, repo := range repos {
		repoPath := filepath.Join(tempDir, repo)
		gitPath := filepath.Join(repoPath, ".git")
		err := os.MkdirAll(gitPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test repo %s: %v", repo, err)
		}
	}

	// Create a non-git directory
	nonGitPath := filepath.Join(tempDir, "not-a-repo")
	err = os.MkdirAll(nonGitPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create non-git dir: %v", err)
	}

	t.Run("Find all repositories", func(t *testing.T) {
		found, err := findRepositories(tempDir, "")
		if err != nil {
			t.Errorf("findRepositories() error = %v", err)
		}

		if len(found) != len(repos) {
			t.Errorf("findRepositories() found %d repos, want %d", len(found), len(repos))
		}

		for _, repo := range repos {
			repoFound := false
			for _, f := range found {
				if f == repo {
					repoFound = true
					break
				}
			}
			if !repoFound {
				t.Errorf("Expected to find %s in results", repo)
			}
		}
	})

	t.Run("Find repositories with pattern", func(t *testing.T) {
		found, err := findRepositories(tempDir, "github")
		if err != nil {
			t.Errorf("findRepositories() error = %v", err)
		}

		expectedCount := 2 // repo1 and repo2
		if len(found) != expectedCount {
			t.Errorf("findRepositories() found %d repos, want %d", len(found), expectedCount)
		}
	})

	t.Run("Non-existent root", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non-existent")
		found, err := findRepositories(nonExistentPath, "")
		if err != nil {
			t.Errorf("findRepositories() error = %v", err)
		}

		if len(found) != 0 {
			t.Errorf("findRepositories() found %d repos, want 0", len(found))
		}
	})
}
