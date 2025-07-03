package git

import (
	"fmt"
	"os"
	"os/exec"
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
