package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Root            string
	DefaultProtocol string
}

func New() *Config {
	return &Config{
		Root:            getDefaultRoot(),
		DefaultProtocol: "https",
	}
}

func getDefaultRoot() string {
	if root := os.Getenv("GHM_ROOT"); root != "" {
		return root
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "./ghm"
	}

	return filepath.Join(home, "ghm")
}
