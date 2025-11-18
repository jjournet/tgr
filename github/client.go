package github

import (
	"context"
	"os/exec"
	"strings"

	gh "github.com/google/go-github/v69/github"
)

// GitHubService centralizes all GitHub API interactions
type GitHubService struct {
	client *gh.Client
}

// NewGitHubService creates a new GitHub service with the provided token
func NewGitHubService(token string) *GitHubService {
	return &GitHubService{
		client: gh.NewClient(nil).WithAuthToken(token),
	}
}

// NewGitHubServiceFromCLI creates a service using the gh CLI token
func NewGitHubServiceFromCLI() (*GitHubService, error) {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return nil, err
	}
	token := strings.TrimSuffix(string(out), "\n")
	return NewGitHubService(token), nil
}

// Context returns a background context for API calls
// Can be extended later to support cancellation
func (s *GitHubService) Context() context.Context {
	return context.Background()
}
