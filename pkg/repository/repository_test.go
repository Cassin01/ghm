package repository

import (
	"path/filepath"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected *Repository
		wantErr  bool
	}{
		{
			name: "GitHub HTTPS URL",
			url:  "https://github.com/user/repo",
			expected: &Repository{
				URL:      "https://github.com/user/repo",
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 0,
			},
			wantErr: false,
		},
		{
			name: "GitHub HTTPS URL with .git",
			url:  "https://github.com/user/repo.git",
			expected: &Repository{
				URL:      "https://github.com/user/repo.git",
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 0,
			},
			wantErr: false,
		},
		{
			name: "GitHub short URL",
			url:  "github.com/user/repo",
			expected: &Repository{
				URL:      "https://github.com/user/repo",
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 0,
			},
			wantErr: false,
		},
		{
			name:     "Invalid URL",
			url:      "invalid",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "URL with only scheme",
			url:      "https://",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "URL with insufficient path components",
			url:      "https://github.com/user",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "SSH URL",
			url:  "git@github.com:user/repo.git",
			expected: &Repository{
				URL:      "git@github.com:user/repo.git",
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				if got.URL != tt.expected.URL ||
					got.Host != tt.expected.Host ||
					got.Owner != tt.expected.Owner ||
					got.Name != tt.expected.Name ||
					got.Instance != tt.expected.Instance {
					t.Errorf("ParseURL() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestRepository_Path(t *testing.T) {
	tests := []struct {
		name     string
		repo     *Repository
		expected string
	}{
		{
			name: "Instance 0",
			repo: &Repository{
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 0,
			},
			expected: filepath.Join("github.com", "user", "repo"),
		},
		{
			name: "Instance 1",
			repo: &Repository{
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 1,
			},
			expected: filepath.Join("github.com", "user", "repo.1"),
		},
		{
			name: "Instance 2",
			repo: &Repository{
				Host:     "github.com",
				Owner:    "user",
				Name:     "repo",
				Instance: 2,
			},
			expected: filepath.Join("github.com", "user", "repo.2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repo.Path()
			if got != tt.expected {
				t.Errorf("Repository.Path() = %v, want %v", got, tt.expected)
			}
		})
	}
}
