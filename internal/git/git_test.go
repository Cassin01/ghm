package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(string) error
		expected bool
	}{
		{
			name: "Valid git repository",
			setup: func(path string) error {
				gitDir := filepath.Join(path, ".git")
				return os.MkdirAll(gitDir, 0755)
			},
			expected: true,
		},
		{
			name: "Directory without .git",
			setup: func(path string) error {
				return nil
			},
			expected: false,
		},
		{
			name: "Non-existent directory",
			setup: func(path string) error {
				return nil
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "ghm-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			testPath := tempDir
			if tt.name != "Non-existent directory" {
				if err := tt.setup(testPath); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			} else {
				testPath = filepath.Join(tempDir, "non-existent")
			}

			got := IsGitRepository(testPath)
			if got != tt.expected {
				t.Errorf("IsGitRepository() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClone(t *testing.T) {
	t.Run("Invalid destination directory parent", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "ghm-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		invalidPath := filepath.Join(tempDir, "non-existent-parent", "repo")
		
		err = Clone("https://invalid-url.com/repo.git", invalidPath)
		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
		}
	})

	t.Run("Create destination directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "ghm-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		destPath := filepath.Join(tempDir, "test-repo")
		
		err = Clone("https://invalid-url.com/repo.git", destPath)
		
		if _, statErr := os.Stat(destPath); os.IsNotExist(statErr) {
			t.Error("Expected destination directory to be created")
		}
	})
}