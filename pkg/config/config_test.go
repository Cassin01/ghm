package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Test default configuration
	cfg := New()
	if cfg.DefaultProtocol != "https" {
		t.Errorf("Expected default protocol to be https, got %s", cfg.DefaultProtocol)
	}
	
	// Test with GHM_ROOT environment variable
	originalRoot := os.Getenv("GHM_ROOT")
	defer func() {
		if originalRoot != "" {
			os.Setenv("GHM_ROOT", originalRoot)
		} else {
			os.Unsetenv("GHM_ROOT")
		}
	}()
	
	testRoot := "/tmp/test-ghm"
	os.Setenv("GHM_ROOT", testRoot)
	
	cfg = New()
	if cfg.Root != testRoot {
		t.Errorf("Expected root to be %s, got %s", testRoot, cfg.Root)
	}
}

func TestGetDefaultRoot(t *testing.T) {
	// Test without GHM_ROOT environment variable
	originalRoot := os.Getenv("GHM_ROOT")
	defer func() {
		if originalRoot != "" {
			os.Setenv("GHM_ROOT", originalRoot)
		} else {
			os.Unsetenv("GHM_ROOT")
		}
	}()
	
	os.Unsetenv("GHM_ROOT")
	
	root := getDefaultRoot()
	home, err := os.UserHomeDir()
	if err != nil {
		if root != "./ghm" {
			t.Errorf("Expected root to be ./ghm when home dir is not available, got %s", root)
		}
	} else {
		expected := filepath.Join(home, "ghm")
		if root != expected {
			t.Errorf("Expected root to be %s, got %s", expected, root)
		}
	}
}