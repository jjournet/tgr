package profile

// a profile is either a user or an organization

import (
	"context"

	"github.com/google/go-github/v61/github"
)

type Profile struct {
	owner  string
	Org    string
	repos  []string
	Client *github.Client
}

func (p *Profile) Owner() string {
	return p.owner
}

func (p *Profile) Repos() []string {
	return p.repos
}

func NewProfile(owner string, org string, client *github.Client) *Profile {
	//list option to get all repositories
	listOpt := &github.ListOptions{PerPage: 100}
	opts := &github.RepositoryListByOrgOptions{ListOptions: *listOpt}
	repoList, _, err := client.Repositories.ListByOrg(context.Background(), org, opts)
	if err != nil {
		panic(err)
	}
	repos := make([]string, len(repoList))
	for i, repo := range repoList {
		repos[i] = *repo.Name
	}
	return &Profile{owner: owner, Org: org, repos: repos, Client: client}
}
