package github

import (
	"context"

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

// Context returns a background context for API calls
// Can be extended later to support cancellation
func (s *GitHubService) Context() context.Context {
	return context.Background()
}
