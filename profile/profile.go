package profile

// a profile is either a user or an organization

import (
	"context"
	"log"

	"github.com/google/go-github/v61/github"
)

type Repo struct {
	Name        string
	Description string
}
type Profile struct {
	owner    string
	Profile  string
	repos    []string
	RepoList []Repo
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
	rList := make([]Repo, len(repoList))
	final := 0
	for i, repo := range repoList {
		repos[i] = *repo.Name
		if repo.Description == nil {
			repo.Description = new(string)
			*repo.Description = "N/A"
		}
		rList[i] = Repo{Name: *repo.Name, Description: *repo.Description}
		final = i
	}
	log.Printf("final repo count: %d", final)
	return &Profile{owner: owner, Profile: org, repos: repos, RepoList: rList}
}
