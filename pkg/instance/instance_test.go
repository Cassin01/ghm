package instance

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseInstanceFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "No instance",
			path:     "github.com/user/repo",
			expected: 0,
		},
		{
			name:     "Instance 1",
			path:     "github.com/user/repo.1",
			expected: 1,
		},
		{
			name:     "Instance 2",
			path:     "github.com/user/repo.2",
			expected: 2,
		},
		{
			name:     "Instance 10",
			path:     "github.com/user/repo.10",
			expected: 10,
		},
		{
			name:     "Non-numeric suffix",
			path:     "github.com/user/repo.git",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInstanceFromPath(tt.path)
			if got != tt.expected {
				t.Errorf("ParseInstanceFromPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSaveAndLoadInstanceInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	
	info := &InstanceInfo{
		URL:         "https://github.com/user/repo.git",
		Instance:    1,
		CreatedAt:   testTime,
		LastUpdated: testTime,
	}

	t.Run("Save instance info", func(t *testing.T) {
		err := SaveInstanceInfo(tempDir, info)
		if err != nil {
			t.Errorf("SaveInstanceInfo() error = %v", err)
		}

		infoPath := filepath.Join(tempDir, ".ghm")
		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			t.Error("Instance info file was not created")
		}
	})

	t.Run("Load instance info", func(t *testing.T) {
		loaded, err := LoadInstanceInfo(tempDir)
		if err != nil {
			t.Errorf("LoadInstanceInfo() error = %v", err)
		}

		if loaded == nil {
			t.Fatal("LoadInstanceInfo() returned nil")
		}

		if loaded.URL != info.URL {
			t.Errorf("URL = %v, want %v", loaded.URL, info.URL)
		}
		if loaded.Instance != info.Instance {
			t.Errorf("Instance = %v, want %v", loaded.Instance, info.Instance)
		}
	})

	t.Run("Load non-existent instance info", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "non-existent")
		loaded, err := LoadInstanceInfo(nonExistentDir)
		if err != nil {
			t.Errorf("LoadInstanceInfo() error = %v", err)
		}
		if loaded != nil {
			t.Errorf("Expected nil for non-existent info, got %v", loaded)
		}
	})

	t.Run("Load invalid JSON", func(t *testing.T) {
		invalidDir := filepath.Join(tempDir, "invalid")
		os.MkdirAll(invalidDir, 0755)
		
		infoPath := filepath.Join(invalidDir, ".ghm")
		os.WriteFile(infoPath, []byte("invalid json"), 0644)

		loaded, err := LoadInstanceInfo(invalidDir)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
		if loaded != nil {
			t.Error("Expected nil for invalid JSON")
		}
	})
}

func TestFindNextInstance(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ghm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	host := "github.com"
	owner := "user"
	name := "repo"

	// Create test directory structure
	baseDir := filepath.Join(tempDir, host, owner)
	os.MkdirAll(baseDir, 0755)

	t.Run("No existing instances", func(t *testing.T) {
		next, err := FindNextInstance(tempDir, host, owner, name)
		if err != nil {
			t.Errorf("FindNextInstance() error = %v", err)
		}
		if next != 1 {
			t.Errorf("FindNextInstance() = %v, want 1", next)
		}
	})

	t.Run("With existing main instance", func(t *testing.T) {
		mainDir := filepath.Join(baseDir, name)
		os.MkdirAll(mainDir, 0755)

		next, err := FindNextInstance(tempDir, host, owner, name)
		if err != nil {
			t.Errorf("FindNextInstance() error = %v", err)
		}
		if next != 1 {
			t.Errorf("FindNextInstance() = %v, want 1", next)
		}
	})

	t.Run("With existing numbered instances", func(t *testing.T) {
		// Create instance 1 and 2
		instance1Dir := filepath.Join(baseDir, name+".1")
		instance2Dir := filepath.Join(baseDir, name+".2")
		os.MkdirAll(instance1Dir, 0755)
		os.MkdirAll(instance2Dir, 0755)

		next, err := FindNextInstance(tempDir, host, owner, name)
		if err != nil {
			t.Errorf("FindNextInstance() error = %v", err)
		}
		if next != 3 {
			t.Errorf("FindNextInstance() = %v, want 3", next)
		}
	})
}