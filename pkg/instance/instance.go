package instance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type InstanceInfo struct {
	URL         string    `json:"url"`
	Instance    int       `json:"instance"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func SaveInstanceInfo(repoPath string, info *InstanceInfo) error {
	infoPath := filepath.Join(repoPath, ".ghm")

	if err := os.MkdirAll(filepath.Dir(infoPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal instance info: %w", err)
	}

	if err := os.WriteFile(infoPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write instance info: %w", err)
	}

	return nil
}

func LoadInstanceInfo(repoPath string) (*InstanceInfo, error) {
	infoPath := filepath.Join(repoPath, ".ghm")

	data, err := os.ReadFile(infoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read instance info: %w", err)
	}

	var info InstanceInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instance info: %w", err)
	}

	return &info, nil
}

func FindNextInstance(rootPath, host, owner, name string) (int, error) {
	basePattern := filepath.Join(rootPath, host, owner, name)

	maxInstance := 0

	// Check if base directory exists
	if _, err := os.Stat(basePattern); err == nil {
		maxInstance = 0
	}

	// Check numbered instances
	for i := 1; i <= 100; i++ {
		instancePath := fmt.Sprintf("%s.%d", basePattern, i)
		if _, err := os.Stat(instancePath); err == nil {
			maxInstance = i
		}
	}

	return maxInstance + 1, nil
}

func ParseInstanceFromPath(path string) int {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return 0
	}

	lastPart := parts[len(parts)-1]
	if instance, err := strconv.Atoi(lastPart); err == nil {
		return instance
	}

	return 0
}
