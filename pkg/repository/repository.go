package repository

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type Repository struct {
	URL      string
	Host     string
	Owner    string
	Name     string
	Instance int
}

func ParseURL(repoURL string) (*Repository, error) {
	if repoURL == "" {
		return nil, fmt.Errorf("repository URL cannot be empty")
	}
	
	// Handle SSH URLs (git@host:owner/repo.git)
	if strings.HasPrefix(repoURL, "git@") {
		return parseSSHURL(repoURL)
	}
	
	// Handle HTTP/HTTPS URLs
	if !strings.Contains(repoURL, "://") {
		repoURL = "https://" + repoURL
	}
	
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("invalid repository URL: %w", err)
	}
	
	if u.Host == "" {
		return nil, fmt.Errorf("invalid repository URL: missing host")
	}
	
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository path: %s", u.Path)
	}
	
	owner := parts[0]
	name := parts[1]
	
	if owner == "" || name == "" {
		return nil, fmt.Errorf("invalid repository path: owner and name cannot be empty")
	}
	
	// Remove .git suffix if present
	if strings.HasSuffix(name, ".git") {
		name = strings.TrimSuffix(name, ".git")
	}
	
	return &Repository{
		URL:      repoURL,
		Host:     u.Host,
		Owner:    owner,
		Name:     name,
		Instance: 0,
	}, nil
}

func parseSSHURL(sshURL string) (*Repository, error) {
	// Format: git@host:owner/repo.git
	if !strings.Contains(sshURL, ":") {
		return nil, fmt.Errorf("invalid SSH URL format: %s", sshURL)
	}
	
	parts := strings.Split(sshURL, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SSH URL format: %s", sshURL)
	}
	
	hostPart := parts[0]
	pathPart := parts[1]
	
	// Extract host from git@host
	if !strings.HasPrefix(hostPart, "git@") {
		return nil, fmt.Errorf("invalid SSH URL format: %s", sshURL)
	}
	host := strings.TrimPrefix(hostPart, "git@")
	
	// Extract owner/repo from path
	pathComponents := strings.Split(pathPart, "/")
	if len(pathComponents) < 2 {
		return nil, fmt.Errorf("invalid SSH URL path: %s", pathPart)
	}
	
	owner := pathComponents[0]
	name := pathComponents[1]
	
	if owner == "" || name == "" {
		return nil, fmt.Errorf("invalid SSH URL path: owner and name cannot be empty")
	}
	
	// Remove .git suffix if present
	if strings.HasSuffix(name, ".git") {
		name = strings.TrimSuffix(name, ".git")
	}
	
	return &Repository{
		URL:      sshURL,
		Host:     host,
		Owner:    owner,
		Name:     name,
		Instance: 0,
	}, nil
}

func (r *Repository) Path() string {
	if r.Instance == 0 {
		return filepath.Join(r.Host, r.Owner, r.Name)
	}
	return filepath.Join(r.Host, r.Owner, fmt.Sprintf("%s.%d", r.Name, r.Instance))
}

func (r *Repository) FullPath(root string) string {
	return filepath.Join(root, r.Path())
}