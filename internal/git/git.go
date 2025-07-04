package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Clone(url, destination string) error {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	cmd := exec.Command("git", "clone", url, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func IsGitRepository(path string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/.git", path))
	return err == nil
}

func GetCurrentBranch(path string) (string, error) {
	if !IsGitRepository(path) {
		return "", fmt.Errorf("not a git repository: %s", path)
	}

	cmd := exec.Command("git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))

	// Handle detached HEAD state
	if branch == "HEAD" {
		// Try to get the commit hash for detached HEAD
		cmd = exec.Command("git", "-C", path, "rev-parse", "--short", "HEAD")
		output, err = cmd.Output()
		if err != nil {
			return "HEAD", nil
		}
		return fmt.Sprintf("HEAD@%s", strings.TrimSpace(string(output))), nil
	}

	return branch, nil
}
