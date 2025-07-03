package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "ghm")
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ghm")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Create temporary root directory
	tempRoot, err := os.MkdirTemp("", "ghm-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempRoot) }()

	// Set environment variable
	_ = os.Setenv("GHM_ROOT", tempRoot)
	defer func() { _ = os.Unsetenv("GHM_ROOT") }()

	t.Run("Root command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "root")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("root command failed: %v", err)
		}

		if strings.TrimSpace(string(output)) != tempRoot {
			t.Errorf("root command output = %q, want %q", strings.TrimSpace(string(output)), tempRoot)
		}
	})

	t.Run("List empty", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "list")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("list command failed: %v", err)
		}

		if len(strings.TrimSpace(string(output))) != 0 {
			t.Errorf("list command should return empty output, got: %q", string(output))
		}
	})

	t.Run("Get invalid URL", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "get", "invalid-url")
		_, err := cmd.Output()
		if err == nil {
			t.Error("get command should fail with invalid URL")
		}
	})

	t.Run("Get without arguments", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "get")
		_, err := cmd.Output()
		if err == nil {
			t.Error("get command should fail without arguments")
		}
	})

	t.Run("Remove non-existent repository", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "remove", "github.com/user/nonexistent")
		_, err := cmd.Output()
		if err == nil {
			t.Error("remove command should fail for non-existent repository")
		}
	})

	t.Run("Help command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("help command failed: %v", err)
		}

		helpText := string(output)
		if !strings.Contains(helpText, "GitHub Manager") {
			t.Error("help text should contain 'GitHub Manager'")
		}
		if !strings.Contains(helpText, "COMMANDS:") {
			t.Error("help text should contain 'COMMANDS:'")
		}
		if !strings.Contains(helpText, "get") {
			t.Error("help text should contain 'get' command")
		}
		if !strings.Contains(helpText, "list") {
			t.Error("help text should contain 'list' command")
		}
		if !strings.Contains(helpText, "root") {
			t.Error("help text should contain 'root' command")
		}
		if !strings.Contains(helpText, "remove") {
			t.Error("help text should contain 'remove' command")
		}
	})

	// Test workflow with mock repositories (using directories instead of real git repos)
	t.Run("Workflow with mock repositories", func(t *testing.T) {
		// Create mock repositories manually
		mockRepos := []string{
			"github.com/user/repo1",
			"github.com/user/repo2",
			"github.com/user/repo2.1",
			"gitlab.com/user/project",
		}

		for _, repo := range mockRepos {
			repoPath := filepath.Join(tempRoot, repo)
			gitPath := filepath.Join(repoPath, ".git")
			err := os.MkdirAll(gitPath, 0755)
			if err != nil {
				t.Fatalf("Failed to create mock repo %s: %v", repo, err)
			}
		}

		// Test list command
		cmd := exec.Command(binaryPath, "list")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("list command failed: %v", err)
		}

		outputStr := string(output)
		for _, repo := range mockRepos {
			if !strings.Contains(outputStr, repo) {
				t.Errorf("list output should contain %s, got: %s", repo, outputStr)
			}
		}

		// Test list with pattern
		cmd = exec.Command(binaryPath, "list", "repo2")
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("list command with pattern failed: %v", err)
		}

		outputStr = string(output)
		if !strings.Contains(outputStr, "github.com/user/repo2") {
			t.Error("list with pattern should contain repo2")
		}
		if !strings.Contains(outputStr, "github.com/user/repo2.1") {
			t.Error("list with pattern should contain repo2.1")
		}
		if strings.Contains(outputStr, "github.com/user/repo1") {
			t.Error("list with pattern should not contain repo1")
		}

		// Test remove command
		cmd = exec.Command(binaryPath, "remove", "github.com/user/repo1")
		_, err = cmd.Output()
		if err != nil {
			t.Fatalf("remove command failed: %v", err)
		}

		// Verify repository is removed
		removedPath := filepath.Join(tempRoot, "github.com/user/repo1")
		if _, err := os.Stat(removedPath); !os.IsNotExist(err) {
			t.Error("Repository should be removed")
		}

		// Verify other repositories still exist
		cmd = exec.Command(binaryPath, "list")
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("list command failed after remove: %v", err)
		}

		outputStr = string(output)
		if strings.Contains(outputStr, "github.com/user/repo1") {
			t.Error("removed repository should not appear in list")
		}
		if !strings.Contains(outputStr, "github.com/user/repo2") {
			t.Error("remaining repositories should still appear in list")
		}
	})
}
